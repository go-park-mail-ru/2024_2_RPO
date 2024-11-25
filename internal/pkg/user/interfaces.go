package user

import (
	"RPO_back/internal/models"
	"context"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go
type UserUsecase interface {
	GetMyProfile(ctx context.Context, userID int64) (profile *models.UserProfile, err error)
	UpdateMyProfile(ctx context.Context, userID int64, data *models.UserProfileUpdateRequest) (updatedProfile *models.UserProfile, err error)
	SetMyAvatar(ctx context.Context, userID int64, file *models.UploadedFile) (updated *models.UserProfile, err error)
	LoginUser(ctx context.Context, email string, password string) (sessionID string, err error)
	RegisterUser(ctx context.Context, user *models.UserRegisterRequest) (sessionID string, err error)
	LogoutUser(ctx context.Context, sessionID string) error
	ChangePassword(ctx context.Context, userID int64, oldPassword string, newPassword string) error
	SubmitPoll(ctx context.Context, userID int64, pollQuestion *models.PollSubmit) error
	GetPollResults(ctx context.Context) (pollResults *models.PollResults, err error)
}

type UserRepo interface {
	GetUserProfile(ctx context.Context, userID int64) (profile *models.UserProfile, err error)
	UpdateUserProfile(ctx context.Context, userID int64, data models.UserProfileUpdateRequest) (newProfile *models.UserProfile, err error)
	SetUserAvatar(ctx context.Context, userID int64, avatarFileID int64) (err error)
	GetUserByEmail(ctx context.Context, email string) (user *models.UserProfile, err error)
	CreateUser(ctx context.Context, user *models.UserRegisterRequest) (newUser *models.UserProfile, err error)
	CheckUniqueCredentials(ctx context.Context, nickname string, email string) error
	DeduplicateFile(ctx context.Context, file *models.UploadedFile) (fileNames []string, fileIDs []int64, err error)
	RegisterFile(ctx context.Context, file *models.UploadedFile) error
	SubmitPoll(ctx context.Context, userID int64, pollSubmit *models.PollSubmit) error
	GetRatingResults(ctx context.Context) (results []models.RatingResults, err error)
	GetTextResults(ctx context.Context) (results []models.AnswerResults, err error)
	SetNextPollDT(ctx context.Context, userID int64) error
	PickPollQuestions(ctx context.Context) (pollQuestions []models.PollQuestion, err error)
}
