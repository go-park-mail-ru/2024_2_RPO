package user

import (
	"RPO_back/internal/models"
	"mime/multipart"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go
type UserUsecase interface {
	GetMyProfile(userID int) (profile *models.UserProfile, err error)
	UpdateMyProfile(userID int, data *models.UserProfileUpdate) (updatedProfile *models.UserProfile, err error)
	SetMyAvatar(userID int, file *multipart.File, fileHeader *multipart.FileHeader) (updated *models.UserProfile, err error)
}

type UserRepo interface {
	GetUserProfile(userID int) (profile *models.UserProfile, err error)
	UpdateUserProfile(userID int, data models.UserProfileUpdate) (newProfile *models.UserProfile, err error)
	SetUserAvatar(userID int, fileExtension string, fileSize int) (fileName string, err error)
}
