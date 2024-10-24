package usecase

import (
	"RPO_back/internal/pkg/auth/repository"
)

type AuthUsecase struct {
	authRepo *repository.AuthRepository
}
