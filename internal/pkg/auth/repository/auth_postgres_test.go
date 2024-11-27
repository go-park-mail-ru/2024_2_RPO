package repository

import (
	mocks "RPO_back/internal/pkg/auth/mocks"
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestSetNewPasswordHash(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authRepo := mocks.NewMockAuthRepo(ctrl)

	ctx := context.Background()
	userID := int64(1)
	newPasswordHash := "new_password_hash"

	authRepo.EXPECT().SetNewPasswordHash(ctx, int(userID), newPasswordHash).Return(nil)

	err := authRepo.SetNewPasswordHash(ctx, int(userID), newPasswordHash)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestGetUserPasswordHash(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authRepo := mocks.NewMockAuthRepo(ctrl)

	ctx := context.Background()
	userID := int64(1)
	passwordHash := "hashed_password"

	authRepo.EXPECT().GetUserPasswordHash(ctx, int(userID)).Return(&passwordHash, nil)

	_, err := authRepo.GetUserPasswordHash(ctx, int(userID))
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
