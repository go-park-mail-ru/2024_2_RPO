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
	CreateSession(ctx context.Context, userID int64, password string) (sessionID string, err error)
	CheckSession(ctx context.Context, sessionID string) (userID int64, err error)
	KillSession(ctx context.Context, sessionID string) (err error)
	ChangePassword(ctx context.Context, oldPassword string, newPassword string) (err error)
}

type AuthRepo interface {
	RegisterSessionRedis(ctx context.Context, cookie string, userID int) error
	KillSessionRedis(ctx context.Context, sessionID string) error
	RetrieveUserIDFromSession(ctx context.Context, sessionID string) (userID int, err error)
	GetUserByEmail(ctx context.Context, email string) (user *models.UserProfile, err error)
	GetUserByID(ctx context.Context, userID int) (user *models.UserProfile, err error)
	CreateUser(ctx context.Context, user *models.UserRegisterRequest, hashedPassword string) (newUser *models.UserProfile, err error)
	CheckUniqueCredentials(ctx context.Context, nickname string, email string) error
	SetNewPasswordHash(ctx context.Context, userID int, newPasswordHash string) error
}
