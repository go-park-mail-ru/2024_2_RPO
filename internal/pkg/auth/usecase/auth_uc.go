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
	_, err = uc.authRepo.GetUserPasswordHash(ctx, int(userID))
	if err != nil {
		return "", err
	}

	// if passwordHash != nil {
	// 	ok := encrypt.CheckPassword(password, *passwordHash)
	// 	if !ok {
	// 		return "", fmt.Errorf("LoginUser: passwords not match: %w", errs.ErrWrongCredentials)
	// 	}
	// }

	sessionID = encrypt.GenerateSessionID()

	err = uc.authRepo.RegisterSessionRedis(ctx, sessionID, int(userID))
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (uc *AuthUsecase) CheckSession(ctx context.Context, sessionID string) (userID int, err error) {
	userID, err = uc.authRepo.CheckSession(ctx, sessionID)
	fmt.Println("CHECK SESSION => u_id=", userID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return 0, errs.ErrNotFound
		}

		return 0, fmt.Errorf("CheckSession: %w", err)
	}

	return userID, nil
}

func (uc *AuthUsecase) KillSession(ctx context.Context, sessionID string) (err error) {
	err = uc.authRepo.KillSessionRedis(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("KillSession: %w", err)
	}

	return nil
}

func (uc *AuthUsecase) ChangePassword(ctx context.Context, oldPassword string, newPassword string, sessionID string) (err error) {
	userID, err := uc.authRepo.CheckSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("ChangePassword (CheckSession): %w", err)
	}

	_, err = uc.authRepo.GetUserPasswordHash(ctx, int(userID))
	if err != nil {
		return fmt.Errorf("ChangePassword (GetUserPasswordHash): %w", err)
	}

	// if oldPasswordHash != nil {
	// 	ok := encrypt.CheckPassword(oldPassword, *oldPasswordHash)
	// 	if !ok {
	// 		return fmt.Errorf("ChangePassword (CheckPassword): passwords do not match: %w", errs.ErrWrongCredentials)
	// 	}
	// }

	err = uc.authRepo.DisplaceUserSessions(ctx, sessionID, int64(userID))
	if err != nil {
		return fmt.Errorf("ChangePassword (DisplaceUserSessions): %w", err)
	}

	newPasswordHash, err := encrypt.SaltAndHashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("ChangePassword (SaltAndHashPassword): %w", err)
	}

	err = uc.authRepo.SetNewPasswordHash(ctx, int(userID), newPasswordHash)
	if err != nil {
		return fmt.Errorf("ChangePassword (SetNewPasswordHash): %w", err)
	}

	err = uc.authRepo.RegisterSessionRedis(ctx, sessionID, int(userID))
	if err != nil {
		return fmt.Errorf("ChangePassword (RegisterSessionRedis): %w", err)
	}

	return nil
}
