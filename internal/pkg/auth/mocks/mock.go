// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package mock_auth is a generated GoMock package.
package mock_auth

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAuthUsecase is a mock of AuthUsecase interface.
type MockAuthUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockAuthUsecaseMockRecorder
}

// MockAuthUsecaseMockRecorder is the mock recorder for MockAuthUsecase.
type MockAuthUsecaseMockRecorder struct {
	mock *MockAuthUsecase
}

// NewMockAuthUsecase creates a new mock instance.
func NewMockAuthUsecase(ctrl *gomock.Controller) *MockAuthUsecase {
	mock := &MockAuthUsecase{ctrl: ctrl}
	mock.recorder = &MockAuthUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthUsecase) EXPECT() *MockAuthUsecaseMockRecorder {
	return m.recorder
}

// ChangePassword mocks base method.
func (m *MockAuthUsecase) ChangePassword(ctx context.Context, oldPassword, newPassword, sessionID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangePassword", ctx, oldPassword, newPassword, sessionID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangePassword indicates an expected call of ChangePassword.
func (mr *MockAuthUsecaseMockRecorder) ChangePassword(ctx, oldPassword, newPassword, sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangePassword", reflect.TypeOf((*MockAuthUsecase)(nil).ChangePassword), ctx, oldPassword, newPassword, sessionID)
}

// CheckSession mocks base method.
func (m *MockAuthUsecase) CheckSession(ctx context.Context, sessionID string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckSession", ctx, sessionID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckSession indicates an expected call of CheckSession.
func (mr *MockAuthUsecaseMockRecorder) CheckSession(ctx, sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckSession", reflect.TypeOf((*MockAuthUsecase)(nil).CheckSession), ctx, sessionID)
}

// CreateSession mocks base method.
func (m *MockAuthUsecase) CreateSession(ctx context.Context, userID int64, password string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", ctx, userID, password)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockAuthUsecaseMockRecorder) CreateSession(ctx, userID, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockAuthUsecase)(nil).CreateSession), ctx, userID, password)
}

// KillSession mocks base method.
func (m *MockAuthUsecase) KillSession(ctx context.Context, sessionID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KillSession", ctx, sessionID)
	ret0, _ := ret[0].(error)
	return ret0
}

// KillSession indicates an expected call of KillSession.
func (mr *MockAuthUsecaseMockRecorder) KillSession(ctx, sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KillSession", reflect.TypeOf((*MockAuthUsecase)(nil).KillSession), ctx, sessionID)
}

// MockAuthRepo is a mock of AuthRepo interface.
type MockAuthRepo struct {
	ctrl     *gomock.Controller
	recorder *MockAuthRepoMockRecorder
}

// MockAuthRepoMockRecorder is the mock recorder for MockAuthRepo.
type MockAuthRepoMockRecorder struct {
	mock *MockAuthRepo
}

// NewMockAuthRepo creates a new mock instance.
func NewMockAuthRepo(ctrl *gomock.Controller) *MockAuthRepo {
	mock := &MockAuthRepo{ctrl: ctrl}
	mock.recorder = &MockAuthRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthRepo) EXPECT() *MockAuthRepoMockRecorder {
	return m.recorder
}

// CheckSession mocks base method.
func (m *MockAuthRepo) CheckSession(ctx context.Context, sessionID string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckSession", ctx, sessionID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckSession indicates an expected call of CheckSession.
func (mr *MockAuthRepoMockRecorder) CheckSession(ctx, sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckSession", reflect.TypeOf((*MockAuthRepo)(nil).CheckSession), ctx, sessionID)
}

// DisplaceUserSessions mocks base method.
func (m *MockAuthRepo) DisplaceUserSessions(ctx context.Context, sessionID string, userID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DisplaceUserSessions", ctx, sessionID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DisplaceUserSessions indicates an expected call of DisplaceUserSessions.
func (mr *MockAuthRepoMockRecorder) DisplaceUserSessions(ctx, sessionID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisplaceUserSessions", reflect.TypeOf((*MockAuthRepo)(nil).DisplaceUserSessions), ctx, sessionID, userID)
}

// GetUserPasswordHash mocks base method.
func (m *MockAuthRepo) GetUserPasswordHash(ctx context.Context, userID int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserPasswordHash", ctx, userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserPasswordHash indicates an expected call of GetUserPasswordHash.
func (mr *MockAuthRepoMockRecorder) GetUserPasswordHash(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserPasswordHash", reflect.TypeOf((*MockAuthRepo)(nil).GetUserPasswordHash), ctx, userID)
}

// KillSessionRedis mocks base method.
func (m *MockAuthRepo) KillSessionRedis(ctx context.Context, sessionID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KillSessionRedis", ctx, sessionID)
	ret0, _ := ret[0].(error)
	return ret0
}

// KillSessionRedis indicates an expected call of KillSessionRedis.
func (mr *MockAuthRepoMockRecorder) KillSessionRedis(ctx, sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KillSessionRedis", reflect.TypeOf((*MockAuthRepo)(nil).KillSessionRedis), ctx, sessionID)
}

// RegisterSessionRedis mocks base method.
func (m *MockAuthRepo) RegisterSessionRedis(ctx context.Context, cookie string, userID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterSessionRedis", ctx, cookie, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterSessionRedis indicates an expected call of RegisterSessionRedis.
func (mr *MockAuthRepoMockRecorder) RegisterSessionRedis(ctx, cookie, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterSessionRedis", reflect.TypeOf((*MockAuthRepo)(nil).RegisterSessionRedis), ctx, cookie, userID)
}

// SetNewPasswordHash mocks base method.
func (m *MockAuthRepo) SetNewPasswordHash(ctx context.Context, userID int, newPasswordHash string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetNewPasswordHash", ctx, userID, newPasswordHash)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetNewPasswordHash indicates an expected call of SetNewPasswordHash.
func (mr *MockAuthRepoMockRecorder) SetNewPasswordHash(ctx, userID, newPasswordHash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetNewPasswordHash", reflect.TypeOf((*MockAuthRepo)(nil).SetNewPasswordHash), ctx, userID, newPasswordHash)
}
