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
	if err != nil {
		return nil, fmt.Errorf("GetMyProfile: %w", err)
	}

	return
}

// UpdateMyProfile обновляет профиль пользователя и возвращает обновлённый профиль
func (uc *UserUsecase) UpdateMyProfile(ctx context.Context, userID int64, data *models.UserProfileUpdateRequest) (updatedProfile *models.UserProfile, err error) {
	updatedProfile, err = uc.userRepo.UpdateUserProfile(ctx, userID, *data)
	if err != nil {
		return nil, fmt.Errorf("GetMyProfile: %w", err)
	}

	return
}

// SetMyAvatar устанавливает пользователю аватарку
func (uc *UserUsecase) SetMyAvatar(ctx context.Context, userID int64, file *models.UploadedFile) (updated *models.UserProfile, err error) {
	funcName := "SetMyAvatar"
	fileID, err := uploads.UsecaseUploadFile(ctx, file, uc.userRepo)
	if err != nil {
		return nil, fmt.Errorf("%s (upload file): %w", funcName, err)
	}

	err = uc.userRepo.SetUserAvatar(ctx, userID, fileID)

	return uc.userRepo.GetUserProfile(ctx, userID)
}

func (uc *UserUsecase) ChangePassword(ctx context.Context, sessionID string, oldPassword string, newPassword string) error {
	responce, err := uc.authClient.ChangePassword(ctx, &authGRPC.ChangePasswordRequest{
		PasswordOld: oldPassword,
		PasswordNew: newPassword,
		SessionID:   sessionID,
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
	funcName := "RegisterUser"

	err = uc.userRepo.CheckUniqueCredentials(ctx, user.Name, user.Email)
	if err != nil {
		return "", fmt.Errorf("%s (uniq): %w", funcName, err)
	}

	newUser, err := uc.userRepo.CreateUser(ctx, user)
	if err != nil {
		return "", fmt.Errorf("%s (create): %w", funcName, err)
	}

	response, err := uc.authClient.CreateSession(ctx, &authGRPC.UserDataRequest{
		UserID:   int64(newUser.ID),
		Password: user.Password,
	})
	if err != nil {
		return "", fmt.Errorf("%s (session): %w", funcName, err)
	}
	errGRPC := response.GetError()
	if errGRPC == authGRPC.Error_INVALID_CREDENTIALS {
		return "", fmt.Errorf("%s (GRPC response 1): %w", funcName, errs.ErrWrongCredentials)
	} else if errGRPC == authGRPC.Error_INTERNAL_SERVER_ERROR {
		return "", fmt.Errorf("%s (GRPC response 1): internal error at auth service", funcName)
	}
	sessionID = response.GetSessionID()

	response2, err := uc.authClient.ChangePassword(ctx, &authGRPC.ChangePasswordRequest{
		PasswordOld: "",
		PasswordNew: user.Password,
		SessionID:   sessionID,
	})
	if err != nil {
		return "", fmt.Errorf("%s (password): %w", funcName, err)
	}
	errGRPC = response2.GetError()
	if errGRPC == authGRPC.Error_INVALID_CREDENTIALS {
		return "", fmt.Errorf("%s (GRPC response 2): %w", funcName, errs.ErrWrongCredentials)
	} else if errGRPC == authGRPC.Error_INTERNAL_SERVER_ERROR {
		return "", fmt.Errorf("%s (GRPC response 2): internal error at auth service", funcName)
	}

	return sessionID, nil
}
