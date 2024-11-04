package usecase_test

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	mocks "RPO_back/internal/pkg/board/mocks"
	BoardUsecase "RPO_back/internal/pkg/board/usecase"
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
		userID               int
		request              models.CreateBoardRequest
		setupMock            func()
		expectedError        bool
		expectedBoardChecker func(*models.Board) bool
	}{
		{
			name:    "successful board creation",
			userID:  1,
			request: models.CreateBoardRequest{Name: "New Board"},
			setupMock: func() {
				mockBoardRepo.EXPECT().CreateBoard("New Board", 1).Return(&models.Board{ID: 1, Name: "New Board"}, nil)
				mockBoardRepo.EXPECT().AddMember(1, 1, 1).Return(nil, nil)
				mockBoardRepo.EXPECT().SetMemberRole(1, 1, "admin").Return(nil, nil)
			},
			expectedError: false,
			expectedBoardChecker: func(b *models.Board) bool {
				return b != nil && b.Name == "New Board"
			},
		},
		{
			name:    "failed to create board",
			userID:  1,
			request: models.CreateBoardRequest{Name: "New Board"},
			setupMock: func() {
				mockBoardRepo.EXPECT().CreateBoard("New Board", 1).Return(nil, errors.New("creation error"))
			},
			expectedError:        true,
			expectedBoardChecker: func(b *models.Board) bool { return b == nil },
		},
		{
			name:    "failed to add member",
			userID:  1,
			request: models.CreateBoardRequest{Name: "New Board"},
			setupMock: func() {
				mockBoardRepo.EXPECT().CreateBoard("New Board", 1).Return(&models.Board{ID: 1, Name: "New Board"}, nil)
				mockBoardRepo.EXPECT().AddMember(1, 1, 1).Return(nil, errors.New("add member error"))
			},
			expectedError:        true,
			expectedBoardChecker: func(b *models.Board) bool { return b == nil },
		},
		{
			name:    "failed to set member role",
			userID:  1,
			request: models.CreateBoardRequest{Name: "New Board"},
			setupMock: func() {
				mockBoardRepo.EXPECT().CreateBoard("New Board", 1).Return(&models.Board{ID: 1, Name: "New Board"}, nil)
				mockBoardRepo.EXPECT().AddMember(1, 1, 1).Return(nil, nil)
				mockBoardRepo.EXPECT().SetMemberRole(1, 1, "admin").Return(nil, errors.New("set role error"))
			},
			expectedError:        true,
			expectedBoardChecker: func(b *models.Board) bool { return b == nil },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			newBoard, err := boardUsecase.CreateNewBoard(tt.userID, tt.request)
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
		userID               int
		boardID              int
		request              models.BoardPutRequest
		setupMock            func()
		expectedError        bool
		expectedBoardChecker func(*models.Board) bool
	}{
		{
			name:    "successful board update by admin",
			userID:  1,
			boardID: 1,
			request: models.BoardPutRequest{NewName: "Updated Board"},
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().UpdateBoard(1, gomock.Any()).Return(&models.Board{Name: "Updated Board"},nil)
			},
			expectedError: false,
			expectedBoardChecker: func(b *models.Board) bool {
				return b != nil && b.Name == "Updated Board"
			},
		},
		{
			name:    "successful board update by editor chief",
			userID:  2,
			boardID: 1,
			request: models.BoardPutRequest{NewName: "Updated Board by Editor Chief"},
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 2, false).Return(&models.MemberWithPermissions{Role: "editor_chief"}, nil)
				mockBoardRepo.EXPECT().UpdateBoard(1, gomock.Any()).Return( &models.Board{Name: "Updated Board by Editor Chief"}, nil)
			},
			expectedError: false,
			expectedBoardChecker: func(b *models.Board) bool {
				return b != nil && b.Name == "Updated Board by Editor Chief"
			},
		},
		{
			name:    "permission denied to update board",
			userID:  3,
			boardID: 1,
			request: models.BoardPutRequest{NewName: "Unauthorized Update"},
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 3, false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
			},
			expectedError: true,
			expectedBoardChecker: func(b *models.Board) bool {
				return b == nil
			},
		},
		{
			name:    "error fetching permissions",
			userID:  4,
			boardID: 1,
			request: models.BoardPutRequest{NewName: "Error Fetching Permissions"},
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 4, false).Return(nil, errors.New("error fetching permissions"))
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

			updatedBoard, err := boardUsecase.UpdateBoard(tt.userID, tt.boardID, tt.request)
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
		userID        int
		boardID       int
		setupMock     func()
		expectedError bool
	}{
		{
			name:    "successful board deletion by admin",
			userID:  1,
			boardID: 1,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().DeleteBoard(1).Return(nil)
			},
			expectedError: false,
		},
		{
			name:    "permission denied to delete board",
			userID:  2,
			boardID: 1,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 2, false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
			},
			expectedError: true,
		},
		{
			name:    "error fetching permissions during board deletion",
			userID:  3,
			boardID: 1,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 3, false).Return(nil, errors.New("error fetching permissions"))
			},
			expectedError: true,
		},
		{
			name:    "error during board deletion",
			userID:  1,
			boardID: 2,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(2, 1, false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().DeleteBoard(2).Return(errors.New("deletion error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := boardUsecase.DeleteBoard(tt.userID, tt.boardID)
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
		name          string
		userID        int
		setupMock     func()
		expectedError bool
		expectedBoards []models.Board
	}{
		{
			name:   "successful retrieval of boards",
			userID: 1,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetBoardsForUser(1).Return([]models.Board{{ID: 1, Name: "Board 1"}, {ID: 2, Name: "Board 2"}}, nil)
			},
			expectedError: false,
			expectedBoards: []models.Board{{ID: 1, Name: "Board 1"}, {ID: 2, Name: "Board 2"}},
		},
		{
			name:   "no boards for user",
			userID: 2,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetBoardsForUser(2).Return([]models.Board{}, nil)
			},
			expectedError: false,
			expectedBoards: []models.Board{},
		},
		{
			name:   "error retrieving boards",
			userID: 3,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetBoardsForUser(3).Return(nil, errors.New("retrieval error"))
			},
			expectedError: true,
			expectedBoards: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			boards, err := boardUsecase.GetMyBoards(tt.userID)
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
		name              string
		ID           	  int
		boardID           int
		setupMock         func()
		expectedError     bool
		expectedData      []models.MemberWithPermissions
	}{
		{
			name:    "successful retrieval of member permissions",
			ID:      1,
			boardID: 1,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().GetMembersWithPermissions(1).Return([]models.MemberWithPermissions{
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
			ID:  2,
			boardID: 1,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 2, false).Return(nil, errors.New("permissions error"))
			},
			expectedError: true,
			expectedData:  nil,
		},
		{
			name:    "error fetching all members' permissions",
			ID:  1,
			boardID: 1,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().GetMembersWithPermissions(1).Return(nil, errors.New("query error"))
			},
			expectedError: true,
			expectedData:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			
			data, err := boardUsecase.GetMembersPermissions(tt.ID, tt.boardID)
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
		ID         int
		boardID        int
		addRequest     *models.AddMemberRequest
		setupMock      func()
		expectedError  bool
		expectedMember *models.MemberWithPermissions
	}{
		{
			name:       "successful addition of member by admin",
			ID:     1,
			boardID:    1,
			addRequest: &models.AddMemberRequest{MemberNickname: "user123"},
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().GetUserByNickname("user123").Return(&models.UserProfile{ID: 2}, nil)
				mockBoardRepo.EXPECT().AddMember(1, 1, 2).Return(&models.MemberWithPermissions{
					User:   &models.UserProfile{ID: 2},
					Role:   "viewer",
				}, nil)
			},
			expectedError: false,
			expectedMember: &models.MemberWithPermissions{
				User:   &models.UserProfile{ID: 2},
				Role:   "viewer",
			},
		},
		{
			name:       "permission denied for non-admin",
			ID:     2,
			boardID:    1,
			addRequest: &models.AddMemberRequest{MemberNickname: "user123"},
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 2, false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
			},
			expectedError:  true,
			expectedMember: nil,
		},
		{
			name:       "error fetching new user's profile",
			ID:     1,
			boardID:    1,
			addRequest: &models.AddMemberRequest{MemberNickname: "user123"},
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().GetUserByNickname("user123").Return(nil, errors.New("user not found"))
			},
			expectedError:  true,
			expectedMember: nil,
		},
		{
			name:       "error adding new user to board",
			ID:     1,
			boardID:    1,
			addRequest: &models.AddMemberRequest{MemberNickname: "user123"},
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().GetUserByNickname("user123").Return(&models.UserProfile{ID: 2}, nil)
				mockBoardRepo.EXPECT().AddMember(1, 1, 2).Return(nil, errors.New("addition error"))
			},
			expectedError:  true,
			expectedMember: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			newMember, err := boardUsecase.AddMember(tt.ID, tt.boardID, tt.addRequest)
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
		name            string
		ID          int
		boardID         int
		memberID        int
		newRole         string
		setupMock       func()
		expectedError   bool
		expectedMember  *models.MemberWithPermissions
	}{
		{
			name:       "successful role update by admin",
			ID:     1,
			boardID:    1,
			memberID:   2,
			newRole:    "editor",
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 2, false).Return(&models.MemberWithPermissions{Role: "viewer"}, nil)
				mockBoardRepo.EXPECT().SetMemberRole(1, 2, "editor").Return(&models.MemberWithPermissions{User: &models.UserProfile{ID: 2}, Role: "editor"}, nil)
			},
			expectedError: false,
			expectedMember: &models.MemberWithPermissions{
				User:   &models.UserProfile{ID: 2},
				Role:   "editor",
			},
		},
		{
			name:       "permission denied for insufficient privileges",
			ID:     3,
			boardID:    1,
			memberID:   2,
			newRole:    "editor",
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 3, false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 2, false).Return(&models.MemberWithPermissions{Role: "viewer"}, nil)
			},
			expectedError: true,
			expectedMember: nil,
		},
		{
			name:       "error fetching updater's permissions",
			ID:     4,
			boardID:    1,
			memberID:   2,
			newRole:    "editor",
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 4, false).Return(nil, errors.New("permissions error"))
			},
			expectedError: true,
			expectedMember: nil,
		},
		{
			name:       "error updating member's role",
			ID:     1,
			boardID:    1,
			memberID:   2,
			newRole:    "editor",
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 2, false).Return(&models.MemberWithPermissions{Role: "viewer"}, nil)
				mockBoardRepo.EXPECT().SetMemberRole(1, 2, "editor").Return(nil, errors.New("update error"))
			},
			expectedError: true,
			expectedMember: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			
			updatedMember, err := boardUsecase.UpdateMemberRole(tt.ID, tt.boardID, tt.memberID, tt.newRole)
			assert.Equal(t, err != nil, tt.expectedError)
			assert.Equal(t, updatedMember, tt.expectedMember)
		})
	}
}

