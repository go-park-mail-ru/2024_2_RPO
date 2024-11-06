package usecase_test

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	mocks "RPO_back/internal/pkg/auth/mocks"
	"RPO_back/internal/pkg/utils/encrypt"
	"context"
	"errors"
	"fmt"
	"testing"

	AuthUsecase "RPO_back/internal/pkg/auth/usecase"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func isStringLengthMoreThan10(s string) bool {
	return len(s) > 10
}

func TestAuthUsecase_LoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthRepo := mocks.NewMockAuthRepo(ctrl)
	authUsecase := AuthUsecase.CreateAuthUsecase(mockAuthRepo)

	tests := []struct {
		name                  string
		email                 string
		password              string
		setupMock             func()
		expectedError         bool
		expectedResultChecker func(string) bool
	}{
		{
			name:     "successful login",
			email:    "test@example.com",
			password: "11111111",
			setupMock: func() {
				pHash, _ := encrypt.SaltAndHashPassword("11111111")
				mockAuthRepo.EXPECT().GetUserByEmail(gomock.Any(), "test@example.com").Return(&models.UserProfile{
					ID:           123,
					PasswordHash: pHash,
				}, nil)

				mockAuthRepo.EXPECT().RegisterSessionRedis(gomock.Any(), gomock.Any(), 123).Return(nil)
			},
			expectedError:         false,
			expectedResultChecker: isStringLengthMoreThan10,
		},
		{
			name:     "user not found",
			email:    "notfound@example.com",
			password: "any-password",
			setupMock: func() {
				mockAuthRepo.EXPECT().GetUserByEmail(gomock.Any(), "notfound@example.com").Return(nil, errs.ErrWrongCredentials)
			},
			expectedError:         true,
			expectedResultChecker: func(value string) bool { return value == "" },
		},
		{
			name:     "wrong password",
			email:    "test@example.com",
			password: "wrongpassword",
			setupMock: func() {
				hash, _ := encrypt.SaltAndHashPassword("11111111")
				mockAuthRepo.EXPECT().GetUserByEmail(gomock.Any(), "test@example.com").Return(&models.UserProfile{
					ID:           123,
					PasswordHash: hash,
				}, nil)
			},
			expectedError:         true,
			expectedResultChecker: func(value string) bool { return value == "" },
		},
		{
			name:     "redis session registration fails",
			email:    "test@example.com",
			password: "11111111",
			setupMock: func() {
				hash, _ := encrypt.SaltAndHashPassword("11111111")
				mockAuthRepo.EXPECT().GetUserByEmail(gomock.Any(), "test@example.com").Return(&models.UserProfile{
					ID:           123,
					PasswordHash: hash,
				}, nil)
				mockAuthRepo.EXPECT().RegisterSessionRedis(gomock.Any(), gomock.Any(), 123).Return(errors.New("redis error"))
			},
			expectedError:         true,
			expectedResultChecker: func(value string) bool { return value == "" },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypt.GenerateSessionID() // Mock session ID generation

			tt.setupMock()

			result, err := authUsecase.LoginUser(context.Background(), tt.email, tt.password)
			assert.Equal(t, err != nil, tt.expectedError)
			assert.Equal(t, tt.expectedResultChecker(result), true)
		})
	}
}

func TestAuthUsecase_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthRepo := mocks.NewMockAuthRepo(ctrl)
	authUsecase := AuthUsecase.CreateAuthUsecase(mockAuthRepo)

	tests := []struct {
		name                  string
		user                  *models.UserRegistration
		setupMock             func()
		expectedError         bool
		expectedResultChecker func(string) bool
	}{
		{
			name: "successful registration",
			user: &models.UserRegistration{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "11111111",
			},
			setupMock: func() {
				mockAuthRepo.EXPECT().CheckUniqueCredentials(gomock.Any(), "Test User", "test@example.com").Return(nil)
				mockAuthRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(&models.UserProfile{ID: 1}, nil)
				mockAuthRepo.EXPECT().RegisterSessionRedis(gomock.Any(), gomock.Any(), 1).Return(nil)
			},
			expectedError:         false,
			expectedResultChecker: isStringLengthMoreThan10,
		},
		{
			name: "non-unique credentials",
			user: &models.UserRegistration{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "11111111",
			},
			setupMock: func() {
				mockAuthRepo.EXPECT().CheckUniqueCredentials(gomock.Any(), "Test User", "test@example.com").Return(errors.New("credentials not unique"))
			},
			expectedError:         true,
			expectedResultChecker: func(value string) bool { return value == "" },
		},
		{
			name: "failed to create user",
			user: &models.UserRegistration{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "11111111",
			},
			setupMock: func() {
				mockAuthRepo.EXPECT().CheckUniqueCredentials(gomock.Any(), "Test User", "test@example.com").Return(nil)
				mockAuthRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("create user error"))
			},
			expectedError:         true,
			expectedResultChecker: func(value string) bool { return value == "" },
		},
		{
			name: "failed to register session",
			user: &models.UserRegistration{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "11111111",
			},
			setupMock: func() {
				mockAuthRepo.EXPECT().CheckUniqueCredentials(gomock.Any(), "Test User", "test@example.com").Return(nil)
				mockAuthRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(&models.UserProfile{ID: 1}, nil)
				mockAuthRepo.EXPECT().RegisterSessionRedis(gomock.Any(), gomock.Any(), 1).Return(errors.New("session registration error"))
			},
			expectedError:         true,
			expectedResultChecker: func(value string) bool { return value == "" },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypt.GenerateSessionID()

			tt.setupMock()

			result, err := authUsecase.RegisterUser(context.Background(), tt.user)
			assert.Equal(t, err != nil, tt.expectedError)
			assert.Equal(t, tt.expectedResultChecker(result), true)
		})
	}
}

func TestAuthUsecase_LogoutUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthRepo := mocks.NewMockAuthRepo(ctrl)
	authUsecase := AuthUsecase.CreateAuthUsecase(mockAuthRepo)

	tests := []struct {
		name          string
		sessionID     string
		setupMock     func()
		expectedError bool
	}{
		{
			name:      "successful logout",
			sessionID: "valid-session-id",
			setupMock: func() {
				mockAuthRepo.EXPECT().KillSessionRedis(gomock.Any(), "valid-session-id").Return(nil)
			},
			expectedError: false,
		},
		{
			name:      "session not found",
			sessionID: "invalid-session-id",
			setupMock: func() {
				mockAuthRepo.EXPECT().KillSessionRedis(gomock.Any(), "invalid-session-id").Return(fmt.Errorf("Error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := authUsecase.LogoutUser(context.Background(), tt.sessionID)
			assert.Equal(t, err != nil, tt.expectedError)
		})
	}
}
