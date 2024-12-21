package repository

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/utils/logging"
	"context"
	"encoding/json"
	"fmt"

	"github.com/olivere/elastic/v7"
)

const ElasticIdxName = "board_index"
const ElasticSearchValueMinLength = 3 // минимальная длина запроса для поиска

type BoardElasticRepository struct {
	elastic *elastic.Client
}

func CreateBoardElasticRepository(el *elastic.Client) *BoardElasticRepository {
	return &BoardElasticRepository{
		elastic: el,
	}
}

func (be *BoardElasticRepository) CreateCard(ctx context.Context, boardID int64, cardID int64, cardTitle string) error {
	funcName := "CreateCard"

	docID := fmt.Sprintf("%d", cardID)

	_, err := be.elastic.Index().
		Index(ElasticIdxName).
		Id(docID).
		BodyJson(map[string]interface{}{
			"card_id":  cardID,
			"title":    cardTitle,
			"board_id": boardID,
		}).
		Refresh("wait_for").
		Do(ctx)
	logging.Debug(ctx, funcName, "CreateCard has error: ", err)
	if err != nil {
		return fmt.Errorf("%s create has error", funcName)
	}

	return nil
}

func (be *BoardElasticRepository) UpdateCard(ctx context.Context, boardID int64, cardID int64, cardTitle string) error {
	funcName := "UpdateCard"

	docID := fmt.Sprintf("%d", cardID)

	_, err := be.elastic.Update().
		Index(ElasticIdxName).
		Id(docID).
		Doc(map[string]interface{}{
			"title":    cardTitle,
			"board_id": boardID,
		}).
		Refresh("wait_for").
		Do(ctx)
	logging.Debug(ctx, funcName, "UpdateCard has error: ", err)
	if err != nil {
		return fmt.Errorf("%s update has error", funcName)
	}

	return nil
}

func (be *BoardElasticRepository) Search(ctx context.Context, boards []models.Board, searchValue string) (foundCards []int64, err error) {
	funcName := "Search"
	if len(searchValue) < ElasticSearchValueMinLength {
		return nil, fmt.Errorf("%s: query must be at least %d characters", funcName, ElasticSearchValueMinLength)
	}

	boardIDs := make([]interface{}, len(boards))
	for i, board := range boards {
		boardIDs[i] = board.ID
	}

	// boardQuery := elastic.NewTermsQuery("board_id", boardIDs...)
	searchQuery := elastic.NewMatchQuery("title", searchValue).Fuzziness("AUTO")

	// fullQuery := elastic.NewBoolQuery().Filter(boardQuery).Must(searchQuery)
	fullQuery := elastic.NewBoolQuery().Must(searchQuery)

	searchResult, err := be.elastic.Search().
		Index(ElasticIdxName).
		Query(fullQuery).
		Do(ctx)
	logging.Debug(ctx, funcName, "error performing search: ", err)
	if err != nil {
		return nil, fmt.Errorf("%s: search operation failed", funcName)
	}

	foundCards = make([]int64, 0, len(searchResult.Hits.Hits))

	for _, hit := range searchResult.Hits.Hits {
		var card struct {
			CardID int64 `json:"card_id"`
		}
		if err := json.Unmarshal(hit.Source, &card); err != nil {
			logging.Debug(ctx, funcName, "failed to unmarshal card ID ", err)
			continue
		}
		foundCards = append(foundCards, card.CardID)
	}

	fmt.Println("===FOUND CARDS HERE===")
	fmt.Printf("%#v\n", foundCards)
	fmt.Println("===FOUND CARDS HERE===")

	return foundCards, nil
}

func (be *BoardElasticRepository) DeleteCard(ctx context.Context, cardID int64) (err error) {
	funcName := "DeleteCard"
	docID := fmt.Sprintf("%d", cardID)
	_, err = be.elastic.Delete().
		Index(ElasticIdxName).
		Id(docID).
		Do(ctx)
	logging.Debug(ctx, funcName, "delete has error: ", err)
	if err != nil {
		return fmt.Errorf("%s (delete)", funcName)
	}

	return nil
}
