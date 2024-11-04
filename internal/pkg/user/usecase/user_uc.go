package usecase

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/user/repository"
	"mime/multipart"
)

type UserUsecase struct {
	userRepo *repository.UserRepository
}

func CreateUserUsecase(userRepo *repository.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

// GetMyProfile возвращает пользователю его профиль
func (uc *UserUsecase) GetMyProfile(userID int) (profile *models.UserProfile, err error) {
	profile, err = uc.userRepo.GetUserProfile(userID)
	return
}

// UpdateMyProfile обновляет профиль пользователя и возвращает обновлённый профиль
func (uc *UserUsecase) UpdateMyProfile(userID int, data *models.UserProfileUpdate) (updatedProfile *models.UserProfile, err error) {
	updatedProfile, err = uc.userRepo.UpdateUserProfile(userID, *data)
	return
}

// SetMyAvatar устанавливает пользователю аватарку
func (uc *UserUsecase) SetMyAvatar(userID int, file *multipart.File) (updated *models.UserProfile, err error) {
	uc.userRepo.SetUserAvatar()
}
