package usecase_test

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	mocks "RPO_back/internal/pkg/board/mocks"
	BoardUsecase "RPO_back/internal/pkg/board/usecase"
	"RPO_back/internal/pkg/utils/misc"
	"context"
	"errors"
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBoardUsecase_CreateNewBoard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBoardRepo := mocks.NewMockBoardRepo(ctrl)
	boardUsecase := BoardUsecase.CreateBoardUsecase(mockBoardRepo)

	tests := []struct {
		name                 string
		userID               int64
		request              models.BoardRequest
		setupMock            func()
		expectedError        bool
		expectedBoardChecker func(*models.Board) bool
	}{
		{
			name:    "successful board creation",
			userID:  int64(1),
			request: models.BoardRequest{NewName: "New Board"},
			setupMock: func() {
				mockBoardRepo.EXPECT().CreateBoard(gomock.Any(), "New Board", 1).Return(&models.Board{ID: int64(1), Name: "New Board"}, nil)
				mockBoardRepo.EXPECT().AddMember(gomock.Any(), int64(1), int64(1), 1).Return(nil, nil)
				mockBoardRepo.EXPECT().SetMemberRole(gomock.Any(), int64(1), int64(1), int64(1), "admin").Return(nil, nil)
			},
			expectedError: false,
			expectedBoardChecker: func(b *models.Board) bool {
				return b != nil && b.Name == "New Board"
			},
		},
		{
			name:    "failed to create board",
			userID:  int64(1),
			request: models.BoardRequest{NewName: "New Board"},
			setupMock: func() {
				mockBoardRepo.EXPECT().CreateBoard(gomock.Any(), "New Board", 1).Return(nil, errors.New("creation error"))
			},
			expectedError:        true,
			expectedBoardChecker: func(b *models.Board) bool { return b == nil },
		},
		{
			name:    "failed to add member",
			userID:  int64(1),
			request: models.BoardRequest{NewName: "New Board"},
			setupMock: func() {
				mockBoardRepo.EXPECT().CreateBoard(gomock.Any(), "New Board", 1).Return(&models.Board{ID: int64(1), Name: "New Board"}, nil)
				mockBoardRepo.EXPECT().AddMember(gomock.Any(), int64(1), int64(1), 1).Return(nil, errors.New("add member error"))
			},
			expectedError:        true,
			expectedBoardChecker: func(b *models.Board) bool { return b == nil },
		},
		{
			name:    "failed to set member role",
			userID:  int64(1),
			request: models.BoardRequest{NewName: "New Board"},
			setupMock: func() {
				mockBoardRepo.EXPECT().CreateBoard(gomock.Any(), "New Board", 1).Return(&models.Board{ID: int64(1), Name: "New Board"}, nil)
				mockBoardRepo.EXPECT().AddMember(gomock.Any(), int64(1), int64(1), 1).Return(nil, nil)
				mockBoardRepo.EXPECT().SetMemberRole(gomock.Any(), int64(1), int64(1), int64(1), "admin").Return(nil, errors.New("set role error"))
			},
			expectedError:        true,
			expectedBoardChecker: func(b *models.Board) bool { return b == nil },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			newBoard, err := boardUsecase.CreateNewBoard(context.Background(), int64(tt.userID), tt.request)
			assert.Equal(t, err != nil, tt.expectedError)
			assert.Equal(t, tt.expectedBoardChecker(newBoard), true)
		})
	}
}

