package repository

import (
	mocks "RPO_back/internal/pkg/auth/mocks"
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestRegisterSessionRedis(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authRepo := mocks.NewMockAuthRepo(ctrl)

	ctx := context.Background()
	sessionID := "session_1"
	userID := 1

	authRepo.EXPECT().RegisterSessionRedis(ctx, sessionID, userID).Return(nil)

	err := authRepo.RegisterSessionRedis(ctx, sessionID, userID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestKillSessionRedis(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authRepo := mocks.NewMockAuthRepo(ctrl)

	ctx := context.Background()
	sessionID := "session_1"

	authRepo.EXPECT().KillSessionRedis(ctx, sessionID).Return(nil)

	err := authRepo.KillSessionRedis(ctx, sessionID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestDisplaceUserSessions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authRepo := mocks.NewMockAuthRepo(ctrl)

	ctx := context.Background()
	sessionID := "session_1"
	userID := int64(1)

	authRepo.EXPECT().DisplaceUserSessions(ctx, sessionID, userID).Return(nil)

	err := authRepo.DisplaceUserSessions(ctx, sessionID, userID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestCheckSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authRepo := mocks.NewMockAuthRepo(ctrl)

	ctx := context.Background()
	sessionID := "session_1"
	userID := 1

	authRepo.EXPECT().CheckSession(ctx, sessionID).Return(userID, nil)

	result, err := authRepo.CheckSession(ctx, sessionID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result != userID {
		t.Errorf("expected userID %v, got %v", userID, result)
	}
}
