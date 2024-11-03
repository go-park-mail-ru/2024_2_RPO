package user

import (
	"RPO_back/internal/models"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go


type BoardUsecase interface {
	GetMyProfile(userID int) (profile *models.UserProfile, err error)
	UpdateMyProfile(userID int, data *models.UserProfileUpdate) (updatedProfile *models.UserProfile, err error)
}

type BoardRepo interface {
	GetUserProfile(userID int) (profile *models.UserProfile, err error)
	UpdateUserProfile(userID int, data models.UserProfileUpdate) (newProfile *models.UserProfile, err error)
}
