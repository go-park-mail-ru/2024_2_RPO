package usecase

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/user"
	"RPO_back/internal/pkg/utils/uploads"
	"context"
	"fmt"
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
func (uc *UserUsecase) GetMyProfile(ctx context.Context, userID int64) (profile *models.UserProfile, err error) {
	profile, err = uc.userRepo.GetUserProfile(ctx, userID)
	return
}

// UpdateMyProfile обновляет профиль пользователя и возвращает обновлённый профиль
func (uc *UserUsecase) UpdateMyProfile(ctx context.Context, userID int64, data *models.UserProfileUpdateRequest) (updatedProfile *models.UserProfile, err error) {
	updatedProfile, err = uc.userRepo.UpdateUserProfile(ctx, userID, *data)
	return
}

// SetMyAvatar устанавливает пользователю аватарку
func (uc *UserUsecase) SetMyAvatar(ctx context.Context, userID int64, file *models.UploadedFile) (updated *models.UserProfile, err error) {
	fileNames, fileIDs, err := uc.userRepo.DeduplicateFile(ctx, file)
	existingID, err := uploads.CompareFiles(fileNames, fileIDs, file)
	if existingID != nil {
		file.FileID = existingID
	} else {
		uc.userRepo.RegisterFile(ctx, file)
		uploads.SaveFile(file)
	}
	if err != nil {
		return nil, fmt.Errorf("SetMyAvatar: %w", err)
	}

	uploads.SaveFile(file)

	return uc.userRepo.GetUserProfile(ctx, userID)
}

func (uc *UserUsecase) ChangePassword(ctx context.Context, userID int64, oldPassword string, newPassword string) error {
	panic("not implemented")
}

func (uc *UserUsecase) LoginUser(ctx context.Context, email string, password string) (sessionID string, err error) {
	panic("not implemented")
}

func (uc *UserUsecase) LogoutUser(ctx context.Context, sessionID string) error {
	panic("not implemented")
}

func (uc *UserUsecase) RegisterUser(ctx context.Context, user *models.UserRegisterRequest) (sessionID string, err error) {
	panic("not implemented")
}
