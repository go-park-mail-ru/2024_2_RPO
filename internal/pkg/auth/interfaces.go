package auth

import (
	"RPO_back/internal/models"
	"context"
)

const (
	SessionCookieName = "session_id"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type AuthUsecase interface {
	LoginUser(ctx context.Context, email string, password string) (sessionID string, err error)
	RegisterUser(ctx context.Context, user *models.UserRegistration) (sessionID string, err error)
	LogoutUser(ctx context.Context, sessionID string) error
}

type AuthRepo interface {
	RegisterSessionRedis(ctx context.Context, cookie string, userID int) error
	KillSessionRedis(ctx context.Context, sessionID string) error
	RetrieveUserIdFromSessionId(ctx context.Context, sessionId string) (userID int, err error)
	GetUserByEmail(ctx context.Context, email string) (user *models.UserProfile, err error)
	GetUserByID(ctx context.Context, userID int) (user *models.UserProfile, err error)
	CreateUser(ctx context.Context, user *models.UserRegistration, hashedPassword string) (newUser *models.UserProfile, err error)
	CheckUniqueCredentials(ctx context.Context, nickname string, email string) error
	SetNewPasswordHash(ctx context.Context, userID int, newPasswordHash string) error
}
