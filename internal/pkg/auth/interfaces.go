package auth

import (
	"RPO_back/internal/models"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go


type AuthUsecase interface {
	LoginUser(email string, password string) (sessionId string, err error)
	RegisterUser(user *models.UserRegistration) (sessionId string, err error)
	LogoutUser(sessionId string) error
}

type AuthRepo interface {
	RegisterSessionRedis(cookie string, userID int) error
	KillSessionRedis(sessionId string) error
	RetrieveUserIdFromSessionId(sessionId string) (userId int, err error)
	GetUserByEmail(email string) (user *models.UserProfile, err error)
	CreateUser(user *models.UserRegistration, hashedPassword string) (newUser *models.UserProfile, err error)
	CheckUniqueCredentials(nickname string, email string) error
	SetNewPasswordHash(userID int, newPasswordHash string)
}
