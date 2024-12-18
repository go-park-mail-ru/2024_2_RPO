package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Configuration holds the necessary configuration for PostgreSQL and Elasticsearch
type Configuration struct {
	PostgreSQL struct {
		URI     string
		Table   string
		Columns []string
	}
	Elasticsearch struct {
		Addresses   []string
		Index       string
		MappingFile string
	}
}

func loadConfig() Configuration {
	// For simplicity, configuration is hardcoded.
	// You can enhance this by reading from environment variables, flags, or config files.
	var cfg Configuration

	// PostgreSQL configuration
	cfg.PostgreSQL.URI = "postgres://username:password@localhost:5432/dbname"
	cfg.PostgreSQL.Table = "your_table"
	cfg.PostgreSQL.Columns = []string{"id", "name", "created_at"} // Customize columns as needed

	// Elasticsearch configuration
	cfg.Elasticsearch.Addresses = []string{
		"http://localhost:9200",
	}
	cfg.Elasticsearch.Index = "your_index"
	cfg.Elasticsearch.MappingFile = "mapping.json"

	return cfg
}

func main() {
	// Load configuration
	cfg := loadConfig()

	// Initialize Elasticsearch client
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: cfg.Elasticsearch.Addresses,
		// Add authentication if needed
		// Username: "user",
		// Password: "pass",
	})
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	ctx := context.Background()

	// Step 1: Delete existing index if it exists
	deleteIndex(cfg.Elasticsearch.Index, es)

	// Step 2: Create new index with mapping
	createIndexWithMapping(cfg.Elasticsearch.Index, cfg.Elasticsearch.MappingFile, es)

	// Step 3: Load data from PostgreSQL and index into Elasticsearch
	err = loadDataToElasticsearch(ctx, cfg, es)
	if err != nil {
		log.Fatalf("Error loading data to Elasticsearch: %s", err)
	}

	log.Println("Data successfully loaded to Elasticsearch")
}

func deleteIndex(index string, es *elasticsearch.Client) {
	// Check if the index exists
	res, err := es.Indices.Exists([]string{index})
	if err != nil {
		log.Fatalf("Error checking if index exists: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		// Delete the index
		res, err := es.Indices.Delete([]string{index}, es.Indices.Delete.WithIgnoreUnavailable(true))
		if err != nil {
			log.Fatalf("Error deleting index: %s", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Fatalf("Error response while deleting index: %s", res.String())
		}

		log.Printf("Deleted existing index: %s\n", index)
	} else if res.StatusCode == 404 {
		log.Printf("Index %s does not exist. No need to delete.\n", index)
	} else {
		log.Fatalf("Unexpected response while checking index existence: %s", res.String())
	}
}

func createIndexWithMapping(index, mappingFile string, es *elasticsearch.Client) {
	// Read mapping JSON from file
	mapping, err := os.ReadFile(mappingFile)
	if err != nil {
		log.Fatalf("Error reading mapping file: %s", err)
	}

	// Create the index with the provided mapping
	res, err := es.Indices.Create(
		index,
		es.Indices.Create.WithBody(bytes.NewReader(mapping)),
	)
	if err != nil {
		log.Fatalf("Error creating index: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error response while creating index: %s", res.String())
	}

	log.Printf("Created index %s with mapping.\n", index)
}

func loadDataToElasticsearch(ctx context.Context, cfg Configuration, es *elasticsearch.Client) error {
	// Initialize PostgreSQL connection pool
	pool, err := pgxpool.New(ctx, cfg.PostgreSQL.URI)
	if err != nil {
		return fmt.Errorf("unable to create PostgreSQL pool: %w", err)
	}
	defer pool.Close()

	// Query to fetch all data
	columns := strings.Join(cfg.PostgreSQL.Columns, ", ")
	query := fmt.Sprintf("SELECT %s FROM %s", columns, cfg.PostgreSQL.Table)

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	// Prepare for bulk indexing
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	batchSize := 1000
	count := 0

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return fmt.Errorf("failed to get row values: %w", err)
		}

		// Assume first column is the ID
		docID := fmt.Sprintf("%v", values[0])

		// Convert row to a map for JSON encoding
		doc := make(map[string]interface{})
		for i, col := range cfg.PostgreSQL.Columns {
			doc[col] = values[i]
		}

		// Create the action metadata
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": cfg.Elasticsearch.Index,
				"_id":    docID,
			},
		}
		metaLine, err := json.Marshal(meta)
		if err != nil {
			return fmt.Errorf("failed to marshal meta line: %w", err)
		}

		dataLine, err := json.Marshal(doc)
		if err != nil {
			return fmt.Errorf("failed to marshal data line: %w", err)
		}

		// Write to buffer
		_, err = writer.Write(metaLine)
		if err != nil {
			return fmt.Errorf("elastic: %w", err)
		}
		err = writer.WriteByte('\n')
		if err != nil {
			return fmt.Errorf("elastic: %w", err)
		}
		_, err = writer.Write(dataLine)
		if err != nil {
			return fmt.Errorf("elastic: %w", err)
		}
		err = writer.WriteByte('\n')
		if err != nil {
			return fmt.Errorf("elastic: %w", err)
		}

		count++
		if count%batchSize == 0 {
			writer.Flush()
			// Perform bulk request
			res, err := es.Bulk(bytes.NewReader(buf.Bytes()), es.Bulk.WithContext(ctx))
			if err != nil {
				return fmt.Errorf("bulk request failed: %w", err)
			}
			defer res.Body.Close()

			if res.IsError() {
				// You might want to handle partial failures here
				return fmt.Errorf("bulk request error: %s", res.String())
			}

			// Reset buffer
			buf.Reset()
		}
	}

	// Flush remaining data
	if buf.Len() > 0 {
		writer.Flush()
		res, err := es.Bulk(bytes.NewReader(buf.Bytes()), es.Bulk.WithContext(ctx))
		if err != nil {
			return fmt.Errorf("final bulk request failed: %w", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			return fmt.Errorf("final bulk request error: %s", res.String())
		}
	}

	log.Printf("Successfully indexed %d documents.\n", count)
	return nil
}
