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
	query := `SELECT c.card_id, c.title, b.board_id
	FROM card AS c
	JOIN kanban_column AS kc ON c.col_id = kc.col_id
	JOIN board AS b ON kc.board_id = b.board_id;`

	rows, err := db.Query(ctx, query)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	bulkRequest := el.Bulk()
	count := 0

	for rows.Next() {
		var cardID, cardTitle, boardID string
		err := rows.Scan(&cardID, &cardTitle, &boardID)
		if err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		doc := map[string]interface{}{
			"card_id":  cardID,
			"title":    cardTitle,
			"board_id": boardID,
		}

		_, err = el.Index().Index(indexName).Id(cardID).BodyJson(doc).Refresh("wait_for").Do(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		count++
		if count%100 == 0 {
			fmt.Printf("Indexed %d docs\n", count)
		}
	}

	if bulkRequest.NumberOfActions() > 0 {
		_, err := bulkRequest.Do(ctx)
		if err != nil {
			return fmt.Errorf("final bulk indexing failed: %w", err)
		}
	}

	log.Printf("Successfully indexed %d documents.\n", count)
	return nil
}
