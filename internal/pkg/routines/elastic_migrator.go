package routines

import (
	"RPO_back/internal/pkg/board/repository"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/olivere/elastic/v7"
)

func ElasticMigrate(el *elastic.Client, db *pgxpool.Pool) {
	ctx := context.Background()

	err := deleteIndex(repository.ElasticIdxName, el, ctx)
	if err != nil {
		log.Fatalf("Error deleting index: %s", err)
	}

	err = createIndexWithMapping(repository.ElasticIdxName, "mapping.json", el, ctx)
	if err != nil {
		log.Fatalf("Error creating index with mapping: %s", err)
	}

	err = loadDataToElasticsearch(ctx, db, el, repository.ElasticIdxName)
	if err != nil {
		log.Fatalf("Error loading data to Elasticsearch: %s", err)
	}

	log.Println("Data successfully loaded to Elasticsearch.")
}

func deleteIndex(index string, el *elastic.Client, ctx context.Context) error {
	exists, err := el.IndexExists(index).Do(ctx)
	if err != nil {
		return fmt.Errorf("error checking if index exists: %w", err)
	}

	if exists {
		_, err := el.DeleteIndex(index).Do(ctx)
		if err != nil {
			return fmt.Errorf("error deleting index: %w", err)
		}
		log.Printf("Deleted existing index: %s\n", index)
	} else {
		log.Printf("Index %s does not exist. No need to delete.\n", index)
	}

	return nil
}

func createIndexWithMapping(index, mappingFile string, el *elastic.Client, ctx context.Context) error {
	mapping, err := os.ReadFile(mappingFile)
	if err != nil {
		return fmt.Errorf("error reading mapping file: %w", err)
	}

	_, err = el.CreateIndex(index).BodyString(string(mapping)).Do(ctx)
	if err != nil {
		return fmt.Errorf("error creating index: %w", err)
	}

	log.Printf("Created index %s with mapping.\n", index)
	return nil
}

func loadDataToElasticsearch(ctx context.Context, db *pgxpool.Pool, el *elastic.Client, indexName string) error {
	query := "SELECT id, name, description FROM some_table"

	rows, err := db.Query(ctx, query)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	bulkRequest := el.Bulk()
	count := 0

	for rows.Next() {
		// Считать строку значений
		var id, name, description string
		err := rows.Scan(&id, &name, &description)
		if err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		// Создаем JSON-документ для индексации
		doc := map[string]interface{}{
			"id":          id,
			"name":        name,
			"description": description,
		}

		// Create bulk request item
		req := elastic.NewBulkIndexRequest().Index(indexName).Id(id).Doc(doc)
		bulkRequest.Add(req)
		count++

		// Execute in batches
		if count%1000 == 0 {
			// Отправить собранные данные в Elasticsearch
			_, err := bulkRequest.Do(ctx)
			if err != nil {
				return fmt.Errorf("bulk indexing failed: %w", err)
			}
			log.Printf("Indexed %d documents...\n", count)
		}
	}

	// Execute remaining requests (если остались неотправленные документы)
	if bulkRequest.NumberOfActions() > 0 {
		_, err := bulkRequest.Do(ctx)
		if err != nil {
			return fmt.Errorf("final bulk indexing failed: %w", err)
		}
	}

	log.Printf("Successfully indexed %d documents.\n", count)
	return nil
}
