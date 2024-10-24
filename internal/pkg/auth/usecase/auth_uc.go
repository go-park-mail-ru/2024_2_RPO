package usecase

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/auth"
	"RPO_back/internal/pkg/auth/repository"
	"RPO_back/internal/pkg/utils/encrypt"
	"errors"
	"net/http"
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
	requestPasswordHash, err := encrypt.SaltAndHashPassword(password, []byte(user.PasswordSalt))
	if err != nil {
		return "", err
	}
	if requestPasswordHash != user.PasswordHash {
		return "", auth.ErrWrongCredentials
	}

	sessionID := encrypt.GenerateSessionID()
	err = this.authRepo.RegisterSessionRedis(sessionID, user.Id)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

func (this *AuthUsecase) RegisterUser(user *models.UserRegistration) (sessionId string, err error) {
	hashedPassword, err := encrypt.SaltAndHashPassword(user.Password)
	if err != nil {
		return "", errors.New("Failed to hash password")
	}

	newUser, err := this.authRepo.CreateUser(user, string(hashedPassword))
	if err != nil {
		return "", errors.New("Internal error")
	}

	sessionId, err = this.authRepo.RegisterSessionRedis(userID, newUser.Id)
	if err != nil {
		return "", errors.New("Failed to register session")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Session cookie is set"))

	return nil
}