func TestBoardUsecase_RemoveMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBoardRepo := mocks.NewMockBoardRepo(ctrl)
	boardUsecase := BoardUsecase.CreateBoardUsecase(mockBoardRepo)

	tests := []struct {
		name          string
		userID        int
		boardID       int
		memberID      int
		setupMock     func()
		expectedError bool
	}{
		{
			name:     "successful removal by admin",
			userID:   1,
			boardID:  1,
			memberID: 2,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 2, false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
				mockBoardRepo.EXPECT().RemoveMember(1, 2).Return(nil)
			},
			expectedError: false,
		},
		{
			name:     "permission denied for non-admin trying to remove someone",
			userID:   3,
			boardID:  1,
			memberID: 2,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 3, false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 2, false).Return(&models.MemberWithPermissions{Role: "viewer"}, nil)
			},
			expectedError: true,
		},
		{
			name:     "error fetching remover's permissions",
			userID:   4,
			boardID:  1,
			memberID: 2,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 4, false).Return(nil, errors.New("permissions error"))
			},
			expectedError: true,
		},
		{
			name:     "error during member removal",
			userID:   1,
			boardID:  1,
			memberID: 2,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 2, false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
				mockBoardRepo.EXPECT().RemoveMember(1, 2).Return(errors.New("removal error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := boardUsecase.RemoveMember(tt.userID, tt.boardID, tt.memberID)
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
		userID        int
		boardID       int
		setupMock     func()
		expectedError bool
		expectedRole  string
	}{
		{
			name:    "successful content retrieval",
			userID:  1,
			boardID: 1,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().GetCardsForBoard(1).Return([]models.Card{{ID: 1}}, nil)
				mockBoardRepo.EXPECT().GetColumnsForBoard(1).Return([]models.Column{{Id: 1}}, nil)
				mockBoardRepo.EXPECT().GetBoard(1).Return(&models.Board{ID: 1}, nil)
			},
			expectedError: false,
			expectedRole:  "admin",
		},
		{
			name:    "permission error",
			userID:  2,
			boardID: 1,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 2, false).Return(nil, errs.ErrNotPermitted)
			},
			expectedError: true,
		},
		{
			name:    "error getting cards",
			userID:  1,
			boardID: 1,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().GetCardsForBoard(1).Return(nil, errors.New("cards error"))
			},
			expectedError: true,
		},
		{
			name:    "error getting columns",
			userID:  1,
			boardID: 1,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().GetCardsForBoard(1).Return([]models.Card{{ID: 1}}, nil)
				mockBoardRepo.EXPECT().GetColumnsForBoard(1).Return(nil, errors.New("columns error"))
			},
			expectedError: true,
		},
		{
			name:    "error getting board info",
			userID:  1,
			boardID: 1,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
				mockBoardRepo.EXPECT().GetCardsForBoard(1).Return([]models.Card{{ID: 1}}, nil)
				mockBoardRepo.EXPECT().GetColumnsForBoard(1).Return([]models.Column{{Id: 1}}, nil)
				mockBoardRepo.EXPECT().GetBoard(1).Return(nil, errors.New("board error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			content, err := boardUsecase.GetBoardContent(tt.userID, tt.boardID)
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

    cardRequest := &models.CardPutRequest{NewColumnId: 10, NewTitle: "New Card"}

    tests := []struct {
        name          string
        userID        int
        boardID       int
        setupMock     func()
        expectedError bool
        expectedCard  *models.Card
    }{
        {
            name:    "successful card creation",
            userID:  1,
            boardID: 1,
            setupMock: func() {
                mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
                mockBoardRepo.EXPECT().CreateNewCard(1, cardRequest.NewColumnId, cardRequest.NewTitle).Return(&models.Card{ID: 1, Title: "New Card", ColumnID: 10}, nil)
            },
            expectedError: false,
            expectedCard:  &models.Card{ID: 1, Title: "New Card", ColumnID: 10},
        },
        {
            name:    "permission error for viewer",
            userID:  2,
            boardID: 1,
            setupMock: func() {
                mockBoardRepo.EXPECT().GetMemberPermissions(1, 2, false).Return(&models.MemberWithPermissions{Role: "viewer"}, nil)
            },
            expectedError: true,
        },
        {
            name:    "error during card creation",
            userID:  1,
            boardID: 1,
            setupMock: func() {
                mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
                mockBoardRepo.EXPECT().CreateNewCard(1, cardRequest.NewColumnId, cardRequest.NewTitle).Return(nil, errors.New("creation error"))
            },
            expectedError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.setupMock()

            newCard, err := boardUsecase.CreateNewCard(tt.userID, tt.boardID, cardRequest)
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

    cardRequest := &models.CardPutRequest{NewColumnId: 10, NewTitle: "Updated Card"}

    tests := []struct {
        name          string
        userID        int
        boardID       int
        cardID        int
        setupMock     func()
        expectedError bool
        expectedCard  *models.Card
    }{
        {
            name:    "successful card update",
            userID:  1,
            boardID: 1,
            cardID:  1,
            setupMock: func() {
                mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
                mockBoardRepo.EXPECT().UpdateCard(1, 1, *cardRequest).Return(&models.Card{ID: 1, Title: "Updated Card", ColumnID: 10}, nil)
            },
            expectedError: false,
            expectedCard:  &models.Card{ID: 1, Title: "Updated Card", ColumnID: 10},
        },
        {
            name:    "permission error for viewer",
            userID:  2,
            boardID: 1,
            cardID:  1,
            setupMock: func() {
                mockBoardRepo.EXPECT().GetMemberPermissions(1, 2, false).Return(&models.MemberWithPermissions{Role: "viewer"}, nil)
            },
            expectedError: true,
        },
        {
            name:    "error during card update",
            userID:  1,
            boardID: 1,
            cardID:  1,
            setupMock: func() {
                mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
                mockBoardRepo.EXPECT().UpdateCard(1, 1, *cardRequest).Return(nil, errors.New("update error"))
            },
            expectedError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.setupMock()

            updatedCard, err := boardUsecase.UpdateCard(tt.userID, tt.boardID, tt.cardID, cardRequest)
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
        userID        int
        boardID       int
        cardID        int
        setupMock     func()
        expectedError bool
    }{
        {
            name:    "successful card deletion",
            userID:  1,
            boardID: 1,
            cardID:  1,
            setupMock: func() {
                mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
                mockBoardRepo.EXPECT().DeleteCard(1, 1).Return(nil)
            },
            expectedError: false,
        },
        {
            name:    "permission error for viewer",
            userID:  2,
            boardID: 1,
            cardID:  1,
            setupMock: func() {
                mockBoardRepo.EXPECT().GetMemberPermissions(1, 2, false).Return(&models.MemberWithPermissions{Role: "viewer"}, nil)
            },
            expectedError: true,
        },
        {
            name:    "error during card deletion",
            userID:  1,
            boardID: 1,
            cardID:  1,
            setupMock: func() {
                mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
                mockBoardRepo.EXPECT().DeleteCard(1, 1).Return(errors.New("delete error"))
            },
            expectedError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.setupMock()

            err := boardUsecase.DeleteCard(tt.userID, tt.boardID, tt.cardID)
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
		userID        int
		boardID       int
		setupMock     func()
		expectedError bool
		expectedCol   *models.Column
	}{
		{
			name:    "successful column creation",
			userID:  1,
			boardID: 1,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
				mockBoardRepo.EXPECT().CreateColumn(1, columnRequest.NewTitle).Return(&models.Column{Id: 1, Title: "New Column"}, nil)
			},
			expectedError: false,
			expectedCol:   &models.Column{Id: 1, Title: "New Column"},
		},
		{
			name:    "permission error for viewer",
			userID:  2,
			boardID: 1,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 2, false).Return(&models.MemberWithPermissions{Role: "viewer"}, nil)
			},
			expectedError: true,
		},
		{
			name:    "error during column creation",
			userID:  1,
			boardID: 1,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
				mockBoardRepo.EXPECT().CreateColumn(1, columnRequest.NewTitle).Return(nil, errors.New("creation error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			newCol, err := boardUsecase.CreateColumn(tt.userID, tt.boardID, columnRequest)
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
		userID        int
		boardID       int
		columnID      int
		setupMock     func()
		expectedError bool
		expectedCol   *models.Column
	}{
		{
			name:    "successful column update",
			userID:  1,
			boardID: 1,
			columnID: 1,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
				mockBoardRepo.EXPECT().UpdateColumn(1, 1, *columnRequest).Return(&models.Column{Id: 1, Title: "Updated Column"}, nil)
			},
			expectedError: false,
			expectedCol:   &models.Column{Id: 1, Title: "Updated Column"},
		},
		{
			name:    "permission error for viewer",
			userID:  2,
			boardID: 1,
			columnID: 1,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 2, false).Return(&models.MemberWithPermissions{Role: "viewer"}, nil)
			},
			expectedError: true,
		},
		{
			name:    "error during column update",
			userID:  1,
			boardID: 1,
			columnID: 1,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
				mockBoardRepo.EXPECT().UpdateColumn(1, 1, *columnRequest).Return(nil, errors.New("update error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			updatedCol, err := boardUsecase.UpdateColumn(tt.userID, tt.boardID, tt.columnID, columnRequest)
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
		userID        int
		boardID       int
		columnID      int
		setupMock     func()
		expectedError bool
	}{
		{
			name:    "successful column deletion",
			userID:  1,
			boardID: 1,
			columnID: 1,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
				mockBoardRepo.EXPECT().DeleteColumn(1, 1).Return(nil)
			},
			expectedError: false,
		},
		{
			name:    "permission error for viewer",
			userID:  2,
			boardID: 1,
			columnID: 1,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 2, false).Return(&models.MemberWithPermissions{Role: "viewer"}, nil)
			},
			expectedError: true,
		},
		{
			name:    "error during column deletion",
			userID:  1,
			boardID: 1,
			columnID: 1,
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "editor"}, nil)
				mockBoardRepo.EXPECT().DeleteColumn(1, 1).Return(errors.New("delete error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := boardUsecase.DeleteColumn(tt.userID, tt.boardID, tt.columnID)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// func TestBoardUsecase_SetBoardBackground(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockBoardRepo := mocks.NewMockBoardRepo(ctrl)
// 	boardUsecase := BoardUsecase.CreateBoardUsecase(mockBoardRepo)

// 	fileContent := []byte("fake image content")
// 	file := io.NopCloser(bytes.NewReader(fileContent))
// 	fileHeader := &multipart.FileHeader{
// 		Filename: "background.png",
// 		Size:     int64(len(fileContent)),
// 	}

// 	tests := []struct {
// 		name          string
// 		userID        int
// 		boardID       int
// 		setupMock     func()
// 		expectedError bool
// 		expectedBoard *models.Board
// 	}{
// 		{
// 			name:    "successful background set",
// 			userID:  1,
// 			boardID: 1,
// 			setupMock: func() {
// 				mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
// 				mockBoardRepo.EXPECT().SetBoardBackground(1, 1, "png", len(fileContent)).Return("path/to/background.png", nil)
// 				mockBoardRepo.EXPECT().GetBoard(1).Return(&models.Board{ID: 1, Name: "Demo Board"}, nil)
// 			},
// 			expectedError: false,
// 			expectedBoard: &models.Board{ID: 1, Name: "Demo Board"},
// 		},
// 		{
// 			name:    "permission error",
// 			userID:  2,
// 			boardID: 1,
// 			setupMock: func() {
// 				mockBoardRepo.EXPECT().GetMemberPermissions(1, 2, false).Return(&models.MemberWithPermissions{Role: "viewer"}, nil)
// 			},
// 			expectedError: true,
// 		},
// 		{
// 			name:    "error during SetBoardBackground",
// 			userID:  1,
// 			boardID: 1,
// 			setupMock: func() {
// 				mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
// 				mockBoardRepo.EXPECT().SetBoardBackground(1, 1, "png", len(fileContent)).Return("", errors.New("background set error"))
// 			},
// 			expectedError: true,
// 		},
// 		{
// 			name:    "file creation error",
// 			userID:  1,
// 			boardID: 1,
// 			setupMock: func() {
// 				mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
// 				mockBoardRepo.EXPECT().SetBoardBackground(1, 1, "png", len(fileContent)).Return("path/to/background.png", nil)
// 			},
// 			expectedError: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.setupMock()

// 			updatedBoard, err := boardUsecase.SetBoardBackground(tt.userID, tt.boardID, &file, fileHeader)
// 			if tt.expectedError {
// 				assert.Nil(t, updatedBoard)
// 				assert.Error(t, err)
// 			} else {
// 				assert.NotNil(t, updatedBoard)
// 				assert.NoError(t, err)
// 				assert.Equal(t, tt.expectedBoard, updatedBoard)
// 			}
// 		})
// 	}
// }
