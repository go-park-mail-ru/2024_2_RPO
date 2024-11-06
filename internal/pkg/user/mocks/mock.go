// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package mock_user is a generated GoMock package.
package mock_user

import (
	models "RPO_back/internal/models"
	context "context"
	multipart "mime/multipart"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserUsecase is a mock of UserUsecase interface.
type MockUserUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockUserUsecaseMockRecorder
}

// MockUserUsecaseMockRecorder is the mock recorder for MockUserUsecase.
type MockUserUsecaseMockRecorder struct {
	mock *MockUserUsecase
}

// NewMockUserUsecase creates a new mock instance.
func NewMockUserUsecase(ctrl *gomock.Controller) *MockUserUsecase {
	mock := &MockUserUsecase{ctrl: ctrl}
	mock.recorder = &MockUserUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserUsecase) EXPECT() *MockUserUsecaseMockRecorder {
	return m.recorder
}

// GetMyProfile mocks base method.
func (m *MockUserUsecase) GetMyProfile(ctx context.Context, userID int) (*models.UserProfile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMyProfile", ctx, userID)
	ret0, _ := ret[0].(*models.UserProfile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMyProfile indicates an expected call of GetMyProfile.
func (mr *MockUserUsecaseMockRecorder) GetMyProfile(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMyProfile", reflect.TypeOf((*MockUserUsecase)(nil).GetMyProfile), ctx, userID)
}

// SetMyAvatar mocks base method.
func (m *MockUserUsecase) SetMyAvatar(ctx context.Context, userID int, file *multipart.File, fileHeader *multipart.FileHeader) (*models.UserProfile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetMyAvatar", ctx, userID, file, fileHeader)
	ret0, _ := ret[0].(*models.UserProfile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetMyAvatar indicates an expected call of SetMyAvatar.
func (mr *MockUserUsecaseMockRecorder) SetMyAvatar(ctx, userID, file, fileHeader interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMyAvatar", reflect.TypeOf((*MockUserUsecase)(nil).SetMyAvatar), ctx, userID, file, fileHeader)
}

// UpdateMyProfile mocks base method.
func (m *MockUserUsecase) UpdateMyProfile(ctx context.Context, userID int, data *models.UserProfileUpdate) (*models.UserProfile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMyProfile", ctx, userID, data)
	ret0, _ := ret[0].(*models.UserProfile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateMyProfile indicates an expected call of UpdateMyProfile.
func (mr *MockUserUsecaseMockRecorder) UpdateMyProfile(ctx, userID, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMyProfile", reflect.TypeOf((*MockUserUsecase)(nil).UpdateMyProfile), ctx, userID, data)
}

// MockUserRepo is a mock of UserRepo interface.
type MockUserRepo struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepoMockRecorder
}

// MockUserRepoMockRecorder is the mock recorder for MockUserRepo.
type MockUserRepoMockRecorder struct {
	mock *MockUserRepo
}

// NewMockUserRepo creates a new mock instance.
func NewMockUserRepo(ctrl *gomock.Controller) *MockUserRepo {
	mock := &MockUserRepo{ctrl: ctrl}
	mock.recorder = &MockUserRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepo) EXPECT() *MockUserRepoMockRecorder {
	return m.recorder
}

// GetUserProfile mocks base method.
func (m *MockUserRepo) GetUserProfile(ctx context.Context, userID int) (*models.UserProfile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserProfile", ctx, userID)
	ret0, _ := ret[0].(*models.UserProfile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserProfile indicates an expected call of GetUserProfile.
func (mr *MockUserRepoMockRecorder) GetUserProfile(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserProfile", reflect.TypeOf((*MockUserRepo)(nil).GetUserProfile), ctx, userID)
}

// SetUserAvatar mocks base method.
func (m *MockUserRepo) SetUserAvatar(ctx context.Context, userID int, fileExtension string, fileSize int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetUserAvatar", ctx, userID, fileExtension, fileSize)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetUserAvatar indicates an expected call of SetUserAvatar.
func (mr *MockUserRepoMockRecorder) SetUserAvatar(ctx, userID, fileExtension, fileSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetUserAvatar", reflect.TypeOf((*MockUserRepo)(nil).SetUserAvatar), ctx, userID, fileExtension, fileSize)
}

// UpdateUserProfile mocks base method.
func (m *MockUserRepo) UpdateUserProfile(ctx context.Context, userID int, data models.UserProfileUpdate) (*models.UserProfile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserProfile", ctx, userID, data)
	ret0, _ := ret[0].(*models.UserProfile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUserProfile indicates an expected call of UpdateUserProfile.
func (mr *MockUserRepoMockRecorder) UpdateUserProfile(ctx, userID, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserProfile", reflect.TypeOf((*MockUserRepo)(nil).UpdateUserProfile), ctx, userID, data)
}
