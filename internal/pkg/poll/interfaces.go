package poll

import (
	"RPO_back/internal/models"
	"context"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go
type PollUsecase interface {
	SubmitPoll(ctx context.Context, userID int64, pollQuestion *models.PollSubmit) error
	GetPollResults(ctx context.Context) (pollResults *models.PollResults, err error)
	GetPollQuestions(ctx context.Context, userID int64) (pollQuestions []models.PollQuestion, err error)
}

type PollRepo interface {
	SubmitPoll(ctx context.Context, userID int64, pollSubmit *models.PollSubmit) error
	GetRatingResults(ctx context.Context) (results []models.RatingResults, err error)
	GetTextResults(ctx context.Context) (results []models.AnswerResults, err error)
	SetNextPollDT(ctx context.Context, userID int64) error
	GetPollQuestions(ctx context.Context) (pollQuestions []models.PollQuestion, err error)
}
