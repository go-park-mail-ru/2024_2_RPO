package auth

import (
	"context"
)

const (
	SessionCookieName = "session_id"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type AuthUsecase interface {
	CreateSession(ctx context.Context, userID int64, password string) (sessionID string, err error)
	CheckSession(ctx context.Context, sessionID string) (userID int, err error)
	KillSession(ctx context.Context, sessionID string) (err error)
	ChangePassword(ctx context.Context, oldPassword string, newPassword string, sessionID string) (err error)
}

type AuthRepo interface {
	RegisterSessionRedis(ctx context.Context, cookie string, userID int) error
	KillSessionRedis(ctx context.Context, sessionID string) error
	CheckSession(ctx context.Context, sessionID string) (userID int, err error)
	SetNewPasswordHash(ctx context.Context, userID int, newPasswordHash string) error
	GetUserPasswordHashForUser(ctx context.Context, userID int) (passwordHash string, err error)
	DisplaceUserSessions(ctx context.Context, sessionID string, userID int64) error
}
