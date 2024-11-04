package usecase_test

import (
	"RPO_back/internal/models"
	mocks "RPO_back/internal/pkg/board/mocks"
	BoardUsecase "RPO_back/internal/pkg/board/usecase"
	"errors"
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
		// {
		// 	name:    "successful board update by admin",
		// 	userID:  1,
		// 	boardID: 1,
		// 	request: models.BoardPutRequest{NewName: "Updated Board"},
		// 	setupMock: func() {
		// 		mockBoardRepo.EXPECT().GetMemberPermissions(1, 1, false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
		// 		mockBoardRepo.EXPECT().UpdateBoard(1, gomock.Any()).Return(nil).Times(1)
		// 	},
		// 	expectedError: false,
		// 	expectedBoardChecker: func(b *models.Board) bool {
		// 		return b != nil && b.Name == "Updated Board"
		// 	},
		// },
		{
			name:    "successful board update by editor chief",
			userID:  2,
			boardID: 1,
			request: models.BoardPutRequest{NewName: "Updated Board by Editor Chief"},
			setupMock: func() {
				mockBoardRepo.EXPECT().GetMemberPermissions(1, 2, false).Return(&models.MemberWithPermissions{Role: "can_edit"}, nil)
				mockBoardRepo.EXPECT().UpdateBoard(2, gomock.Any()).Return( nil)
			},
			expectedError: true,
			expectedBoardChecker: func(b *models.Board) bool {
				return b != nil && b.Name == "Updated Board by Editor Chief"
			},
		},
		// {
		// 	name:    "permission denied to update board",
		// 	userID:  3,
		// 	boardID: 1,
		// 	request: models.BoardPutRequest{NewName: "Unauthorized Update"},
		// 	setupMock: func() {
		// 		mockBoardRepo.EXPECT().GetMemberPermissions(1, 3, false).Return(&models.MemberPermissions{CanEdit: true}, nil)
		// 	},
		// 	expectedError: true,
		// 	expectedBoardChecker: func(b *models.Board) bool {
		// 		return b == nil
		// 	},
		// },
		// {
		// 	name:    "error fetching permissions",
		// 	userID:  4,
		// 	boardID: 1,
		// 	request: models.BoardPutRequest{NewName: "Error Fetching Permissions"},
		// 	setupMock: func() {
		// 		mockBoardRepo.EXPECT().GetMemberPermissions(1, 4, false).Return(nil, errors.New("error fetching permissions"))
		// 	},
		// 	expectedError: true,
		// 	expectedBoardChecker: func(b *models.Board) bool {
		// 		return b == nil
		// 	},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			updatedBoard, err := boardUsecase.UpdateBoard(tt.userID, tt.boardID, tt.request)
			assert.Equal(t, err != nil, tt.expectedError)
			assert.Equal(t, tt.expectedBoardChecker(updatedBoard), true)
		})
	}
}