func TestBoardUsecase_UpdateBoard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBoardRepo := mocks.NewMockBoardRepo(ctrl)
	boardUsecase := BoardUsecase.CreateBoardUsecase(mockBoardRepo)

	tests := []struct {
		name                 string
		userID               int64
		boardID              int64
		request              models.BoardRequest
		setupMock            func()
		expectedError        bool
		expectedBoardChecker func(*models.Board) bool
	}{
		{
			name:    "successful board update by admin",
			userID:  int64(1),
			boardID: int64(1),
			request: models.BoardRequest{NewName: "Updated Board"},
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().UpdateBoard(gomock.Any(), int64(1), int64(1), gomock.Any()).Return(&models.Board{Name: "Updated Board"}, nil)
			},
			expectedError: false,
			expectedBoardChecker: func(b *models.Board) bool {
				return b != nil && b.Name == "Updated Board"
			},
		},
		{
			name:    "successful board update by editor chief",
			userID:  2,
			boardID: int64(1),
			request: models.BoardRequest{NewName: "Updated Board by Editor Chief"},
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), 2, false).Return(&models.MemberWithPermissions{Role: "editor_chief"}, nil)
				mockBoardRepo.EXPECT().UpdateBoard(gomock.Any(), int64(1), int64(1), gomock.Any()).Return(&models.Board{Name: "Updated Board by Editor Chief"}, nil)
			},
			expectedError: false,
			expectedBoardChecker: func(b *models.Board) bool {
				return b != nil && b.Name == "Updated Board by Editor Chief"
			},
		},
		{
			name:    "permission denied to update board",
			userID:  3,
			boardID: int64(1),
			request: models.BoardRequest{NewName: "Unauthorized Update"},
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), 3, false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
			},
			expectedError: true,
			expectedBoardChecker: func(b *models.Board) bool {
				return b == nil
			},
		},
		{
			name:    "error fetching permissions",
			userID:  4,
			boardID: int64(1),
			request: models.BoardRequest{NewName: "Error Fetching Permissions"},
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), 4, false).Return(nil, errors.New("error fetching permissions"))
			},
			expectedError: true,
			expectedBoardChecker: func(b *models.Board) bool {
				return b == nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			updatedBoard, err := boardUsecase.UpdateBoard(context.Background(), int64(tt.userID), int64(tt.boardID), tt.request)
			fmt.Println("Error: ", err)
			assert.Equal(t, err != nil, tt.expectedError)
			assert.Equal(t, tt.expectedBoardChecker(updatedBoard), true)
		})
	}
}

func TestBoardUsecase_DeleteBoard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBoardRepo := mocks.NewMockBoardRepo(ctrl)
	boardUsecase := BoardUsecase.CreateBoardUsecase(mockBoardRepo)

	tests := []struct {
		name          string
		userID        int64
		boardID       int64
		setupMock     func()
		expectedError bool
	}{
		{
			name:    "successful board deletion by admin",
			userID:  int64(1),
			boardID: int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().DeleteBoard(gomock.Any(), 1).Return(nil)
			},
			expectedError: false,
		},
		{
			name:    "permission denied to delete board",
			userID:  2,
			boardID: int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), 2, false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
			},
			expectedError: true,
		},
		{
			name:    "error fetching permissions during board deletion",
			userID:  3,
			boardID: int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), 3, false).Return(nil, errors.New("error fetching permissions"))
			},
			expectedError: true,
		},
		{
			name:    "error during board deletion",
			userID:  int64(1),
			boardID: 2,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), 2, int64(1), false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().DeleteBoard(gomock.Any(), 2).Return(errors.New("deletion error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := boardUsecase.DeleteBoard(context.Background(), int64(tt.userID), int64(tt.boardID))
			assert.Equal(t, err != nil, tt.expectedError)
		})
	}
}

func TestBoardUsecase_GetMyBoards(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBoardRepo := mocks.NewMockBoardRepo(ctrl)
	boardUsecase := BoardUsecase.CreateBoardUsecase(mockBoardRepo)

	tests := []struct {
		name           string
		userID         int64
		setupMock      func()
		expectedError  bool
		expectedBoards []models.Board
	}{
		{
			name:   "successful retrieval of boards",
			userID: int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetBoardsForUser(gomock.Any(), 1).Return([]models.Board{{ID: int64(1), Name: "Board 1"}, {ID: 2, Name: "Board 2"}}, nil)
			},
			expectedError:  false,
			expectedBoards: []models.Board{{ID: int64(1), Name: "Board 1"}, {ID: 2, Name: "Board 2"}},
		},
		{
			name:   "no boards for user",
			userID: 2,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetBoardsForUser(gomock.Any(), 2).Return([]models.Board{}, nil)
			},
			expectedError:  false,
			expectedBoards: []models.Board{},
		},
		{
			name:   "error retrieving boards",
			userID: 3,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetBoardsForUser(gomock.Any(), 3).Return(nil, errors.New("retrieval error"))
			},
			expectedError:  true,
			expectedBoards: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			boards, err := boardUsecase.GetMyBoards(context.Background(), int64(tt.userID))
			assert.Equal(t, err != nil, tt.expectedError)
			assert.Equal(t, boards, tt.expectedBoards)
		})
	}
}

