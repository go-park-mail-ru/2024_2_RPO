package usecase_test

import (
	"RPO_back/internal/models"
	mocks "RPO_back/internal/pkg/user/mocks"
	"RPO_back/internal/pkg/user/usecase"
	"errors"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserUsecase_GetMyProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	userUsecase := usecase.CreateUserUsecase(mockUserRepo)

	t.Run("successful profile retrieval", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserProfile(1).Return(&models.UserProfile{
			ID:             1,
			Name:           "Test User",
			Email:          "testuser@example.com",
			Description:    "User Description",
			JoinedAt:       time.Now(),
			UpdatedAt:      time.Now(),
			AvatarImageURL: "",
		}, nil)

		profile, err := userUsecase.GetMyProfile(1)
		assert.NoError(t, err)
		assert.Equal(t, "Test User", profile.Name)
		assert.Equal(t, "testuser@example.com", profile.Email)
	})

	t.Run("failed profile retrieval", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserProfile(1).Return(nil, errors.New("user not found"))

		profile, err := userUsecase.GetMyProfile(1)
		assert.Error(t, err)
		assert.Nil(t, profile)
	})
}

func TestUserUsecase_UpdateMyProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	userUsecase := usecase.CreateUserUsecase(mockUserRepo)

	t.Run("successful profile update", func(t *testing.T) {
		updateData := &models.UserProfileUpdate{
			NewName: "Updated User",
			Email:   "updateduser@example.com",
		}
		mockUserRepo.EXPECT().UpdateUserProfile(1, *updateData).Return(&models.UserProfile{
			ID:             1,
			Name:           "Updated User",
			Email:          "updateduser@example.com",
			Description:    "Updated Description",
			JoinedAt:       time.Now(),
			UpdatedAt:      time.Now(),
			AvatarImageURL: "",
		}, nil)

		profile, err := userUsecase.UpdateMyProfile(1, updateData)
		assert.NoError(t, err)
		assert.Equal(t, "Updated User", profile.Name)
		assert.Equal(t, "updateduser@example.com", profile.Email)
	})

	t.Run("failed profile update", func(t *testing.T) {
		updateData := &models.UserProfileUpdate{
			NewName: "Updated User",
			Email:   "updateduser@example.com",
		}
		mockUserRepo.EXPECT().UpdateUserProfile(1, *updateData).Return(nil, errors.New("update failed"))

		profile, err := userUsecase.UpdateMyProfile(1, updateData)
		assert.Error(t, err)
		assert.Nil(t, profile)
	})
}

// func TestUserUsecase_SetMyAvatar(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockUserRepo := mocks.NewMockUserRepo(ctrl)
// 	userUsecase := usecase.CreateUserUsecase(mockUserRepo)

// 	// Мокаем загрузку аватара
// 	t.Run("successful avatar upload", func(t *testing.T) {
// 		// Предполагаемая настройка окружения для загрузки файлов
// 		os.Setenv("USER_UPLOADS_DIR", "/tmp/uploads")

// 		file := multipart.File(nil) // Заглушка для теста
// 		fileHeader := &multipart.FileHeader{
// 			Filename: "avatar.jpg",
// 			Size:     1024,
// 		}

// 		mockUserRepo.EXPECT().SetUserAvatar(1, ".jpg", int(fileHeader.Size)).Return("user1_avatar.jpg", nil)
// 		mockUserRepo.EXPECT().GetUserProfile(1).Return(&models.UserProfile{
// 			ID:             1,
// 			Name:           "Test User",
// 			Email:          "testuser@example.com",
// 			AvatarImageURL: "/tmp/uploads/user1_avatar.jpg",
// 		}, nil)

// 		profile, err := userUsecase.SetMyAvatar(1, &file, fileHeader)
// 		assert.NoError(t, err)
// 		assert.Equal(t, "/tmp/uploads/user1_avatar.jpg", profile.AvatarImageURL)
// 	})

// 	t.Run("failed avatar upload - unable to save", func(t *testing.T) {
// 		file := multipart.File(nil) // Заглушка для теста
// 		fileHeader := &multipart.FileHeader{
// 			Filename: "avatar.jpg",
// 			Size:     1024,
// 		}

// 		mockUserRepo.EXPECT().SetUserAvatar(1, ".jpg", int(fileHeader.Size)).Return("", errors.New("upload error"))

// 		profile, err := userUsecase.SetMyAvatar(1, &file, fileHeader)
// 		assert.Error(t, err)
// 		assert.Nil(t, profile)
// 	})
// }
