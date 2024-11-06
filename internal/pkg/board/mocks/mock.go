// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package mock_board is a generated GoMock package.
package mock_board

import (
	models "RPO_back/internal/models"
	context "context"
	multipart "mime/multipart"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockBoardUsecase is a mock of BoardUsecase interface.
type MockBoardUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockBoardUsecaseMockRecorder
}

// MockBoardUsecaseMockRecorder is the mock recorder for MockBoardUsecase.
type MockBoardUsecaseMockRecorder struct {
	mock *MockBoardUsecase
}

// NewMockBoardUsecase creates a new mock instance.
func NewMockBoardUsecase(ctrl *gomock.Controller) *MockBoardUsecase {
	mock := &MockBoardUsecase{ctrl: ctrl}
	mock.recorder = &MockBoardUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBoardUsecase) EXPECT() *MockBoardUsecaseMockRecorder {
	return m.recorder
}

// AddMember mocks base method.
func (m *MockBoardUsecase) AddMember(ctx context.Context, userID, boardID int, addRequest *models.AddMemberRequest) (*models.MemberWithPermissions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMember", ctx, userID, boardID, addRequest)
	ret0, _ := ret[0].(*models.MemberWithPermissions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddMember indicates an expected call of AddMember.
func (mr *MockBoardUsecaseMockRecorder) AddMember(ctx, userID, boardID, addRequest interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMember", reflect.TypeOf((*MockBoardUsecase)(nil).AddMember), ctx, userID, boardID, addRequest)
}

// CreateColumn mocks base method.
func (m *MockBoardUsecase) CreateColumn(ctx context.Context, userID, boardID int, data *models.ColumnRequest) (*models.Column, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateColumn", ctx, userID, boardID, data)
	ret0, _ := ret[0].(*models.Column)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateColumn indicates an expected call of CreateColumn.
func (mr *MockBoardUsecaseMockRecorder) CreateColumn(ctx, userID, boardID, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateColumn", reflect.TypeOf((*MockBoardUsecase)(nil).CreateColumn), ctx, userID, boardID, data)
}

// CreateNewBoard mocks base method.
func (m *MockBoardUsecase) CreateNewBoard(ctx context.Context, userID int, data models.CreateBoardRequest) (*models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNewBoard", ctx, userID, data)
	ret0, _ := ret[0].(*models.Board)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNewBoard indicates an expected call of CreateNewBoard.
func (mr *MockBoardUsecaseMockRecorder) CreateNewBoard(ctx, userID, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNewBoard", reflect.TypeOf((*MockBoardUsecase)(nil).CreateNewBoard), ctx, userID, data)
}

// CreateNewCard mocks base method.
func (m *MockBoardUsecase) CreateNewCard(ctx context.Context, userID, boardID int, data *models.CardPutRequest) (*models.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNewCard", ctx, userID, boardID, data)
	ret0, _ := ret[0].(*models.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNewCard indicates an expected call of CreateNewCard.
func (mr *MockBoardUsecaseMockRecorder) CreateNewCard(ctx, userID, boardID, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNewCard", reflect.TypeOf((*MockBoardUsecase)(nil).CreateNewCard), ctx, userID, boardID, data)
}

// DeleteBoard mocks base method.
func (m *MockBoardUsecase) DeleteBoard(ctx context.Context, userID, boardID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBoard", ctx, userID, boardID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBoard indicates an expected call of DeleteBoard.
func (mr *MockBoardUsecaseMockRecorder) DeleteBoard(ctx, userID, boardID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBoard", reflect.TypeOf((*MockBoardUsecase)(nil).DeleteBoard), ctx, userID, boardID)
}

// DeleteCard mocks base method.
func (m *MockBoardUsecase) DeleteCard(ctx context.Context, userID, boardID, cardID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCard", ctx, userID, boardID, cardID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCard indicates an expected call of DeleteCard.
func (mr *MockBoardUsecaseMockRecorder) DeleteCard(ctx, userID, boardID, cardID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCard", reflect.TypeOf((*MockBoardUsecase)(nil).DeleteCard), ctx, userID, boardID, cardID)
}

// DeleteColumn mocks base method.
func (m *MockBoardUsecase) DeleteColumn(ctx context.Context, userID, boardID, columnID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteColumn", ctx, userID, boardID, columnID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteColumn indicates an expected call of DeleteColumn.
func (mr *MockBoardUsecaseMockRecorder) DeleteColumn(ctx, userID, boardID, columnID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteColumn", reflect.TypeOf((*MockBoardUsecase)(nil).DeleteColumn), ctx, userID, boardID, columnID)
}

// GetBoardContent mocks base method.
func (m *MockBoardUsecase) GetBoardContent(ctx context.Context, userID, boardID int) (*models.BoardContent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBoardContent", ctx, userID, boardID)
	ret0, _ := ret[0].(*models.BoardContent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBoardContent indicates an expected call of GetBoardContent.
func (mr *MockBoardUsecaseMockRecorder) GetBoardContent(ctx, userID, boardID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBoardContent", reflect.TypeOf((*MockBoardUsecase)(nil).GetBoardContent), ctx, userID, boardID)
}

// GetMembersPermissions mocks base method.
func (m *MockBoardUsecase) GetMembersPermissions(ctx context.Context, userID, boardID int) ([]models.MemberWithPermissions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMembersPermissions", ctx, userID, boardID)
	ret0, _ := ret[0].([]models.MemberWithPermissions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMembersPermissions indicates an expected call of GetMembersPermissions.
func (mr *MockBoardUsecaseMockRecorder) GetMembersPermissions(ctx, userID, boardID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMembersPermissions", reflect.TypeOf((*MockBoardUsecase)(nil).GetMembersPermissions), ctx, userID, boardID)
}

// GetMyBoards mocks base method.
func (m *MockBoardUsecase) GetMyBoards(ctx context.Context, userID int) ([]models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMyBoards", ctx, userID)
	ret0, _ := ret[0].([]models.Board)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMyBoards indicates an expected call of GetMyBoards.
func (mr *MockBoardUsecaseMockRecorder) GetMyBoards(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMyBoards", reflect.TypeOf((*MockBoardUsecase)(nil).GetMyBoards), ctx, userID)
}

// RemoveMember mocks base method.
func (m *MockBoardUsecase) RemoveMember(ctx context.Context, userID, boardID, memberID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveMember", ctx, userID, boardID, memberID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveMember indicates an expected call of RemoveMember.
func (mr *MockBoardUsecaseMockRecorder) RemoveMember(ctx, userID, boardID, memberID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveMember", reflect.TypeOf((*MockBoardUsecase)(nil).RemoveMember), ctx, userID, boardID, memberID)
}

// SetBoardBackground mocks base method.
func (m *MockBoardUsecase) SetBoardBackground(ctx context.Context, userID, boardID int, file *multipart.File, fileHeader *multipart.FileHeader) (*models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetBoardBackground", ctx, userID, boardID, file, fileHeader)
	ret0, _ := ret[0].(*models.Board)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetBoardBackground indicates an expected call of SetBoardBackground.
func (mr *MockBoardUsecaseMockRecorder) SetBoardBackground(ctx, userID, boardID, file, fileHeader interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBoardBackground", reflect.TypeOf((*MockBoardUsecase)(nil).SetBoardBackground), ctx, userID, boardID, file, fileHeader)
}

// UpdateBoard mocks base method.
func (m *MockBoardUsecase) UpdateBoard(ctx context.Context, userID, boardID int, data models.BoardPutRequest) (*models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBoard", ctx, userID, boardID, data)
	ret0, _ := ret[0].(*models.Board)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateBoard indicates an expected call of UpdateBoard.
func (mr *MockBoardUsecaseMockRecorder) UpdateBoard(ctx, userID, boardID, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBoard", reflect.TypeOf((*MockBoardUsecase)(nil).UpdateBoard), ctx, userID, boardID, data)
}

// UpdateCard mocks base method.
func (m *MockBoardUsecase) UpdateCard(ctx context.Context, userID, boardID, cardID int, data *models.CardPutRequest) (*models.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCard", ctx, userID, boardID, cardID, data)
	ret0, _ := ret[0].(*models.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateCard indicates an expected call of UpdateCard.
func (mr *MockBoardUsecaseMockRecorder) UpdateCard(ctx, userID, boardID, cardID, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCard", reflect.TypeOf((*MockBoardUsecase)(nil).UpdateCard), ctx, userID, boardID, cardID, data)
}

// UpdateColumn mocks base method.
func (m *MockBoardUsecase) UpdateColumn(ctx context.Context, userID, boardID, columnID int, data *models.ColumnRequest) (*models.Column, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateColumn", ctx, userID, boardID, columnID, data)
	ret0, _ := ret[0].(*models.Column)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateColumn indicates an expected call of UpdateColumn.
func (mr *MockBoardUsecaseMockRecorder) UpdateColumn(ctx, userID, boardID, columnID, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateColumn", reflect.TypeOf((*MockBoardUsecase)(nil).UpdateColumn), ctx, userID, boardID, columnID, data)
}

// UpdateMemberRole mocks base method.
func (m *MockBoardUsecase) UpdateMemberRole(ctx context.Context, userID, boardID, memberID int, newRole string) (*models.MemberWithPermissions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMemberRole", ctx, userID, boardID, memberID, newRole)
	ret0, _ := ret[0].(*models.MemberWithPermissions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateMemberRole indicates an expected call of UpdateMemberRole.
func (mr *MockBoardUsecaseMockRecorder) UpdateMemberRole(ctx, userID, boardID, memberID, newRole interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMemberRole", reflect.TypeOf((*MockBoardUsecase)(nil).UpdateMemberRole), ctx, userID, boardID, memberID, newRole)
}

// MockBoardRepo is a mock of BoardRepo interface.
type MockBoardRepo struct {
	ctrl     *gomock.Controller
	recorder *MockBoardRepoMockRecorder
}

// MockBoardRepoMockRecorder is the mock recorder for MockBoardRepo.
type MockBoardRepoMockRecorder struct {
	mock *MockBoardRepo
}

// NewMockBoardRepo creates a new mock instance.
func NewMockBoardRepo(ctrl *gomock.Controller) *MockBoardRepo {
	mock := &MockBoardRepo{ctrl: ctrl}
	mock.recorder = &MockBoardRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBoardRepo) EXPECT() *MockBoardRepoMockRecorder {
	return m.recorder
}

// AddMember mocks base method.
func (m *MockBoardRepo) AddMember(ctx context.Context, boardID, adderID, memberUserID int) (*models.MemberWithPermissions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMember", ctx, boardID, adderID, memberUserID)
	ret0, _ := ret[0].(*models.MemberWithPermissions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddMember indicates an expected call of AddMember.
func (mr *MockBoardRepoMockRecorder) AddMember(ctx, boardID, adderID, memberUserID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMember", reflect.TypeOf((*MockBoardRepo)(nil).AddMember), ctx, boardID, adderID, memberUserID)
}

// CreateBoard mocks base method.
func (m *MockBoardRepo) CreateBoard(ctx context.Context, name string, userID int) (*models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBoard", ctx, name, userID)
	ret0, _ := ret[0].(*models.Board)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBoard indicates an expected call of CreateBoard.
func (mr *MockBoardRepoMockRecorder) CreateBoard(ctx, name, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBoard", reflect.TypeOf((*MockBoardRepo)(nil).CreateBoard), ctx, name, userID)
}

// CreateColumn mocks base method.
func (m *MockBoardRepo) CreateColumn(ctx context.Context, boardId int, title string) (*models.Column, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateColumn", ctx, boardId, title)
	ret0, _ := ret[0].(*models.Column)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateColumn indicates an expected call of CreateColumn.
func (mr *MockBoardRepoMockRecorder) CreateColumn(ctx, boardId, title interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateColumn", reflect.TypeOf((*MockBoardRepo)(nil).CreateColumn), ctx, boardId, title)
}

// CreateNewCard mocks base method.
func (m *MockBoardRepo) CreateNewCard(ctx context.Context, boardID, columnID int, title string) (*models.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNewCard", ctx, boardID, columnID, title)
	ret0, _ := ret[0].(*models.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNewCard indicates an expected call of CreateNewCard.
func (mr *MockBoardRepoMockRecorder) CreateNewCard(ctx, boardID, columnID, title interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNewCard", reflect.TypeOf((*MockBoardRepo)(nil).CreateNewCard), ctx, boardID, columnID, title)
}

// DeleteBoard mocks base method.
func (m *MockBoardRepo) DeleteBoard(ctx context.Context, boardId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBoard", ctx, boardId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBoard indicates an expected call of DeleteBoard.
func (mr *MockBoardRepoMockRecorder) DeleteBoard(ctx, boardId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBoard", reflect.TypeOf((*MockBoardRepo)(nil).DeleteBoard), ctx, boardId)
}

// DeleteCard mocks base method.
func (m *MockBoardRepo) DeleteCard(ctx context.Context, boardID, cardID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCard", ctx, boardID, cardID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCard indicates an expected call of DeleteCard.
func (mr *MockBoardRepoMockRecorder) DeleteCard(ctx, boardID, cardID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCard", reflect.TypeOf((*MockBoardRepo)(nil).DeleteCard), ctx, boardID, cardID)
}

// DeleteColumn mocks base method.
func (m *MockBoardRepo) DeleteColumn(ctx context.Context, boardID, columnID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteColumn", ctx, boardID, columnID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteColumn indicates an expected call of DeleteColumn.
func (mr *MockBoardRepoMockRecorder) DeleteColumn(ctx, boardID, columnID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteColumn", reflect.TypeOf((*MockBoardRepo)(nil).DeleteColumn), ctx, boardID, columnID)
}

// GetBoard mocks base method.
func (m *MockBoardRepo) GetBoard(ctx context.Context, boardID int) (*models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBoard", ctx, boardID)
	ret0, _ := ret[0].(*models.Board)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBoard indicates an expected call of GetBoard.
func (mr *MockBoardRepoMockRecorder) GetBoard(ctx, boardID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBoard", reflect.TypeOf((*MockBoardRepo)(nil).GetBoard), ctx, boardID)
}

// GetBoardsForUser mocks base method.
func (m *MockBoardRepo) GetBoardsForUser(ctx context.Context, userID int) ([]models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBoardsForUser", ctx, userID)
	ret0, _ := ret[0].([]models.Board)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBoardsForUser indicates an expected call of GetBoardsForUser.
func (mr *MockBoardRepoMockRecorder) GetBoardsForUser(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBoardsForUser", reflect.TypeOf((*MockBoardRepo)(nil).GetBoardsForUser), ctx, userID)
}

// GetCardsForBoard mocks base method.
func (m *MockBoardRepo) GetCardsForBoard(ctx context.Context, boardID int) ([]models.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCardsForBoard", ctx, boardID)
	ret0, _ := ret[0].([]models.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCardsForBoard indicates an expected call of GetCardsForBoard.
func (mr *MockBoardRepoMockRecorder) GetCardsForBoard(ctx, boardID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCardsForBoard", reflect.TypeOf((*MockBoardRepo)(nil).GetCardsForBoard), ctx, boardID)
}

// GetColumnsForBoard mocks base method.
func (m *MockBoardRepo) GetColumnsForBoard(ctx context.Context, boardID int) ([]models.Column, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetColumnsForBoard", ctx, boardID)
	ret0, _ := ret[0].([]models.Column)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetColumnsForBoard indicates an expected call of GetColumnsForBoard.
func (mr *MockBoardRepoMockRecorder) GetColumnsForBoard(ctx, boardID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetColumnsForBoard", reflect.TypeOf((*MockBoardRepo)(nil).GetColumnsForBoard), ctx, boardID)
}

// GetMemberPermissions mocks base method.
func (m *MockBoardRepo) GetMemberPermissions(ctx context.Context, boardID, memberUserID int, getAdderInfo bool) (*models.MemberWithPermissions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMemberPermissions", ctx, boardID, memberUserID, getAdderInfo)
	ret0, _ := ret[0].(*models.MemberWithPermissions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMemberPermissions indicates an expected call of GetMemberPermissions.
func (mr *MockBoardRepoMockRecorder) GetMemberPermissions(ctx, boardID, memberUserID, getAdderInfo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMemberPermissions", reflect.TypeOf((*MockBoardRepo)(nil).GetMemberPermissions), ctx, boardID, memberUserID, getAdderInfo)
}

// GetMembersWithPermissions mocks base method.
func (m *MockBoardRepo) GetMembersWithPermissions(ctx context.Context, boardID int) ([]models.MemberWithPermissions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMembersWithPermissions", ctx, boardID)
	ret0, _ := ret[0].([]models.MemberWithPermissions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMembersWithPermissions indicates an expected call of GetMembersWithPermissions.
func (mr *MockBoardRepoMockRecorder) GetMembersWithPermissions(ctx, boardID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMembersWithPermissions", reflect.TypeOf((*MockBoardRepo)(nil).GetMembersWithPermissions), ctx, boardID)
}

// GetUserByNickname mocks base method.
func (m *MockBoardRepo) GetUserByNickname(ctx context.Context, nickname string) (*models.UserProfile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByNickname", ctx, nickname)
	ret0, _ := ret[0].(*models.UserProfile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByNickname indicates an expected call of GetUserByNickname.
func (mr *MockBoardRepoMockRecorder) GetUserByNickname(ctx, nickname interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByNickname", reflect.TypeOf((*MockBoardRepo)(nil).GetUserByNickname), ctx, nickname)
}

// GetUserProfile mocks base method.
func (m *MockBoardRepo) GetUserProfile(ctx context.Context, userID int) (*models.UserProfile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserProfile", ctx, userID)
	ret0, _ := ret[0].(*models.UserProfile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserProfile indicates an expected call of GetUserProfile.
func (mr *MockBoardRepoMockRecorder) GetUserProfile(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserProfile", reflect.TypeOf((*MockBoardRepo)(nil).GetUserProfile), ctx, userID)
}

// RemoveMember mocks base method.
func (m *MockBoardRepo) RemoveMember(ctx context.Context, boardID, memberUserID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveMember", ctx, boardID, memberUserID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveMember indicates an expected call of RemoveMember.
func (mr *MockBoardRepoMockRecorder) RemoveMember(ctx, boardID, memberUserID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveMember", reflect.TypeOf((*MockBoardRepo)(nil).RemoveMember), ctx, boardID, memberUserID)
}

// SetBoardBackground mocks base method.
func (m *MockBoardRepo) SetBoardBackground(ctx context.Context, userID, boardID int, fileExtension string, fileSize int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetBoardBackground", ctx, userID, boardID, fileExtension, fileSize)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetBoardBackground indicates an expected call of SetBoardBackground.
func (mr *MockBoardRepoMockRecorder) SetBoardBackground(ctx, userID, boardID, fileExtension, fileSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBoardBackground", reflect.TypeOf((*MockBoardRepo)(nil).SetBoardBackground), ctx, userID, boardID, fileExtension, fileSize)
}

// SetMemberRole mocks base method.
func (m *MockBoardRepo) SetMemberRole(ctx context.Context, boardID, memberUserID int, newRole string) (*models.MemberWithPermissions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetMemberRole", ctx, boardID, memberUserID, newRole)
	ret0, _ := ret[0].(*models.MemberWithPermissions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetMemberRole indicates an expected call of SetMemberRole.
func (mr *MockBoardRepoMockRecorder) SetMemberRole(ctx, boardID, memberUserID, newRole interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMemberRole", reflect.TypeOf((*MockBoardRepo)(nil).SetMemberRole), ctx, boardID, memberUserID, newRole)
}

// UpdateBoard mocks base method.
func (m *MockBoardRepo) UpdateBoard(ctx context.Context, boardID int, data *models.BoardPutRequest) (*models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBoard", ctx, boardID, data)
	ret0, _ := ret[0].(*models.Board)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateBoard indicates an expected call of UpdateBoard.
func (mr *MockBoardRepoMockRecorder) UpdateBoard(ctx, boardID, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBoard", reflect.TypeOf((*MockBoardRepo)(nil).UpdateBoard), ctx, boardID, data)
}

// UpdateCard mocks base method.
func (m *MockBoardRepo) UpdateCard(ctx context.Context, boardID, cardID int, data models.CardPutRequest) (*models.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCard", ctx, boardID, cardID, data)
	ret0, _ := ret[0].(*models.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateCard indicates an expected call of UpdateCard.
func (mr *MockBoardRepoMockRecorder) UpdateCard(ctx, boardID, cardID, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCard", reflect.TypeOf((*MockBoardRepo)(nil).UpdateCard), ctx, boardID, cardID, data)
}

// UpdateColumn mocks base method.
func (m *MockBoardRepo) UpdateColumn(ctx context.Context, boardID, columnID int, data models.ColumnRequest) (*models.Column, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateColumn", ctx, boardID, columnID, data)
	ret0, _ := ret[0].(*models.Column)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateColumn indicates an expected call of UpdateColumn.
func (mr *MockBoardRepoMockRecorder) UpdateColumn(ctx, boardID, columnID, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateColumn", reflect.TypeOf((*MockBoardRepo)(nil).UpdateColumn), ctx, boardID, columnID, data)
}