func TestBoardUsecase_GetMembersPermissions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBoardRepo := mocks.NewMockBoardRepo(ctrl)
	boardUsecase := BoardUsecase.CreateBoardUsecase(mockBoardRepo)

	tests := []struct {
		name          string
		ID            int64
		boardID       int64
		setupMock     func()
		expectedError bool
		expectedData  []models.MemberWithPermissions
	}{
		{
			name:    "successful retrieval of member permissions",
			ID:      int64(1),
			boardID: int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().GetMembersWithPermissions(gomock.Any(), int64(1), 1).Return([]models.MemberWithPermissions{
					{User: &models.UserProfile{ID: 1}, Role: "admin"},
					{User: &models.UserProfile{ID: 2}, Role: "editor"},
				}, nil)
			},
			expectedError: false,
			expectedData: []models.MemberWithPermissions{
				{User: &models.UserProfile{ID: 1}, Role: "admin"},
				{User: &models.UserProfile{ID: 2}, Role: "editor"},
			},
		},
		{
			name:    "error fetching user's own permissions",
			ID:      2,
			boardID: int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), 2, false).Return(nil, errors.New("permissions error"))
			},
			expectedError: true,
			expectedData:  nil,
		},
		{
			name:    "error fetching all members' permissions",
			ID:      int64(1),
			boardID: int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().GetMembersWithPermissions(gomock.Any(), int64(1), 1).Return(nil, errors.New("query error"))
			},
			expectedError: true,
			expectedData:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			data, err := boardUsecase.GetMembersPermissions(context.Background(), int64(tt.ID), int64(tt.boardID))
			assert.Equal(t, err != nil, tt.expectedError)
			assert.Equal(t, data, tt.expectedData)
		})
	}
}

func TestBoardUsecase_AddMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBoardRepo := mocks.NewMockBoardRepo(ctrl)
	boardUsecase := BoardUsecase.CreateBoardUsecase(mockBoardRepo)

	tests := []struct {
		name           string
		ID             int64
		boardID        int64
		addRequest     *models.AddMemberRequest
		setupMock      func()
		expectedError  bool
		expectedMember *models.MemberWithPermissions
	}{
		{
			name:       "successful addition of member by admin",
			ID:         int64(1),
			boardID:    int64(1),
			addRequest: &models.AddMemberRequest{MemberNickname: "user123"},
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().GetUserByNickname(gomock.Any(), "user123").Return(&models.UserProfile{ID: 1}, nil)
				mockBoardRepo.EXPECT().AddMember(gomock.Any(), int64(1), int64(1), int64(1)).Return(&models.MemberWithPermissions{
					User: &models.UserProfile{ID: 1},
					Role: "viewer",
				}, nil)
			},
			expectedError: false,
			expectedMember: &models.MemberWithPermissions{
				User: &models.UserProfile{ID: 1},
				Role: "viewer",
			},
		},
		{
			name:       "permission denied for non-admin",
			ID:         1,
			boardID:    int64(1),
			addRequest: &models.AddMemberRequest{MemberNickname: "user123"},
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
			},
			expectedError:  true,
			expectedMember: nil,
		},
		{
			name:       "error fetching new user's profile",
			ID:         int64(1),
			boardID:    int64(1),
			addRequest: &models.AddMemberRequest{MemberNickname: "user123"},
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().GetUserByNickname(gomock.Any(), "user123").Return(nil, errors.New("user not found"))
			},
			expectedError:  true,
			expectedMember: nil,
		},
		{
			name:       "error adding new user to board",
			ID:         int64(1),
			boardID:    int64(1),
			addRequest: &models.AddMemberRequest{MemberNickname: "user123"},
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().GetUserByNickname(gomock.Any(), "user123").Return(&models.UserProfile{ID: 1}, nil)
				mockBoardRepo.EXPECT().AddMember(gomock.Any(), int64(1), int64(1), int64(1)).Return(nil, errors.New("addition error"))
			},
			expectedError:  true,
			expectedMember: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			newMember, err := boardUsecase.AddMember(context.Background(), int64(tt.ID), int64(tt.boardID), tt.addRequest)
			assert.Equal(t, err != nil, tt.expectedError)
			assert.Equal(t, newMember, tt.expectedMember)
		})
	}
}

