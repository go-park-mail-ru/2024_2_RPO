package delivery

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/pkg/auth/delivery/grpc/gen"
	mocks "RPO_back/internal/pkg/auth/mocks"
	"context"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateSession_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUsecase := mocks.NewMockAuthUsecase(ctrl)
	AuthDelivery := CreateAuthServer(mockAuthUsecase)

	ctx := context.Background()
	userID := int64(1)
	password := "password123"
	sessionID := "session123"

	mockAuthUsecase.EXPECT().
		CreateSession(ctx, userID, password).
		Return(sessionID, nil)

	request := &gen.UserDataRequest{
		UserID:   userID,
		Password: password,
	}

	resp, err := AuthDelivery.CreateSession(ctx, request)
	assert.NoError(t, err)
	assert.Equal(t, sessionID, resp.SessionID)
	assert.Equal(t, gen.Error_NONE, resp.Error)
}

func TestCreateSession_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUsecase := mocks.NewMockAuthUsecase(ctrl)
	AuthDelivery := CreateAuthServer(mockAuthUsecase)

	ctx := context.Background()
	userID := int64(1)
	password := "password123"

	mockAuthUsecase.EXPECT().
		CreateSession(ctx, userID, password).
		Return("", errs.ErrWrongCredentials)

	request := &gen.UserDataRequest{
		UserID:   userID,
		Password: password,
	}

	resp, err := AuthDelivery.CreateSession(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, gen.Error_INVALID_CREDENTIALS, resp.Error)
}

func TestCheckSession_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUsecase := mocks.NewMockAuthUsecase(ctrl)
	AuthDelivery := CreateAuthServer(mockAuthUsecase)

	ctx := context.Background()
	sessionID := "session123"
	userID := 1

	mockAuthUsecase.EXPECT().
		CheckSession(ctx, sessionID).
		Return(userID, nil)

	request := &gen.CheckSessionRequest{
		SessionID: sessionID,
	}

	resp, err := AuthDelivery.CheckSession(ctx, request)
	assert.NoError(t, err)
	assert.Equal(t, int64(userID), resp.UserID)
	assert.Equal(t, gen.Error_NONE, resp.Error)
}

func TestCheckSession_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUsecase := mocks.NewMockAuthUsecase(ctrl)
	AuthDelivery := CreateAuthServer(mockAuthUsecase)

	ctx := context.Background()
	sessionID := "session123"

	mockAuthUsecase.EXPECT().
		CheckSession(ctx, sessionID).
		Return(0, errs.ErrNotFound)

	request := &gen.CheckSessionRequest{
		SessionID: sessionID,
	}

	resp, err := AuthDelivery.CheckSession(ctx, request)
	assert.NoError(t, err)
	assert.Equal(t, gen.Error_INVALID_CREDENTIALS, resp.Error)
}

func TestDeleteSession_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUsecase := mocks.NewMockAuthUsecase(ctrl)
	AuthDelivery := CreateAuthServer(mockAuthUsecase)

	ctx := context.Background()
	sessionID := "session123"

	mockAuthUsecase.EXPECT().
		KillSession(ctx, sessionID).
		Return(nil)

	request := &gen.Session{
		SessionID: sessionID,
	}

	resp, err := AuthDelivery.DeleteSession(ctx, request)
	assert.NoError(t, err)
	assert.Equal(t, gen.Error_NONE, resp.Error)
}

func TestDeleteSession_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUsecase := mocks.NewMockAuthUsecase(ctrl)
	AuthDelivery := CreateAuthServer(mockAuthUsecase)

	ctx := context.Background()
	sessionID := "session123"

	mockAuthUsecase.EXPECT().
		KillSession(ctx, sessionID).
		Return(errors.New("unexpected error"))

	request := &gen.Session{
		SessionID: sessionID,
	}

	resp, err := AuthDelivery.DeleteSession(ctx, request)
	assert.NoError(t, err)
	assert.Equal(t, gen.Error_INTERNAL_SERVER_ERROR, resp.Error)
}

func TestChangePassword_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUsecase := mocks.NewMockAuthUsecase(ctrl)
	AuthDelivery := CreateAuthServer(mockAuthUsecase)

	ctx := context.Background()
	oldPassword := "oldpass"
	newPassword := "newpass"
	sessionID := "session123"

	mockAuthUsecase.EXPECT().
		ChangePassword(ctx, oldPassword, newPassword, sessionID).
		Return(nil)

	request := &gen.ChangePasswordRequest{
		PasswordOld: oldPassword,
		PasswordNew: newPassword,
		SessionID:   sessionID,
	}

	resp, err := AuthDelivery.ChangePassword(ctx, request)
	assert.NoError(t, err)
	assert.Equal(t, gen.Error_NONE, resp.Error)
}

func TestChangePassword_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUsecase := mocks.NewMockAuthUsecase(ctrl)
	AuthDelivery := CreateAuthServer(mockAuthUsecase)

	ctx := context.Background()
	oldPassword := "oldpass"
	newPassword := "newpass"
	sessionID := "session123"

	mockAuthUsecase.EXPECT().
		ChangePassword(ctx, oldPassword, newPassword, sessionID).
		Return(errors.New("unexpected error"))

	request := &gen.ChangePasswordRequest{
		PasswordOld: oldPassword,
		PasswordNew: newPassword,
		SessionID:   sessionID,
	}

	resp, err := AuthDelivery.ChangePassword(ctx, request)
	assert.NoError(t, err)
	assert.Equal(t, gen.Error_INTERNAL_SERVER_ERROR, resp.Error)
}
