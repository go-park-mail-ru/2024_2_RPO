package repository

import (
	"context"

	"github.com/olivere/elastic/v7"
)

const ElasticIdxName = ""
const ElasticSearchValueMinLength = 0

type BoardElasticRepository struct {
	elastic *elastic.Client
}

func CreateBoardElasticRepository(el *elastic.Client) *BoardElasticRepository {
	return &BoardElasticRepository{
		elastic: el,
	}
}

func (be *BoardElasticRepository) PutCard(ctx context.Context, boardID int64, cardID int64, cardText string) (err error) {
	panic("Not implemented")
}

func (be *BoardElasticRepository) Search(ctx context.Context, query string) (cardID []int64, err error) {
	panic("Not implemented")
}

func (be *BoardElasticRepository) DeleteCard(ctx context.Context, cardID int64) (err error) {
	panic("Not implemented")
}
