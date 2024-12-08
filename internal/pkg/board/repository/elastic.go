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

func (be *BoardElasticRepository) PutCard(ctx context.Context, boardID int64, cardID int64, cardText string) error {
	funcName := "PutCard"
	if len(cardText) < ElasticSearchValueMinLength {
		return fmt.Errorf("%s: cardText must be at least %d characters", funcName, ElasticSearchValueMinLength)
	}

	doc := models.ElasticCard{
		BoardID:  boardID,
		CardID:   cardID,
		CardText: cardText,
	}

	docID := fmt.Sprintf("%d:%d", boardID, cardID)

	_, err := be.elastic.Index().Index(ElasticIdxName).Id(docID).BodyJson(doc).Do(ctx)
	if err != nil {
		logging.Warn(ctx, fmt.Sprintf("%s: failed to index card: %v", funcName, err))
		return fmt.Errorf("%s: index operation failed", funcName)
	}

	return nil
}

func (be *BoardElasticRepository) Search(ctx context.Context, query string) ([]int64, error) {
	funcName := "Search"
	if len(query) < ElasticSearchValueMinLength {
		return nil, fmt.Errorf("%s: query must be at least %d characters", funcName, ElasticSearchValueMinLength)
	}

	searchQuery := elastic.NewMatchQuery("cardText", query)
	searchResult, err := be.elastic.Search().Index(ElasticIdxName).Query(searchQuery).Do(ctx)
	if err != nil {
		logging.Error(ctx, fmt.Sprintf("%s: search error: %v", funcName, err))
		return nil, fmt.Errorf("%s: search operation failed", funcName)
	}

	var cardIDs []int64
	for _, hit := range searchResult.Hits.Hits {
		var doc models.ElasticCard
		if err := json.Unmarshal(hit.Source, &doc); err != nil {
			logging.Warn(ctx, fmt.Sprintf("%s: failed to unmarshal hit: %v", funcName, err))
			continue
		}
		cardIDs = append(cardIDs, doc.CardID)
	}

	return cardIDs, nil
}

func (be *BoardElasticRepository) DeleteCard(ctx context.Context, boardID, cardID int64) (err error) {
	funcName := "DeleteCard"
	docID := fmt.Sprintf("%d:%d", boardID, cardID)
	_, err = be.elastic.Delete().Index(ElasticIdxName).Id(docID).Do(ctx)
	logging.Debug(ctx, funcName, "delete has error: ", err)
	if err != nil {
		return fmt.Errorf("%s (delete)", funcName)
	}

	return nil
}
