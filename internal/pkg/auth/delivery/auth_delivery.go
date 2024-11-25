package delivery

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/pkg/auth"
	"RPO_back/internal/pkg/auth/delivery/grpc/gen"
	"context"
	"errors"
	"fmt"
)

type AuthDelivery struct {
	gen.AuthServer
	authUsecase auth.AuthUsecase
}

func CreateAuthServer(uc auth.AuthUsecase) *AuthDelivery {
	return &AuthDelivery{
		authUsecase: uc,
	}
}

func (d *AuthDelivery) CreateSession(ctx context.Context, request *gen.UserDataRequest) (*gen.Session, error) {
	sessionID, err := d.authUsecase.CreateSession(ctx, request.UserID, request.Password)
	if err != nil {
		return &gen.Session{Error: gen.Error_INVALID_CREDENTIALS}, err
	}

	return &gen.Session{SessionID: sessionID, Error: gen.Error_NONE}, nil
}

func (d *AuthDelivery) CheckSession(ctx context.Context, request *gen.CheckSessionRequest) (*gen.UserDataResponse, error) {
	userID, err := d.authUsecase.CheckSession(ctx, request.SessionID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return &gen.UserDataResponse{Error: gen.Error_INVALID_CREDENTIALS}, nil
		}
		return nil, fmt.Errorf("CheckSession: %w", err)
	}

	return &gen.UserDataResponse{UserID: int64(userID), Error: gen.Error_NONE}, nil
}

func (d *AuthDelivery) DeleteSession(ctx context.Context, request *gen.Session) (*gen.StatusResponse, error) {
	err := d.authUsecase.KillSession(ctx, request.SessionID)
	if err != nil {
		return &gen.StatusResponse{Error: gen.Error_INTERNAL_SERVER_ERROR}, nil
	}

	return &gen.StatusResponse{Error: gen.Error_NONE}, nil
}

func (d *AuthDelivery) ChangePassword(ctx context.Context, request *gen.ChangePasswordRequest) (*gen.StatusResponse, error) {
	err := d.authUsecase.ChangePassword(ctx, request.PasswordOld, request.PasswordNew, request.SessionID)
	if err != nil {
		return &gen.StatusResponse{Error: gen.Error_INTERNAL_SERVER_ERROR}, nil
	}

	return &gen.StatusResponse{Error: gen.Error_NONE}, nil
}
