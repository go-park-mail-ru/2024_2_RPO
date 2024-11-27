package usecase

import (
	mocks "RPO_back/internal/pkg/board/mocks"
	"context"
	"testing"
	"time"

	"RPO_back/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateBoard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockBoardRepo(ctrl)
	uc := CreateBoardUsecase(mockRepo)
	ctx := context.Background()
	data := models.BoardRequest{}
	mockRepo.EXPECT().GetMemberPermissions(ctx, int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
	mockRepo.EXPECT().UpdateBoard(ctx, int64(1), int64(1), &data).Return(&models.Board{ID: 1}, nil)
	_, err := uc.UpdateBoard(ctx, 1, 1, data)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestDeleteBoard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBoardRepo(ctrl)
	ctx := context.Background()
	userID := int64(1)
	boardID := int64(1)

	mockRepo.EXPECT().GetMemberPermissions(ctx, boardID, userID, false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
	mockRepo.EXPECT().DeleteBoard(ctx, boardID).Return(nil)

	uc := CreateBoardUsecase(mockRepo)
	err := uc.DeleteBoard(ctx, userID, boardID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestGetMembersPermissions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBoardRepo(ctrl)
	uc := CreateBoardUsecase(mockRepo)

	ctx := context.Background()
	userID := int64(1)
	boardID := int64(1)

	mockRepo.EXPECT().GetMemberPermissions(ctx, boardID, userID, false).Return(nil, nil)
	mockRepo.EXPECT().GetMembersWithPermissions(ctx, boardID, userID).Return([]models.MemberWithPermissions{}, nil)

	_, err := uc.GetMembersPermissions(ctx, userID, boardID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAddMember_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	boardRepo := mocks.NewMockBoardRepo(ctrl)
	uc := CreateBoardUsecase(boardRepo)
	ctx := context.Background()
	userID := int64(1)
	boardID := int64(1)
	addRequest := &models.AddMemberRequest{MemberNickname: "nickname"}

	boardRepo.EXPECT().
		GetMemberPermissions(ctx, boardID, userID, false).
		Return(&models.MemberWithPermissions{Role: "admin"}, nil)
	boardRepo.EXPECT().
		GetUserByNickname(ctx, "nickname").
		Return(&models.UserProfile{ID: 1}, nil)
	boardRepo.EXPECT().
		AddMember(ctx, boardID, userID, int64(1)).
		Return(&models.MemberWithPermissions{User: &models.UserProfile{ID: 1}, Role: "viewer"}, nil)

	_, err := uc.AddMember(ctx, userID, boardID, addRequest)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestUpdateMemberRole_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	boardRepo := mocks.NewMockBoardRepo(ctrl)
	uc := CreateBoardUsecase(boardRepo)
	ctx := context.Background()
	userID := int64(1)
	boardID := int64(1)
	memberID := int64(1)
	newRole := "editor"

	boardRepo.EXPECT().GetMemberPermissions(ctx, boardID, userID, false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
	boardRepo.EXPECT().GetMemberPermissions(ctx, boardID, memberID, false).Return(&models.MemberWithPermissions{Role: "member"}, nil)
	boardRepo.EXPECT().SetMemberRole(ctx, userID, boardID, memberID, newRole).Return(&models.MemberWithPermissions{Role: newRole}, nil)

	_, err := uc.UpdateMemberRole(ctx, userID, boardID, memberID, newRole)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestRemoveMember_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	boardRepo := mocks.NewMockBoardRepo(ctrl)
	uc := CreateBoardUsecase(boardRepo)
	ctx := context.Background()
	boardRepo.EXPECT().GetMemberPermissions(ctx, int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
	boardRepo.EXPECT().GetMemberPermissions(ctx, int64(1), int64(1), false).Return(&models.MemberWithPermissions{Role: "user"}, nil)
	boardRepo.EXPECT().RemoveMember(ctx, int64(1), int64(1)).Return(nil)
	err := uc.RemoveMember(ctx, int64(1), int64(1), int64(1))
	if err != nil {
		t.Fatal(err)
	}
}
func TestGetBoardContent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	boardRepo := mocks.NewMockBoardRepo(ctrl)
	uc := CreateBoardUsecase(boardRepo)

	userID := int64(1)
	boardID := int64(1)
	ctx := context.Background()

	boardRepo.EXPECT().GetMemberPermissions(ctx, boardID, userID, false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
	boardRepo.EXPECT().GetCardsForBoard(ctx, boardID).Return([]models.Card{}, nil)
	boardRepo.EXPECT().GetColumnsForBoard(ctx, boardID).Return([]models.Column{}, nil)
	boardRepo.EXPECT().GetBoard(ctx, boardID, userID).Return(&models.Board{}, nil)

	_, err := uc.GetBoardContent(ctx, userID, boardID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestCreateNewCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	boardRepo := mocks.NewMockBoardRepo(ctrl)
	uc := CreateBoardUsecase(boardRepo)
	ctx := context.Background()
	userID := int64(1)
	boardID := int64(1)
	columnID := int64(1)
	title := "Test Title"
	data := &models.CardPostRequest{
		ColumnID: &columnID,
		Title:    &title,
	}
	boardRepo.EXPECT().GetMemberPermissions(ctx, boardID, userID, false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
	boardRepo.EXPECT().CreateNewCard(ctx, columnID, title).Return(&models.Card{
		ID:        1,
		Title:     title,
		ColumnID:  columnID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil)
	_, err := uc.CreateNewCard(ctx, userID, boardID, data)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	boardRepo := mocks.NewMockBoardRepo(ctrl)
	uc := CreateBoardUsecase(boardRepo)

	ctx := context.Background()
	userID := int64(1)
	cardID := int64(1)
	data := &models.CardPatchRequest{}

	boardRepo.EXPECT().GetMemberFromCard(ctx, userID, cardID).Return("admin", int64(1), nil)
	boardRepo.EXPECT().UpdateCard(ctx, cardID, *data).Return(&models.Card{
		ID:        1,
		Title:     "Updated Title",
		ColumnID:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil)

	_, err := uc.UpdateCard(ctx, userID, cardID, data)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDeleteCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	boardRepo := mocks.NewMockBoardRepo(ctrl)
	uc := CreateBoardUsecase(boardRepo)

	ctx := context.Background()
	userID := int64(1)
	cardID := int64(1)

	boardRepo.EXPECT().GetMemberFromCard(ctx, userID, cardID).Return("admin", int64(1), nil)
	boardRepo.EXPECT().DeleteCard(ctx, cardID).Return(nil)

	err := uc.DeleteCard(ctx, userID, cardID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateColumn_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	boardRepo := mocks.NewMockBoardRepo(ctrl)
	uc := CreateBoardUsecase(boardRepo)
	ctx := context.Background()
	userID := int64(1)
	boardID := int64(1)
	data := &models.ColumnRequest{NewTitle: "New Column"}
	boardRepo.EXPECT().
		GetMemberPermissions(ctx, boardID, userID, false).
		Return(&models.MemberWithPermissions{Role: "admin"}, nil)
	boardRepo.EXPECT().
		CreateColumn(ctx, boardID, data.NewTitle).
		Return(&models.Column{ID: 1, Title: "New Column"}, nil)
	_, err := uc.CreateColumn(ctx, userID, boardID, data)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateColumn_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	boardRepo := mocks.NewMockBoardRepo(ctrl)
	uc := CreateBoardUsecase(boardRepo)

	ctx := context.Background()
	userID := int64(1)
	columnID := int64(1)
	data := &models.ColumnRequest{}

	boardRepo.EXPECT().
		GetMemberFromColumn(ctx, userID, columnID).
		Return("admin", int64(1), nil)

	boardRepo.EXPECT().
		UpdateColumn(ctx, columnID, *data).
		Return(&models.Column{ID: 1}, nil)

	updatedCol, err := uc.UpdateColumn(ctx, userID, columnID, data)
	require.NoError(t, err)
	require.NotNil(t, updatedCol)
}

func TestDeleteColumn_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBoardRepo(ctrl)
	mockRepo.EXPECT().GetMemberFromColumn(context.Background(), int64(1), int64(1)).Return("admin", int64(1), nil)
	mockRepo.EXPECT().DeleteColumn(context.Background(), int64(1)).Return(nil)

	uc := &BoardUsecase{boardRepository: mockRepo}
	err := uc.DeleteColumn(context.Background(), 1, 1)
	assert.NoError(t, err)
}

func TestSetBoardBackground_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockBoardRepo(ctrl)
	ctx := context.Background()
	userID := int64(1)
	boardID := int64(1)
	file := &models.UploadedFile{FileID: &boardID}
	board := &models.Board{ID: boardID}

	mockRepo.EXPECT().GetMemberPermissions(ctx, boardID, userID, false).Return(&models.MemberWithPermissions{Role: "admin"}, nil)
	mockRepo.EXPECT().DeduplicateFile(ctx, file).Return([]string{}, []int64{}, nil)
	mockRepo.EXPECT().SetBoardBackground(ctx, userID, boardID, file).Return(board, nil)
	mockRepo.EXPECT().RegisterFile(ctx, file).Return(nil)

	uc := CreateBoardUsecase(mockRepo)
	_, err := uc.SetBoardBackground(ctx, userID, boardID, file)

	assert.NoError(t, err)
}
