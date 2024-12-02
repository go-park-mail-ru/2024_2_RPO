package delivery

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/pkg/auth"
	"RPO_back/internal/pkg/auth/delivery/grpc/gen"
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=grpc/gen/auth_grpc.pb.go -destination=grpc/mocks/auth_mock.go

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
	funcName := "CreateSession"
	sessionID, err := d.authUsecase.CreateSession(ctx, request.UserID, request.Password)
	if err != nil {
		logrus.Errorf("%s: %v", funcName, err)
		if errors.Is(err, errs.ErrWrongCredentials) {
			return &gen.Session{Error: gen.Error_INVALID_CREDENTIALS}, nil
		}
		return &gen.Session{Error: gen.Error_INTERNAL_SERVER_ERROR}, nil
	}

	return &gen.Session{SessionID: sessionID, Error: gen.Error_NONE}, nil
}

func (d *AuthDelivery) CheckSession(ctx context.Context, request *gen.CheckSessionRequest) (*gen.UserDataResponse, error) {
	funcName := "CheckSession"
	userID, err := d.authUsecase.CheckSession(ctx, request.SessionID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return &gen.UserDataResponse{Error: gen.Error_INVALID_CREDENTIALS}, nil
		}
		logrus.Errorf("%s: %v", funcName, err)
		return nil, fmt.Errorf("CheckSession: %w", err)
	}

	return &gen.UserDataResponse{UserID: int64(userID), Error: gen.Error_NONE}, nil
}

func (d *AuthDelivery) DeleteSession(ctx context.Context, request *gen.Session) (*gen.StatusResponse, error) {
	funcName := "DeleteSession"
	err := d.authUsecase.KillSession(ctx, request.SessionID)
	if err != nil {
		logrus.Errorf("%s: %v", funcName, err)
		return &gen.StatusResponse{Error: gen.Error_INTERNAL_SERVER_ERROR}, nil
	}

	return &gen.StatusResponse{Error: gen.Error_NONE}, nil
}

func (d *AuthDelivery) ChangePassword(ctx context.Context, request *gen.ChangePasswordRequest) (*gen.StatusResponse, error) {
	funcName := "ChangePassword"
	err := d.authUsecase.ChangePassword(ctx, request.PasswordOld, request.PasswordNew, request.SessionID)
	if err != nil {
		logrus.Errorf("%s: %v", funcName, err)
		return &gen.StatusResponse{Error: gen.Error_INTERNAL_SERVER_ERROR}, nil
	}

	return &gen.StatusResponse{Error: gen.Error_NONE}, nil
}