func TestBoardUsecase_UpdateMemberRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBoardRepo := mocks.NewMockBoardRepo(ctrl)
	boardUsecase := BoardUsecase.CreateBoardUsecase(mockBoardRepo)

	tests := []struct {
		name           string
		ID             int64
		boardID        int64
		memberID       int64
		newRole        string
		setupMock      func()
		expectedError  bool
		expectedMember *models.MemberWithPermissions
	}{
		{
			name:     "successful role update by admin",
			ID:       1,
			boardID:  1,
			memberID: 1,
			newRole:  "editor",
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "viewer"}, nil)
				mockBoardRepo.EXPECT().SetMemberRole(gomock.Any(), int64(1), int64(1), int64(1), "editor").Return(&models.MemberWithPermissions{User: &models.UserProfile{ID: 2}, Role: "editor"}, nil)
			},
			expectedError: false,
			expectedMember: &models.MemberWithPermissions{
				User: &models.UserProfile{ID: 1},
				Role: "editor",
			},
		},
		{
			name:     "permission denied for insufficient privileges",
			ID:       1,
			boardID:  1,
			memberID: 1,
			newRole:  "editor",
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "viewer"}, nil)
			},
			expectedError:  true,
			expectedMember: nil,
		},
		{
			name:     "error fetching updater's permissions",
			ID:       1,
			boardID:  1,
			memberID: 1,
			newRole:  "editor",
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(nil, errors.New("permissions error"))
			},
			expectedError:  true,
			expectedMember: nil,
		},
		{
			name:     "error updating member's role",
			ID:       1,
			boardID:  1,
			memberID: 1,
			newRole:  "editor",
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "viewer"}, nil)
				mockBoardRepo.EXPECT().SetMemberRole(gomock.Any(), int64(1), int64(1), int64(1), "editor").Return(nil, errors.New("update error"))
			},
			expectedError:  true,
			expectedMember: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			_, err := boardUsecase.UpdateMemberRole(context.Background(), int64(tt.ID), int64(tt.boardID), int64(tt.memberID), tt.newRole)
			assert.Equal(t, err != nil, tt.expectedError)
		})
	}
}

func TestBoardUsecase_GetBoardContent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBoardRepo := mocks.NewMockBoardRepo(ctrl)
	boardUsecase := BoardUsecase.CreateBoardUsecase(mockBoardRepo)

	tests := []struct {
		name          string
		userID        int64
		boardID       int64
		setupMock     func()
		expectedError bool
		expectedRole  string
	}{
		{
			name:    "successful content retrieval",
			userID:  int64(1),
			boardID: int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().GetCardsForBoard(gomock.Any(), int64(1)).Return([]models.Card{{ID: int64(1)}}, nil)
				mockBoardRepo.EXPECT().GetColumnsForBoard(gomock.Any(), int64(1)).Return([]models.Column{{ID: int64(1)}}, nil)
				mockBoardRepo.EXPECT().GetBoard(gomock.Any(), int64(1), int64(1)).Return(&models.Board{ID: int64(1)}, nil)
			},
			expectedError: false,
			expectedRole:  "admin",
		},
		{
			name:    "permission error",
			userID:  2,
			boardID: int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), 2, false).Return(nil, errs.ErrNotPermitted)
			},
			expectedError: true,
		},
		{
			name:    "error getting cards",
			userID:  int64(1),
			boardID: int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().GetCardsForBoard(gomock.Any(), int64(1)).Return(nil, errors.New("cards error"))
			},
			expectedError: true,
		},
		{
			name:    "error getting columns",
			userID:  int64(1),
			boardID: int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().GetCardsForBoard(gomock.Any(), int64(1)).Return([]models.Card{{ID: int64(1)}}, nil)
				mockBoardRepo.EXPECT().GetColumnsForBoard(gomock.Any(), int64(1)).Return(nil, errors.New("columns error"))
			},
			expectedError: true,
		},
		{
			name:    "error getting board info",
			userID:  int64(1),
			boardID: int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().GetCardsForBoard(gomock.Any(), int64(1)).Return([]models.Card{{ID: int64(1)}}, nil)
				mockBoardRepo.EXPECT().GetColumnsForBoard(gomock.Any(), int64(1)).Return([]models.Column{{ID: int64(1)}}, nil)
				mockBoardRepo.EXPECT().GetBoard(gomock.Any(), int64(1), int64(1)).Return(nil, errors.New("board error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			content, err := boardUsecase.GetBoardContent(context.Background(), int64(tt.userID), int64(tt.boardID))
			if tt.expectedError {
				assert.Nil(t, content)
				assert.Error(t, err)
			} else {
				assert.NotNil(t, content)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRole, content.MyRole)
			}
		})
	}
}

