package usecase

import (
	"RPO_back/internal/pkg/user/repository"
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
func (uc *UserUsecase) GetMyProfile() {
	panic("Not implemented")
}

// UpdateMyProfile обновляет профиль пользователя и возвращает обновлённый профиль
func (uc *UserUsecase) UpdateMyProfile() {
	panic("Not implemented")
}
