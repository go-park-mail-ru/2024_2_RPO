package usecase

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	authGRPC "RPO_back/internal/pkg/auth/delivery/grpc/gen"
	"RPO_back/internal/pkg/user"
	"RPO_back/internal/pkg/utils/uploads"
	"context"
	"fmt"
)

type UserUsecase struct {
	authClient authGRPC.AuthClient
	userRepo   user.UserRepo
}

func CreateUserUsecase(userRepo user.UserRepo, authClient authGRPC.AuthClient) *UserUsecase {
	return &UserUsecase{
		authClient: authClient,
		userRepo:   userRepo,
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
	if err != nil {
		return nil, fmt.Errorf("SetMyAvatar: %w", err)
	}

	existingID, err := uploads.CompareFiles(fileNames, fileIDs, file)
	if err != nil {
		return nil, fmt.Errorf("SetMyAvatar: %w", err)
	}

	if existingID != nil {
		file.FileID = existingID
	} else {
		uc.userRepo.RegisterFile(ctx, file)
		uploads.SaveFile(file)
	}

	uploads.SaveFile(file)

	return uc.userRepo.GetUserProfile(ctx, userID)
}

func (uc *UserUsecase) ChangePassword(ctx context.Context, userID int64, oldPassword string, newPassword string) error {
	responce, err := uc.authClient.ChangePassword(ctx, &authGRPC.ChangePasswordRequest{
		PasswordOld: oldPassword,
		PasswordNew: newPassword,
	})
	if err != nil {
		return fmt.Errorf("ChangePassword: %w", err)
	}

	errGRPC := responce.GetError()
	if errGRPC == authGRPC.Error_INVALID_CREDENTIALS {
		return fmt.Errorf("ChangePassword: %w", errs.ErrWrongCredentials)
	} else if errGRPC == authGRPC.Error_INTERNAL_SERVER_ERROR {
		return fmt.Errorf("ChangePassword: internal error at auth service")
	}

	return nil
}

func (uc *UserUsecase) LoginUser(ctx context.Context, email string, password string) (sessionID string, err error) {
	user, err := uc.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("LoginUser (GetUserByEmail): %w", err)
	}

	userID := user.ID

	responce, err := uc.authClient.CreateSession(ctx, &authGRPC.UserDataRequest{
		UserID:   int64(userID),
		Password: password,
	})
	if err != nil {
		return "", fmt.Errorf("LoginUser (GRPC request): %w", err)
	}

	errGRPC := responce.GetError()
	if errGRPC == authGRPC.Error_INVALID_CREDENTIALS {
		return "", fmt.Errorf("CreateSession (GRPC response): %w", errs.ErrWrongCredentials)
	} else if errGRPC == authGRPC.Error_INTERNAL_SERVER_ERROR {
		return "", fmt.Errorf("CreateSession (GRPC response): internal error at auth service")
	}

	sessionID = responce.GetSessionID()

	return sessionID, nil
}

func (uc *UserUsecase) LogoutUser(ctx context.Context, sessionID string) error {
	responce, err := uc.authClient.DeleteSession(ctx, &authGRPC.Session{SessionID: sessionID})
	if err != nil {
		return fmt.Errorf("LogoutUser (GRPC request): %w", err)
	}

	errGRPC := responce.GetError()
	if errGRPC == authGRPC.Error_INVALID_CREDENTIALS {
		return fmt.Errorf("DeleteSession (GRPC response): %w", errs.ErrWrongCredentials)
	} else if errGRPC == authGRPC.Error_INTERNAL_SERVER_ERROR {
		return fmt.Errorf("DeleteSession (GRPC response): internal error at auth service")
	}

	return nil
}

func (uc *UserUsecase) RegisterUser(ctx context.Context, user *models.UserRegisterRequest) (sessionID string, err error) {
	newUser, err := uc.userRepo.CreateUser(ctx, user)
	if err != nil {
		return "", fmt.Errorf("RegisterUser (CreateUser): %w", err)
	}

	responce, err := uc.authClient.CreateSession(ctx, &authGRPC.UserDataRequest{
		UserID:   int64(newUser.ID),
		Password: user.Password,
	})
	if err != nil {
		return "", fmt.Errorf("RegisterUser (GRPC request): %w", err)
	}

	errGRPC := responce.GetError()
	if errGRPC == authGRPC.Error_INVALID_CREDENTIALS {
		return "", fmt.Errorf("CreateSession (GRPC response): %w", errs.ErrWrongCredentials)
	} else if errGRPC == authGRPC.Error_INTERNAL_SERVER_ERROR {
		return "", fmt.Errorf("CreateSession (GRPC response): internal error at auth service")
	}

	sessionID = responce.GetSessionID()

	return sessionID, nil
}
