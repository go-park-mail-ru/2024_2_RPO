package usecase

import (
	"RPO_back/internal/errs"
	mocks "RPO_back/internal/pkg/auth/mocks"
	"RPO_back/internal/pkg/utils/encrypt"
	"context"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateSession_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepo(ctrl)
	AuthUsecase := CreateAuthUsecase(mockRepo)

	ctx := context.Background()
	userID := int64(1)
	password := "password123"
	passwordHash, err := encrypt.SaltAndHashPassword(password)
	assert.NoError(t, err)

	mockRepo.EXPECT().GetUserPasswordHash(ctx, int(userID)).Return(&passwordHash, nil)
	mockRepo.EXPECT().RegisterSessionRedis(ctx, gomock.Any(), int(userID)).Return(nil).Times(1)

	result, err := AuthUsecase.CreateSession(ctx, userID, password)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestCreateSession_WrongPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepo(ctrl)
	AuthUsecase := CreateAuthUsecase(mockRepo)

	ctx := context.Background()
	userID := int64(1)
	password := "wrongPassword"
	passwordHash, err := encrypt.SaltAndHashPassword("correctPassword")
	assert.NoError(t, err)

	mockRepo.EXPECT().GetUserPasswordHash(ctx, int(userID)).Return(&passwordHash, nil)

	_, err = AuthUsecase.CreateSession(ctx, userID, password)
	assert.ErrorIs(t, err, errs.ErrWrongCredentials)
}

func TestCheckSession_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepo(ctrl)
	AuthUsecase := CreateAuthUsecase(mockRepo)

	ctx := context.Background()
	sessionID := "validSessionID"
	userID := 1

	mockRepo.EXPECT().CheckSession(ctx, sessionID).Return(userID, nil)

	result, err := AuthUsecase.CheckSession(ctx, sessionID)
	assert.NoError(t, err)
	assert.Equal(t, userID, result)
}

func TestCheckSession_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepo(ctrl)
	AuthUsecase := CreateAuthUsecase(mockRepo)

	ctx := context.Background()
	sessionID := "invalidSessionID"

	mockRepo.EXPECT().CheckSession(ctx, sessionID).Return(0, errs.ErrNotFound)

	_, err := AuthUsecase.CheckSession(ctx, sessionID)
	assert.ErrorIs(t, err, errs.ErrNotFound)
}

func TestKillSession_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepo(ctrl)
	AuthUsecase := CreateAuthUsecase(mockRepo)

	ctx := context.Background()
	sessionID := "validSessionID"

	mockRepo.EXPECT().KillSessionRedis(ctx, sessionID).Return(nil)

	err := AuthUsecase.KillSession(ctx, sessionID)
	assert.NoError(t, err)
}

func TestKillSession_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepo(ctrl)
	AuthUsecase := CreateAuthUsecase(mockRepo)

	ctx := context.Background()
	sessionID := "invalidSessionID"

	mockRepo.EXPECT().KillSessionRedis(ctx, sessionID).Return(errors.New("unexpected error"))

	err := AuthUsecase.KillSession(ctx, sessionID)
	assert.Error(t, err)
}

func TestChangePassword_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepo(ctrl)
	AuthUsecase := CreateAuthUsecase(mockRepo)

	ctx := context.Background()
	oldPassword := "oldPassword123"
	newPassword := "newPassword123"
	sessionID := "validSessionID"
	userID := 1
	oldPasswordHash, err := encrypt.SaltAndHashPassword(oldPassword)
	assert.NoError(t, err)
	_, err = encrypt.SaltAndHashPassword(newPassword)
	assert.NoError(t, err)

	mockRepo.EXPECT().CheckSession(ctx, sessionID).Return(userID, nil)
	mockRepo.EXPECT().GetUserPasswordHash(ctx, userID).Return(&oldPasswordHash, nil)
	mockRepo.EXPECT().DisplaceUserSessions(ctx, sessionID, int64(userID)).Return(nil)
	mockRepo.EXPECT().SetNewPasswordHash(ctx, userID, gomock.Any()).Return(nil)
	mockRepo.EXPECT().RegisterSessionRedis(ctx, sessionID, userID).Return(nil)

	err = AuthUsecase.ChangePassword(ctx, oldPassword, newPassword, sessionID)
	assert.NoError(t, err)
}

func TestChangePassword_WrongOldPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepo(ctrl)
	AuthUsecase := CreateAuthUsecase(mockRepo)

	ctx := context.Background()
	oldPassword := "wrongPassword"
	sessionID := "validSessionID"
	userID := 1
	oldPasswordHash, err := encrypt.SaltAndHashPassword("correctPassword")
	assert.NoError(t, err)

	mockRepo.EXPECT().CheckSession(ctx, sessionID).Return(userID, nil)
	mockRepo.EXPECT().GetUserPasswordHash(ctx, userID).Return(&oldPasswordHash, nil)

	err = AuthUsecase.ChangePassword(ctx, oldPassword, "newPassword123", sessionID)
	assert.ErrorIs(t, err, errs.ErrWrongCredentials)
}
