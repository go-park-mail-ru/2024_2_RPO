package usecase

import (
	"RPO_back/internal/models"
	authGRPC "RPO_back/internal/pkg/auth/delivery/grpc/mocks"
	mockUser "RPO_back/internal/pkg/user/mocks"
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMyProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockUser.NewMockUserRepo(ctrl)
	mockAuthClient := authGRPC.NewMockAuthClient(ctrl)
	uc := CreateUserUsecase(mockUserRepo, mockAuthClient)

	userID := int64(1)
	updateData := &models.UserProfileUpdateRequest{NewName: "updated_nickname"}
	expectedProfile := &models.UserProfile{ID: userID, Name: "updated_nickname"}

	mockUserRepo.EXPECT().UpdateUserProfile(context.Background(), userID, *updateData).Return(expectedProfile, nil)

	updatedProfile, err := uc.UpdateMyProfile(context.Background(), userID, updateData)
	assert.NoError(t, err)
	assert.Equal(t, expectedProfile, updatedProfile)
}
