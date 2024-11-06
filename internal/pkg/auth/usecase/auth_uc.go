package usecase

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/auth"
	"RPO_back/internal/pkg/utils/encrypt"
	"context"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type AuthUsecase struct {
	authRepo auth.AuthRepo
}

func CreateAuthUsecase(repo auth.AuthRepo) *AuthUsecase {
	return &AuthUsecase{
		authRepo: repo,
	}
}

func (uc *AuthUsecase) LoginUser(ctx context.Context, email string, password string) (sessionID string, err error) {
	user, err := uc.authRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	ok := encrypt.CheckPassword(password, user.PasswordHash)
	if !ok {
		return "", fmt.Errorf("LoginUser: passwords not match: %w", errs.ErrWrongCredentials)
	}

	sessionID = encrypt.GenerateSessionID()
	err = uc.authRepo.RegisterSessionRedis(ctx, sessionID, user.ID)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (uc *AuthUsecase) RegisterUser(ctx context.Context, user *models.UserRegistration) (sessionID string, err error) {
	err = uc.authRepo.CheckUniqueCredentials(ctx, user.Name, user.Email)
	if err != nil {
		return "", err
	}

	hashedPassword, err := encrypt.SaltAndHashPassword(user.Password)
	if err != nil {
		return "", errors.New("failed to hash password")
	}

	log.Info("Password hash: ", hashedPassword)

	newUser, err := uc.authRepo.CreateUser(ctx, user, string(hashedPassword))
	if err != nil {
		return "", fmt.Errorf("internal error: %w", err)
	}

	sessionID = encrypt.GenerateSessionID()
	err = uc.authRepo.RegisterSessionRedis(ctx, sessionID, newUser.ID)
	if err != nil {
		return "", errors.New("failed to register session")
	}

	return sessionID, nil
}

func (uc *AuthUsecase) LogoutUser(ctx context.Context, sessionID string) error {
	return uc.authRepo.KillSessionRedis(ctx, sessionID)
}

func (uc *AuthUsecase) ChangePassword(ctx context.Context, userID int, oldPassword string, newPassword string) error {
	user, err := uc.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	ok := encrypt.CheckPassword(oldPassword, user.PasswordHash)
	if !ok {
		return fmt.Errorf("ChangePassword: %w", errs.ErrNotPermitted)
	}
	newPasswordHash, err := encrypt.SaltAndHashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("ChangePassword (hashing new): %w", err)
	}
	uc.authRepo.SetNewPasswordHash(ctx, userID, newPasswordHash)
	return nil
}
