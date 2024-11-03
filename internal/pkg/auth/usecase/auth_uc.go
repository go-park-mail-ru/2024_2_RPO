package usecase

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/auth"
	"RPO_back/internal/pkg/utils/encrypt"
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

func (uc *AuthUsecase) LoginUser(email string, password string) (sessionId string, err error) {
	user, err := uc.authRepo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	ok := encrypt.CheckPassword(password, user.PasswordHash)
	if !ok {
		return "", fmt.Errorf("LoginUser: passwords not match: %w", errs.ErrWrongCredentials)
	}

	sessionID := encrypt.GenerateSessionID()
	err = uc.authRepo.RegisterSessionRedis(sessionID, user.ID)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (uc *AuthUsecase) RegisterUser(user *models.UserRegistration) (sessionId string, err error) {
	err = uc.authRepo.CheckUniqueCredentials(user.Name, user.Email)
	if err != nil {
		return "", err
	}

	hashedPassword, err := encrypt.SaltAndHashPassword(user.Password)
	if err != nil {
		return "", errors.New("failed to hash password")
	}

	log.Info("Password hash: ", hashedPassword)

	newUser, err := uc.authRepo.CreateUser(user, string(hashedPassword))
	if err != nil {
		return "", fmt.Errorf("internal error: %w", err)
	}

	sessionId = encrypt.GenerateSessionID()
	err = uc.authRepo.RegisterSessionRedis(sessionId, newUser.ID)
	if err != nil {
		return "", errors.New("failed to register session")
	}

	return sessionId, nil
}

func (uc *AuthUsecase) LogoutUser(sessionId string) error {
	return uc.authRepo.KillSessionRedis(sessionId)
}