func TestBoardUsecase_CreateNewCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBoardRepo := mocks.NewMockBoardRepo(ctrl)
	boardUsecase := BoardUsecase.CreateBoardUsecase(mockBoardRepo)

	cardRequest := &models.CardPostRequest{Title: misc.StringPtr("New Card")}

	tests := []struct {
		name          string
		userID        int64
		boardID       int64
		setupMock     func()
		expectedError bool
		expectedCard  *models.Card
	}{
		{
			name:    "successful card creation",
			userID:  int64(1),
			boardID: int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
				mockBoardRepo.EXPECT().CreateNewCard(gomock.Any(), int64(1), cardRequest.Title).Return(&models.Card{ID: int64(1), Title: "New Card", ColumnID: 10}, nil)
			},
			expectedError: false,
			expectedCard:  &models.Card{ID: int64(1), Title: "New Card", ColumnID: 10},
		},
		{
			name:    "permission error for viewer",
			userID:  2,
			boardID: int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), 2, false).Return(&models.MemberWithPermissions{Role: "viewer"}, nil)
			},
			expectedError: true,
		},
		{
			name:    "error during card creation",
			userID:  int64(1),
			boardID: int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
				mockBoardRepo.EXPECT().CreateNewCard(gomock.Any(), int64(1), cardRequest.Title).Return(nil, errors.New("creation error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			newCard, err := boardUsecase.CreateNewCard(context.Background(), int64(tt.userID), int64(tt.boardID), cardRequest)
			if tt.expectedError {
				assert.Nil(t, newCard)
				assert.Error(t, err)
			} else {
				assert.NotNil(t, newCard)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCard, newCard)
			}
		})
	}
}

func TestBoardUsecase_UpdateCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBoardRepo := mocks.NewMockBoardRepo(ctrl)
	boardUsecase := BoardUsecase.CreateBoardUsecase(mockBoardRepo)

	cardRequest := &models.CardPatchRequest{NewTitle: misc.StringPtr("Updated Card")}

	tests := []struct {
		name          string
		userID        int64
		boardID       int64
		cardID        int64
		setupMock     func()
		expectedError bool
		expectedCard  *models.Card
	}{
		{
			name:    "successful card update",
			userID:  int64(1),
			boardID: int64(1),
			cardID:  int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
				mockBoardRepo.EXPECT().UpdateCard(gomock.Any(), int64(1), *cardRequest).Return(&models.Card{ID: int64(1), Title: "Updated Card", ColumnID: 10}, nil)
			},
			expectedError: false,
			expectedCard:  &models.Card{ID: int64(1), Title: "Updated Card", ColumnID: 10},
		},
		{
			name:    "permission error for viewer",
			userID:  2,
			boardID: int64(1),
			cardID:  int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), 2, false).Return(&models.MemberWithPermissions{Role: "viewer"}, nil)
			},
			expectedError: true,
		},
		{
			name:    "error during card update",
			userID:  int64(1),
			boardID: int64(1),
			cardID:  int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
				mockBoardRepo.EXPECT().UpdateCard(gomock.Any(), int64(1), *cardRequest).Return(nil, errors.New("update error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			updatedCard, err := boardUsecase.UpdateCard(context.Background(), int64(tt.userID), int64(tt.cardID), cardRequest)
			if tt.expectedError {
				assert.Nil(t, updatedCard)
				assert.Error(t, err)
			} else {
				assert.NotNil(t, updatedCard)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCard, updatedCard)
			}
		})
	}
}

