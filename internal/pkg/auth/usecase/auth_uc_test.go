package usecase_test

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	mocks "RPO_back/internal/pkg/auth/mocks"
	"RPO_back/internal/pkg/utils/encrypt"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthUsecase_LoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthRepo := mocks.NewMockAuthRepo(ctrl)
	authUsecase := mocks.NewMockAuthUsecase(ctrl)

	tests := []struct {
		name           string
		email          string
		password       string
		setupMock      func()
		expectedError  error
		expectedResult string
	}{
		{
			name:    "successful login",
			email:   "test@example.com",
			password: "correctpassword",
			setupMock: func() {
				mockAuthRepo.EXPECT().GetUserByEmail("test@example.com").Return(&models.UserProfile{
					ID:          123,
					PasswordHash: "dpksdkfposkfo1341",
				}, nil)

				mockAuthRepo.EXPECT().RegisterSessionRedis(gomock.Any(), "1").Return(nil)
			},
			expectedError:  nil,
			expectedResult: "valid-session-id",
		},
		{
			name:    "user not found",
			email:   "notfound@example.com",
			password: "any-password",
			setupMock: func() {
				mockAuthRepo.EXPECT().GetUserByEmail("notfound@example.com").Return(nil, fmt.Errorf("user not found"))
			},
			expectedError:  fmt.Errorf("user not found"),
			expectedResult: "",
		},
		{
			name:    "wrong password",
			email:   "test@example.com",
			password: "wrongpassword",
			setupMock: func() {
				mockAuthRepo.EXPECT().GetUserByEmail("test@example.com").Return(&models.UserProfile{
					ID:          123,
					PasswordHash: "dpksdkfposkfo1341",
				}, nil)
			},
			expectedError:  fmt.Errorf("LoginUser: passwords not match: %w", errs.ErrWrongCredentials),
			expectedResult: "",
		},
		{
			name:    "redis session registration fails",
			email:   "test@example.com",
			password: "correctpassword",
			setupMock: func() {
				mockAuthRepo.EXPECT().GetUserByEmail("test@example.com").Return(&models.UserProfile{
					ID:          123,
					PasswordHash: "dpksdkfposkfo1341",
				}, nil)
				mockAuthRepo.EXPECT().RegisterSessionRedis(gomock.Any(), "1").Return(errors.New("redis error"))
			},
			expectedError:  errors.New("redis error"),
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypt.GenerateSessionID() // Mock session ID generation

			tt.setupMock()

			result, err := authUsecase.LoginUser(tt.email, tt.password)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
