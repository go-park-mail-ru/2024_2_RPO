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

func TestGetMyProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockUser.NewMockUserRepo(ctrl)
	mockAuthClient := authGRPC.NewMockAuthClient(ctrl)
	uc := CreateUserUsecase(mockUserRepo, mockAuthClient)

	userID := int64(1)
	expectedProfile := &models.UserProfile{ID: userID, Name: "test_user"}
	mockUserRepo.EXPECT().GetUserProfile(context.Background(), userID).Return(expectedProfile, nil)

	profile, err := uc.GetMyProfile(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedProfile, profile)
}

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

func TestSetMyAvatar_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockUser.NewMockUserRepo(ctrl)
	mockAuthClient := authGRPC.NewMockAuthClient(ctrl)
	uc := CreateUserUsecase(mockUserRepo, mockAuthClient)

	userID := int64(1)
	uploadedFile := &models.UploadedFile{FileID: nil, Content: []byte{}, OriginalName: "avatar.png", UUID: nil, FileExtension: "png"}
	expectedProfile := &models.UserProfile{ID: userID, AvatarImageURL: "avatar.png"}

	mockUserRepo.EXPECT().DeduplicateFile(context.Background(), uploadedFile).Return(nil, nil, nil)
	mockUserRepo.EXPECT().RegisterFile(context.Background(), uploadedFile).Return(nil)
	mockUserRepo.EXPECT().GetUserProfile(context.Background(), userID).Return(expectedProfile, nil)

	profile, err := uc.SetMyAvatar(context.Background(), userID, uploadedFile)
	assert.NoError(t, err)
	assert.Equal(t, expectedProfile, profile)
}
