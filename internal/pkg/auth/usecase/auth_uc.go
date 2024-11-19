package usecase

import (
	"RPO_back/internal/pkg/auth"
	"context"
)

type AuthUsecase struct {
	authRepo auth.AuthRepo
}

func CreateAuthUsecase(repo auth.AuthRepo) *AuthUsecase {
	return &AuthUsecase{
		authRepo: repo,
	}
}

func (uc *AuthUsecase) CreateSession(ctx context.Context, userID int64, password string) (sessionID string, err error) {

}

func (uc *AuthUsecase) CheckSession(ctx context.Context, sessionID string) (userID int64, err error) {

}

func (uc *AuthUsecase) KillSession(ctx context.Context, sessionID string) (err error) {

}

func (uc *AuthUsecase) ChangePassword(ctx context.Context, oldPassword string, newPassword string) (err error) {

}
