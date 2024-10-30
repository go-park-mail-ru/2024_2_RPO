package usecase

import "RPO_back/internal/pkg/user/repository"

type UserUsecase struct {
	userRepo *repository.UserRepository
}

func CreateUserUsecase(userRepo *repository.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}
