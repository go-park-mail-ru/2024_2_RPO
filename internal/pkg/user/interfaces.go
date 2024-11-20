package user

import (
	"RPO_back/internal/models"
	"context"
	"mime/multipart"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go
type UserUsecase interface {
	GetMyProfile(ctx context.Context, userID int) (profile *models.UserProfile, err error)
	UpdateMyProfile(ctx context.Context, userID int, data *models.UserProfileUpdateRequest) (updatedProfile *models.UserProfile, err error)
	SetMyAvatar(ctx context.Context, userID int, file *multipart.File, fileHeader *multipart.FileHeader) (updated *models.UserProfile, err error)
	LoginUser(ctx context.Context, email string, password string) (sessionID string, err error)
	RegisterUser(ctx context.Context, user *models.UserRegisterRequest) (sessionID string, err error)
	LogoutUser(ctx context.Context, sessionID string) error
	ChangePassword(ctx context.Context, userID int, oldPassword string, newPassword string) error
}

type UserRepo interface {
	GetUserProfile(ctx context.Context, userID int) (profile *models.UserProfile, err error)
	UpdateUserProfile(ctx context.Context, userID int, data models.UserProfileUpdateRequest) (newProfile *models.UserProfile, err error)
	SetUserAvatar(ctx context.Context, userID int, fileExtension string, fileSize int) (fileName string, err error)
	GetUserByEmail(ctx context.Context, email string) (user *models.UserProfile, err error)
	CreateUser(ctx context.Context, user *models.UserRegisterRequest) (newUser *models.UserProfile, err error)
	CheckUniqueCredentials(ctx context.Context, nickname string, email string) error
}
