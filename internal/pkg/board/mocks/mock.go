// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go
//
// Generated by this command:
//
//	mockgen -source=interfaces.go -destination=mocks/mock.go
//

// Package mock_board is a generated GoMock package.
package mock_board

import (
	models "RPO_back/internal/models"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockBoardUsecase is a mock of BoardUsecase interface.
type MockBoardUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockBoardUsecaseMockRecorder
	isgomock struct{}
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
func (m *MockBoardUsecase) AddMember(userID, boardID int, addRequest *models.AddMemberRequest) (*models.MemberWithPermissions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMember", userID, boardID, addRequest)
	ret0, _ := ret[0].(*models.MemberWithPermissions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddMember indicates an expected call of AddMember.
func (mr *MockBoardUsecaseMockRecorder) AddMember(userID, boardID, addRequest any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMember", reflect.TypeOf((*MockBoardUsecase)(nil).AddMember), userID, boardID, addRequest)
}

// CreateColumn mocks base method.
func (m *MockBoardUsecase) CreateColumn(userID, boardID int, data *models.ColumnRequest) (*models.Column, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateColumn", userID, boardID, data)
	ret0, _ := ret[0].(*models.Column)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateColumn indicates an expected call of CreateColumn.
func (mr *MockBoardUsecaseMockRecorder) CreateColumn(userID, boardID, data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateColumn", reflect.TypeOf((*MockBoardUsecase)(nil).CreateColumn), userID, boardID, data)
}

// CreateNewBoard mocks base method.
func (m *MockBoardUsecase) CreateNewBoard(userID int, data models.CreateBoardRequest) (*models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNewBoard", userID, data)
	ret0, _ := ret[0].(*models.Board)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNewBoard indicates an expected call of CreateNewBoard.
func (mr *MockBoardUsecaseMockRecorder) CreateNewBoard(userID, data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNewBoard", reflect.TypeOf((*MockBoardUsecase)(nil).CreateNewBoard), userID, data)
}

// CreateNewCard mocks base method.
func (m *MockBoardUsecase) CreateNewCard(userID, boardID int, data *models.CardPutRequest) (*models.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNewCard", userID, boardID, data)
	ret0, _ := ret[0].(*models.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNewCard indicates an expected call of CreateNewCard.
func (mr *MockBoardUsecaseMockRecorder) CreateNewCard(userID, boardID, data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNewCard", reflect.TypeOf((*MockBoardUsecase)(nil).CreateNewCard), userID, boardID, data)
}

// DeleteBoard mocks base method.
func (m *MockBoardUsecase) DeleteBoard(userID, boardID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBoard", userID, boardID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBoard indicates an expected call of DeleteBoard.
func (mr *MockBoardUsecaseMockRecorder) DeleteBoard(userID, boardID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBoard", reflect.TypeOf((*MockBoardUsecase)(nil).DeleteBoard), userID, boardID)
}

// DeleteCard mocks base method.
func (m *MockBoardUsecase) DeleteCard(userID, boardID, cardID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCard", userID, boardID, cardID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCard indicates an expected call of DeleteCard.
func (mr *MockBoardUsecaseMockRecorder) DeleteCard(userID, boardID, cardID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCard", reflect.TypeOf((*MockBoardUsecase)(nil).DeleteCard), userID, boardID, cardID)
}

// DeleteColumn mocks base method.
func (m *MockBoardUsecase) DeleteColumn(userID, boardID, columnID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteColumn", userID, boardID, columnID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteColumn indicates an expected call of DeleteColumn.
func (mr *MockBoardUsecaseMockRecorder) DeleteColumn(userID, boardID, columnID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteColumn", reflect.TypeOf((*MockBoardUsecase)(nil).DeleteColumn), userID, boardID, columnID)
}

// GetBoardContent mocks base method.
func (m *MockBoardUsecase) GetBoardContent(userID, boardID int) (*models.BoardContent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBoardContent", userID, boardID)
	ret0, _ := ret[0].(*models.BoardContent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBoardContent indicates an expected call of GetBoardContent.
func (mr *MockBoardUsecaseMockRecorder) GetBoardContent(userID, boardID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBoardContent", reflect.TypeOf((*MockBoardUsecase)(nil).GetBoardContent), userID, boardID)
}

// GetMembersPermissions mocks base method.
func (m *MockBoardUsecase) GetMembersPermissions(userID, boardID int) ([]models.MemberWithPermissions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMembersPermissions", userID, boardID)
	ret0, _ := ret[0].([]models.MemberWithPermissions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMembersPermissions indicates an expected call of GetMembersPermissions.
func (mr *MockBoardUsecaseMockRecorder) GetMembersPermissions(userID, boardID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMembersPermissions", reflect.TypeOf((*MockBoardUsecase)(nil).GetMembersPermissions), userID, boardID)
}

// GetMyBoards mocks base method.
func (m *MockBoardUsecase) GetMyBoards(userID int) ([]models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMyBoards", userID)
	ret0, _ := ret[0].([]models.Board)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMyBoards indicates an expected call of GetMyBoards.
func (mr *MockBoardUsecaseMockRecorder) GetMyBoards(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMyBoards", reflect.TypeOf((*MockBoardUsecase)(nil).GetMyBoards), userID)
}

// RemoveMember mocks base method.
func (m *MockBoardUsecase) RemoveMember(userID, boardID, memberID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveMember", userID, boardID, memberID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveMember indicates an expected call of RemoveMember.
func (mr *MockBoardUsecaseMockRecorder) RemoveMember(userID, boardID, memberID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveMember", reflect.TypeOf((*MockBoardUsecase)(nil).RemoveMember), userID, boardID, memberID)
}

// UpdateBoard mocks base method.
func (m *MockBoardUsecase) UpdateBoard(userID, boardID int, data models.BoardPutRequest) (*models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBoard", userID, boardID, data)
	ret0, _ := ret[0].(*models.Board)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateBoard indicates an expected call of UpdateBoard.
func (mr *MockBoardUsecaseMockRecorder) UpdateBoard(userID, boardID, data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBoard", reflect.TypeOf((*MockBoardUsecase)(nil).UpdateBoard), userID, boardID, data)
}

// UpdateCard mocks base method.
func (m *MockBoardUsecase) UpdateCard(userID, boardID, cardID int, data *models.CardPutRequest) (*models.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCard", userID, boardID, cardID, data)
	ret0, _ := ret[0].(*models.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateCard indicates an expected call of UpdateCard.
func (mr *MockBoardUsecaseMockRecorder) UpdateCard(userID, boardID, cardID, data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCard", reflect.TypeOf((*MockBoardUsecase)(nil).UpdateCard), userID, boardID, cardID, data)
}

// UpdateColumn mocks base method.
func (m *MockBoardUsecase) UpdateColumn(userID, boardID, columnID int, data *models.ColumnRequest) (*models.Column, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateColumn", userID, boardID, columnID, data)
	ret0, _ := ret[0].(*models.Column)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateColumn indicates an expected call of UpdateColumn.
func (mr *MockBoardUsecaseMockRecorder) UpdateColumn(userID, boardID, columnID, data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateColumn", reflect.TypeOf((*MockBoardUsecase)(nil).UpdateColumn), userID, boardID, columnID, data)
}

// UpdateMemberRole mocks base method.
func (m *MockBoardUsecase) UpdateMemberRole(userID, boardID, memberID int, newRole string) (*models.MemberWithPermissions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMemberRole", userID, boardID, memberID, newRole)
	ret0, _ := ret[0].(*models.MemberWithPermissions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateMemberRole indicates an expected call of UpdateMemberRole.
func (mr *MockBoardUsecaseMockRecorder) UpdateMemberRole(userID, boardID, memberID, newRole any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMemberRole", reflect.TypeOf((*MockBoardUsecase)(nil).UpdateMemberRole), userID, boardID, memberID, newRole)
}

// MockBoardRepo is a mock of BoardRepo interface.
type MockBoardRepo struct {
	ctrl     *gomock.Controller
	recorder *MockBoardRepoMockRecorder
	isgomock struct{}
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
func (m *MockBoardRepo) AddMember(boardID, adderID, memberUserID int) (*models.MemberWithPermissions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMember", boardID, adderID, memberUserID)
	ret0, _ := ret[0].(*models.MemberWithPermissions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddMember indicates an expected call of AddMember.
func (mr *MockBoardRepoMockRecorder) AddMember(boardID, adderID, memberUserID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMember", reflect.TypeOf((*MockBoardRepo)(nil).AddMember), boardID, adderID, memberUserID)
}

// CreateBoard mocks base method.
func (m *MockBoardRepo) CreateBoard(name string, userID int) (*models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBoard", name, userID)
	ret0, _ := ret[0].(*models.Board)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBoard indicates an expected call of CreateBoard.
func (mr *MockBoardRepoMockRecorder) CreateBoard(name, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBoard", reflect.TypeOf((*MockBoardRepo)(nil).CreateBoard), name, userID)
}

// CreateColumn mocks base method.
func (m *MockBoardRepo) CreateColumn(boardId int, title string) (*models.Column, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateColumn", boardId, title)
	ret0, _ := ret[0].(*models.Column)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateColumn indicates an expected call of CreateColumn.
func (mr *MockBoardRepoMockRecorder) CreateColumn(boardId, title any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateColumn", reflect.TypeOf((*MockBoardRepo)(nil).CreateColumn), boardId, title)
}

// CreateNewCard mocks base method.
func (m *MockBoardRepo) CreateNewCard(boardID, columnID int, title string) (*models.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNewCard", boardID, columnID, title)
	ret0, _ := ret[0].(*models.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNewCard indicates an expected call of CreateNewCard.
func (mr *MockBoardRepoMockRecorder) CreateNewCard(boardID, columnID, title any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNewCard", reflect.TypeOf((*MockBoardRepo)(nil).CreateNewCard), boardID, columnID, title)
}

// DeleteBoard mocks base method.
func (m *MockBoardRepo) DeleteBoard(boardId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBoard", boardId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBoard indicates an expected call of DeleteBoard.
func (mr *MockBoardRepoMockRecorder) DeleteBoard(boardId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBoard", reflect.TypeOf((*MockBoardRepo)(nil).DeleteBoard), boardId)
}

// DeleteCard mocks base method.
func (m *MockBoardRepo) DeleteCard(boardID, cardID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCard", boardID, cardID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCard indicates an expected call of DeleteCard.
func (mr *MockBoardRepoMockRecorder) DeleteCard(boardID, cardID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCard", reflect.TypeOf((*MockBoardRepo)(nil).DeleteCard), boardID, cardID)
}

// DeleteColumn mocks base method.
func (m *MockBoardRepo) DeleteColumn(boardID, columnID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteColumn", boardID, columnID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteColumn indicates an expected call of DeleteColumn.
func (mr *MockBoardRepoMockRecorder) DeleteColumn(boardID, columnID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteColumn", reflect.TypeOf((*MockBoardRepo)(nil).DeleteColumn), boardID, columnID)
}

// GetBoard mocks base method.
func (m *MockBoardRepo) GetBoard(boardID int) (*models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBoard", boardID)
	ret0, _ := ret[0].(*models.Board)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBoard indicates an expected call of GetBoard.
func (mr *MockBoardRepoMockRecorder) GetBoard(boardID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBoard", reflect.TypeOf((*MockBoardRepo)(nil).GetBoard), boardID)
}

// GetBoardsForUser mocks base method.
func (m *MockBoardRepo) GetBoardsForUser(userID int) ([]models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBoardsForUser", userID)
	ret0, _ := ret[0].([]models.Board)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBoardsForUser indicates an expected call of GetBoardsForUser.
func (mr *MockBoardRepoMockRecorder) GetBoardsForUser(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBoardsForUser", reflect.TypeOf((*MockBoardRepo)(nil).GetBoardsForUser), userID)
}

// GetCardsForBoard mocks base method.
func (m *MockBoardRepo) GetCardsForBoard(boardID int) ([]models.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCardsForBoard", boardID)
	ret0, _ := ret[0].([]models.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCardsForBoard indicates an expected call of GetCardsForBoard.
func (mr *MockBoardRepoMockRecorder) GetCardsForBoard(boardID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCardsForBoard", reflect.TypeOf((*MockBoardRepo)(nil).GetCardsForBoard), boardID)
}

// GetColumnsForBoard mocks base method.
func (m *MockBoardRepo) GetColumnsForBoard(boardID int) ([]models.Column, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetColumnsForBoard", boardID)
	ret0, _ := ret[0].([]models.Column)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetColumnsForBoard indicates an expected call of GetColumnsForBoard.
func (mr *MockBoardRepoMockRecorder) GetColumnsForBoard(boardID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetColumnsForBoard", reflect.TypeOf((*MockBoardRepo)(nil).GetColumnsForBoard), boardID)
}

// GetMemberPermissions mocks base method.
func (m *MockBoardRepo) GetMemberPermissions(boardID, memberUserID int, getAdderInfo bool) (*models.MemberWithPermissions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMemberPermissions", boardID, memberUserID, getAdderInfo)
	ret0, _ := ret[0].(*models.MemberWithPermissions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMemberPermissions indicates an expected call of GetMemberPermissions.
func (mr *MockBoardRepoMockRecorder) GetMemberPermissions(boardID, memberUserID, getAdderInfo any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMemberPermissions", reflect.TypeOf((*MockBoardRepo)(nil).GetMemberPermissions), boardID, memberUserID, getAdderInfo)
}

// GetMembersWithPermissions mocks base method.
func (m *MockBoardRepo) GetMembersWithPermissions(boardID int) ([]models.MemberWithPermissions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMembersWithPermissions", boardID)
	ret0, _ := ret[0].([]models.MemberWithPermissions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMembersWithPermissions indicates an expected call of GetMembersWithPermissions.
func (mr *MockBoardRepoMockRecorder) GetMembersWithPermissions(boardID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMembersWithPermissions", reflect.TypeOf((*MockBoardRepo)(nil).GetMembersWithPermissions), boardID)
}

// GetUserByNickname mocks base method.
func (m *MockBoardRepo) GetUserByNickname(nickname string) (*models.UserProfile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByNickname", nickname)
	ret0, _ := ret[0].(*models.UserProfile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByNickname indicates an expected call of GetUserByNickname.
func (mr *MockBoardRepoMockRecorder) GetUserByNickname(nickname any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByNickname", reflect.TypeOf((*MockBoardRepo)(nil).GetUserByNickname), nickname)
}

// GetUserProfile mocks base method.
func (m *MockBoardRepo) GetUserProfile(userID int) (*models.UserProfile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserProfile", userID)
	ret0, _ := ret[0].(*models.UserProfile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserProfile indicates an expected call of GetUserProfile.
func (mr *MockBoardRepoMockRecorder) GetUserProfile(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserProfile", reflect.TypeOf((*MockBoardRepo)(nil).GetUserProfile), userID)
}

// RemoveMember mocks base method.
func (m *MockBoardRepo) RemoveMember(boardID, memberUserID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveMember", boardID, memberUserID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveMember indicates an expected call of RemoveMember.
func (mr *MockBoardRepoMockRecorder) RemoveMember(boardID, memberUserID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveMember", reflect.TypeOf((*MockBoardRepo)(nil).RemoveMember), boardID, memberUserID)
}

// SetMemberRole mocks base method.
func (m *MockBoardRepo) SetMemberRole(boardID, memberUserID int, newRole string) (*models.MemberWithPermissions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetMemberRole", boardID, memberUserID, newRole)
	ret0, _ := ret[0].(*models.MemberWithPermissions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetMemberRole indicates an expected call of SetMemberRole.
func (mr *MockBoardRepoMockRecorder) SetMemberRole(boardID, memberUserID, newRole any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMemberRole", reflect.TypeOf((*MockBoardRepo)(nil).SetMemberRole), boardID, memberUserID, newRole)
}

// UpdateBoard mocks base method.
func (m *MockBoardRepo) UpdateBoard(boardID int, data *models.BoardPutRequest) (*models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBoard", boardID, data)
	ret0, _ := ret[0].(*models.Board)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateBoard indicates an expected call of UpdateBoard.
func (mr *MockBoardRepoMockRecorder) UpdateBoard(boardID, data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBoard", reflect.TypeOf((*MockBoardRepo)(nil).UpdateBoard), boardID, data)
}

// UpdateCard mocks base method.
func (m *MockBoardRepo) UpdateCard(boardID, cardID int, data models.CardPutRequest) (*models.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCard", boardID, cardID, data)
	ret0, _ := ret[0].(*models.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateCard indicates an expected call of UpdateCard.
func (mr *MockBoardRepoMockRecorder) UpdateCard(boardID, cardID, data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCard", reflect.TypeOf((*MockBoardRepo)(nil).UpdateCard), boardID, cardID, data)
}

// UpdateColumn mocks base method.
func (m *MockBoardRepo) UpdateColumn(boardID, columnID int, data models.ColumnRequest) (*models.Column, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateColumn", boardID, columnID, data)
	ret0, _ := ret[0].(*models.Column)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateColumn indicates an expected call of UpdateColumn.
func (mr *MockBoardRepoMockRecorder) UpdateColumn(boardID, columnID, data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateColumn", reflect.TypeOf((*MockBoardRepo)(nil).UpdateColumn), boardID, columnID, data)
}
