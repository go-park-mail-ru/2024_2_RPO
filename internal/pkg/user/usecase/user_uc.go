package usecase

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/user"
	"RPO_back/internal/pkg/utils/uploads"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type UserUsecase struct {
	userRepo user.UserRepo
}

func CreateUserUsecase(userRepo user.UserRepo) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

// GetMyProfile возвращает пользователю его профиль
func (uc *UserUsecase) GetMyProfile(userID int) (profile *models.UserProfile, err error) {
	profile, err = uc.userRepo.GetUserProfile(userID)
	return
}

// UpdateMyProfile обновляет профиль пользователя и возвращает обновлённый профиль
func (uc *UserUsecase) UpdateMyProfile(userID int, data *models.UserProfileUpdate) (updatedProfile *models.UserProfile, err error) {
	updatedProfile, err = uc.userRepo.UpdateUserProfile(userID, *data)
	return
}

// SetMyAvatar устанавливает пользователю аватарку
func (uc *UserUsecase) SetMyAvatar(userID int, file *multipart.File, fileHeader *multipart.FileHeader) (updated *models.UserProfile, err error) {
	fileName, err := uc.userRepo.SetUserAvatar(userID, uploads.ExtractFileExtension(fileHeader.Filename), int(fileHeader.Size))
	if err != nil {
		return nil, fmt.Errorf("SetMyAvatar: %w", err)
	}
	uploadDir := os.Getenv("USER_UPLOADS_DIR")

	filePath := filepath.Join(uploadDir, fileName)
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("Cant create file on server side: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, *file); err != nil {
		return nil, fmt.Errorf("Cant copy file on server side: %w", err)
	}

	return uc.userRepo.GetUserProfile(userID)
}
