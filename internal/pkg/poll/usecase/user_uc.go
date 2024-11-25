package usecase

import (
	"RPO_back/internal/models"
	authGRPC "RPO_back/internal/pkg/auth/delivery/grpc/gen"
	"RPO_back/internal/pkg/poll"
	"context"
	"fmt"
)

type PollUsecase struct {
	authClient authGRPC.AuthClient
	pollRepo   poll.PollRepo
}

func CreatePollUsecase(pollRepo poll.PollRepo, authClient authGRPC.AuthClient) *PollUsecase {
	return &PollUsecase{
		authClient: authClient,
		pollRepo:   pollRepo,
	}
}

func (uc *PollUsecase) SubmitPoll(ctx context.Context, userID int64, pollSubmit *models.PollSubmit) error {
	err := uc.pollRepo.SubmitPoll(ctx, userID, pollSubmit)
	if err != nil {
		return fmt.Errorf("SubmitPoll: %w", err)
	}

	return nil
}

func (uc *PollUsecase) GetPollResults(ctx context.Context) (pollResults *models.PollResults, err error) {
	pollRating, err := uc.pollRepo.GetRatingResults(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetPollResults (GetRatingResults): %w", err)
	}

	pollText, err := uc.pollRepo.GetTextResults(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetPollResults (GetTextResults): %w", err)
	}

	pollResults = &models.PollResults{
		RatingResults: pollRating,
		TextResults:   pollText,
	}

	return pollResults, nil
}
