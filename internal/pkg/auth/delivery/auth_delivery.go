package delivery

import (
	"RPO_back/internal/pkg/auth"
	"RPO_back/internal/pkg/auth/delivery/grpc/gen"
	"context"
)

type AuthServer struct {
	gen.AuthServer
	authUsecase auth.AuthUsecase
}

func CreateAuthServer(uc auth.AuthUsecase) *AuthServer {
	return &AuthServer{
		authUsecase: uc,
	}
}

func (d *AuthServer) CreateSession(ctx context.Context, request *gen.UserDataRequest) (*gen.Session, error) {
	sessionID, err := d.authUsecase.CreateSession(ctx, request.UserID, request.Password)
	if err != nil {
		return &gen.Session{Error: gen.Error_INVALID_CREDENTIALS}, err
	}

	return &gen.Session{SessionID: sessionID, Error: gen.Error_NONE}, nil
}

func (d *AuthServer) CheckSession(ctx context.Context, request *gen.Session) (*gen.UserDataResponse, error) {
	userID, err := d.authUsecase.CheckSession(ctx, request.SessionID)
	if err != nil {
		return &gen.UserDataResponse{Error: gen.Error_INVALID_CREDENTIALS}, err
	}

	return &gen.UserDataResponse{UserID: userID, Error: gen.Error_NONE}, nil
}

func (d *AuthServer) DeleteSession(ctx context.Context, request *gen.Session) (*gen.StatusResponse, error) {
	err := d.authUsecase.KillSession(ctx, request.SessionID)
	if err != nil {
		return &gen.StatusResponse{Error: gen.Error_INTERNAL_SERVER_ERROR}, nil
	}

	return &gen.StatusResponse{Error: gen.Error_NONE}, nil
}

func (d *AuthServer) ChangePassword(ctx context.Context, request *gen.ChangePasswordRequest) (*gen.StatusResponse, error) {
	err := d.authUsecase.ChangePassword(ctx, request.PasswordOld, request.PasswordNew)
	if err != nil {
		return &gen.StatusResponse{Error: gen.Error_INTERNAL_SERVER_ERROR}, nil
	}

	return &gen.StatusResponse{Error: gen.Error_NONE}, nil
}
