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
	// authUsecase := AuthUsecase.CreateAuthUsecase(mockAuthRepo)

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
