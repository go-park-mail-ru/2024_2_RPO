package user

import (
	"RPO_back/internal/models"
	"context"
	"mime/multipart"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go
type UserUsecase interface {
	GetMyProfile(ctx context.Context, userID int) (profile *models.UserProfile, err error)
	UpdateMyProfile(ctx context.Context, userID int, data *models.UserProfileUpdate) (updatedProfile *models.UserProfile, err error)
	SetMyAvatar(ctx context.Context, userID int, file *multipart.File, fileHeader *multipart.FileHeader) (updated *models.UserProfile, err error)
}

type UserRepo interface {
	GetUserProfile(ctx context.Context, userID int) (profile *models.UserProfile, err error)
	UpdateUserProfile(ctx context.Context, userID int, data models.UserProfileUpdate) (newProfile *models.UserProfile, err error)
	SetUserAvatar(ctx context.Context, userID int, fileExtension string, fileSize int) (fileName string, err error)
}
