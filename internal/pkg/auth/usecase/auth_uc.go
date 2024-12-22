package usecase

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/pkg/auth"
	"RPO_back/internal/pkg/utils/encrypt"
	"context"
	"errors"
	"fmt"
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
	funcName := "CreateSession"
	passwordHash, err := uc.authRepo.GetUserPasswordHash(ctx, userID)
	if err != nil {
		return "", err
	}

	if passwordHash != nil {
		ok := encrypt.CheckPassword(password, *passwordHash)
		if !ok {
			return "", fmt.Errorf("%s: passwords not match: %w", funcName, errs.ErrWrongCredentials)
		}
	}

	sessionID = encrypt.GenerateSessionID()

	err = uc.authRepo.CreateSession(ctx, sessionID, userID)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

// CheckSession определяет, какому пользователю соответствует эта сессия.
// Если сессия не валидна, возвращает errs.ErrNotFound
func (uc *AuthUsecase) CheckSession(ctx context.Context, sessionID string) (userID int64, err error) {
	funcName := "CheckSession"
	userID, err = uc.authRepo.CheckSession(ctx, sessionID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return 0, errs.ErrNotFound
		}

		return 0, fmt.Errorf("%s: %w", funcName, err)
	}

	return userID, nil
}

func (uc *AuthUsecase) RemoveSession(ctx context.Context, sessionID string) (err error) {
	funcName := "RemoveSession"
	err = uc.authRepo.RemoveSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("%s: %w", funcName, err)
	}

	return nil
}

func (uc *AuthUsecase) ChangePassword(ctx context.Context, oldPassword string, newPassword string, sessionID string) (err error) {
	funcName := "ChangePassword"
	userID, err := uc.authRepo.CheckSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("%s (CheckSession): %w", funcName, err)
	}

	oldPasswordHash, err := uc.authRepo.GetUserPasswordHash(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s (GetUserPasswordHash): %w", funcName, err)
	}

	if oldPasswordHash != nil {
		ok := encrypt.CheckPassword(oldPassword, *oldPasswordHash)
		if !ok {
			return fmt.Errorf("%s (CheckPassword): passwords do not match: %w", funcName, errs.ErrWrongCredentials)
		}
	}

	err = uc.authRepo.DisplaceUserSessions(ctx, sessionID, userID)
	if err != nil {
		return fmt.Errorf("%s (DisplaceUserSessions): %w", funcName, err)
	}

	newPasswordHash, err := encrypt.SaltAndHashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("%s (SaltAndHashPassword): %w", funcName, err)
	}

	err = uc.authRepo.SetNewPasswordHash(ctx, userID, newPasswordHash)
	if err != nil {
		return fmt.Errorf("%s (SetNewPasswordHash): %w", funcName, err)
	}

	err = uc.authRepo.CreateSession(ctx, sessionID, userID)
	if err != nil {
		return fmt.Errorf("%s (CreateSession): %w", funcName, err)
	}

	return nil
}