func TestBoardUsecase_DeleteCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBoardRepo := mocks.NewMockBoardRepo(ctrl)
	boardUsecase := BoardUsecase.CreateBoardUsecase(mockBoardRepo)

	tests := []struct {
		name          string
		userID        int64
		boardID       int64
		cardID        int64
		setupMock     func()
		expectedError bool
	}{
		{
			name:    "successful card deletion",
			userID:  int64(1),
			boardID: int64(1),
			cardID:  int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
				mockBoardRepo.EXPECT().DeleteCard(gomock.Any(), 1).Return(nil)
			},
			expectedError: false,
		},
		{
			name:    "permission error for viewer",
			userID:  2,
			boardID: int64(1),
			cardID:  int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), 2, false).Return(&models.MemberWithPermissions{Role: "viewer"}, nil)
			},
			expectedError: true,
		},
		{
			name:    "error during card deletion",
			userID:  int64(1),
			boardID: int64(1),
			cardID:  int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
				mockBoardRepo.EXPECT().DeleteCard(gomock.Any(), 1).Return(errors.New("delete error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := boardUsecase.DeleteCard(context.Background(), int64(tt.userID), int64(tt.cardID))
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBoardUsecase_CreateColumn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBoardRepo := mocks.NewMockBoardRepo(ctrl)
	boardUsecase := BoardUsecase.CreateBoardUsecase(mockBoardRepo)

	columnRequest := &models.ColumnRequest{NewTitle: "New Column"}

	tests := []struct {
		name          string
		userID        int64
		boardID       int64
		setupMock     func()
		expectedError bool
		expectedCol   *models.Column
	}{
		{
			name:    "successful column creation",
			userID:  int64(1),
			boardID: int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
				mockBoardRepo.EXPECT().CreateColumn(gomock.Any(), int64(1), columnRequest.NewTitle).Return(&models.Column{ID: int64(1), Title: "New Column"}, nil)
			},
			expectedError: false,
			expectedCol:   &models.Column{ID: int64(1), Title: "New Column"},
		},
		{
			name:    "permission error for viewer",
			userID:  2,
			boardID: int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), 2, false).Return(&models.MemberWithPermissions{Role: "viewer"}, nil)
			},
			expectedError: true,
		},
		{
			name:    "error during column creation",
			userID:  int64(1),
			boardID: int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
				mockBoardRepo.EXPECT().CreateColumn(gomock.Any(), int64(1), columnRequest.NewTitle).Return(nil, errors.New("creation error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			newCol, err := boardUsecase.CreateColumn(context.Background(), int64(tt.userID), int64(tt.boardID), columnRequest)
			if tt.expectedError {
				assert.Nil(t, newCol)
				assert.Error(t, err)
			} else {
				assert.NotNil(t, newCol)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCol, newCol)
			}
		})
	}
}

func TestBoardUsecase_UpdateColumn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBoardRepo := mocks.NewMockBoardRepo(ctrl)
	boardUsecase := BoardUsecase.CreateBoardUsecase(mockBoardRepo)

	columnRequest := &models.ColumnRequest{NewTitle: "Updated Column"}

	tests := []struct {
		name          string
		userID        int64
		boardID       int64
		columnID      int64
		setupMock     func()
		expectedError bool
		expectedCol   *models.Column
	}{
		{
			name:     "successful column update",
			userID:   int64(1),
			boardID:  int64(1),
			columnID: int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
				mockBoardRepo.EXPECT().UpdateColumn(gomock.Any(), int64(1), *columnRequest).Return(&models.Column{ID: int64(1), Title: "Updated Column"}, nil)
			},
			expectedError: false,
			expectedCol:   &models.Column{ID: int64(1), Title: "Updated Column"},
		},
		{
			name:     "permission error for viewer",
			userID:   2,
			boardID:  int64(1),
			columnID: int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), 2, false).Return(&models.MemberWithPermissions{Role: "viewer"}, nil)
			},
			expectedError: true,
		},
		{
			name:     "error during column update",
			userID:   int64(1),
			boardID:  int64(1),
			columnID: int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
				mockBoardRepo.EXPECT().UpdateColumn(gomock.Any(), int64(1), *columnRequest).Return(nil, errors.New("update error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			updatedCol, err := boardUsecase.UpdateColumn(context.Background(), int64(tt.userID), int64(tt.columnID), columnRequest)
			if tt.expectedError {
				assert.Nil(t, updatedCol)
				assert.Error(t, err)
			} else {
				assert.NotNil(t, updatedCol)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCol, updatedCol)
			}
		})
	}
}

func TestBoardUsecase_DeleteColumn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBoardRepo := mocks.NewMockBoardRepo(ctrl)
	boardUsecase := BoardUsecase.CreateBoardUsecase(mockBoardRepo)

	tests := []struct {
		name          string
		userID        int64
		boardID       int64
		columnID      int64
		setupMock     func()
		expectedError bool
	}{
		{
			name:     "successful column deletion",
			userID:   int64(1),
			boardID:  int64(1),
			columnID: int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
				mockBoardRepo.EXPECT().DeleteColumn(gomock.Any(), 1).Return(nil)
			},
			expectedError: false,
		},
		{
			name:     "permission error for viewer",
			userID:   2,
			boardID:  int64(1),
			columnID: int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), 2, false).Return(&models.MemberWithPermissions{Role: "viewer"}, nil)
			},
			expectedError: true,
		},
		{
			name:     "error during column deletion",
			userID:   int64(1),
			boardID:  int64(1),
			columnID: int64(1),
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(gomock.Any(), int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
				mockBoardRepo.EXPECT().DeleteColumn(gomock.Any(), 1).Return(errors.New("delete error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := boardUsecase.DeleteColumn(context.Background(), int64(tt.userID), int64(tt.columnID))
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
