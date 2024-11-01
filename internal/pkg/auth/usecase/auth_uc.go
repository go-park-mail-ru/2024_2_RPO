package usecase

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/auth"
	"RPO_back/internal/pkg/auth/repository"
	"RPO_back/internal/pkg/utils/encrypt"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type AuthUsecase struct {
	authRepo *repository.AuthRepository
}

func CreateAuthUsecase(repo *repository.AuthRepository) *AuthUsecase {
	return &AuthUsecase{
		authRepo: repo,
	}
}

func (this *AuthUsecase) LoginUser(email string, password string) (sessionId string, err error) {
	user, err := this.authRepo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	ok := encrypt.CheckPassword(password, user.PasswordHash)
	if !ok {
		return "", fmt.Errorf("LoginUser: passwords not match: %w", auth.ErrWrongCredentials)
	}

	sessionID := encrypt.GenerateSessionID()
	err = this.authRepo.RegisterSessionRedis(sessionID, user.ID)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

func (this *AuthUsecase) RegisterUser(user *models.UserRegistration) (sessionId string, err error) {
	err = this.authRepo.CheckUniqueCredentials(user.Name, user.Email)
	if err != nil {
		return "", err
	}
	hashedPassword, err := encrypt.SaltAndHashPassword(user.Password)
	if err != nil {
		return "", errors.New("Failed to hash password")
	}
	log.Info("Password hash: ", hashedPassword)

	newUser, err := this.authRepo.CreateUser(user, string(hashedPassword))
	if err != nil {
		return "", fmt.Errorf("Internal error: %w", err)
	}

	sessionId = encrypt.GenerateSessionID()
	err = this.authRepo.RegisterSessionRedis(sessionId, newUser.ID)
	if err != nil {
		return "", errors.New("Failed to register session")
	}

	return sessionId, nil
}

func (this *AuthUsecase) LogoutUser(sessionId string) error {
	return this.authRepo.KillSessionRedis(sessionId)
}
