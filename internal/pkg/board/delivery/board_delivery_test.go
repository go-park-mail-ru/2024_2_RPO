package delivery

import (
	mocks "RPO_back/internal/pkg/board/mocks"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestCreateNewBoard_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockBoardUsecase(ctrl)
	boardDelivery := CreateBoardDelivery(mockUsecase)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/boards", bytes.NewReader([]byte(`{"title":"Test Board"}`)))

	boardDelivery.CreateNewBoard(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %v, got %v", http.StatusUnauthorized, w.Code)
	}
}

func TestUpdateBoard_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockBoardUsecase(ctrl)
	boardDelivery := CreateBoardDelivery(mockUsecase)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/boards/1", bytes.NewReader([]byte(`{"title":"Updated Board"}`)))

	boardDelivery.UpdateBoard(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %v, got %v", http.StatusUnauthorized, w.Code)
	}
}

func TestDeleteBoard_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockBoardUsecase(ctrl)
	boardDelivery := CreateBoardDelivery(mockUsecase)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/boards/1", nil)

	boardDelivery.DeleteBoard(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %v, got %v", http.StatusUnauthorized, w.Code)
	}
}

func TestGetMyBoards_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockBoardUsecase(ctrl)
	boardDelivery := CreateBoardDelivery(mockUsecase)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/boards", nil)

	boardDelivery.GetMyBoards(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %v, got %v", http.StatusUnauthorized, w.Code)
	}
}

func TestGetMembersPermissions_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockBoardUsecase(ctrl)
	boardDelivery := CreateBoardDelivery(mockUsecase)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/boards/1/members", nil)

	boardDelivery.GetMembersPermissions(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %v, got %v", http.StatusUnauthorized, w.Code)
	}
}

func TestUpdateMemberRole_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockBoardUsecase(ctrl)
	boardDelivery := CreateBoardDelivery(mockUsecase)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/boards/1/members/2", bytes.NewReader([]byte(`{"newRole":"editor"}`)))

	boardDelivery.UpdateMemberRole(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %v, got %v", http.StatusUnauthorized, w.Code)
	}
}

func TestRemoveMember_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockBoardUsecase(ctrl)
	boardDelivery := CreateBoardDelivery(mockUsecase)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/boards/1/members/2", nil)

	boardDelivery.RemoveMember(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %v, got %v", http.StatusUnauthorized, w.Code)
	}
}

func TestGetBoardContent_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockBoardUsecase(ctrl)
	boardDelivery := CreateBoardDelivery(mockUsecase)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/boards/1/content", nil)

	boardDelivery.GetBoardContent(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %v, got %v", http.StatusUnauthorized, w.Code)
	}
}

func TestCreateNewCard_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockBoardUsecase(ctrl)
	boardDelivery := CreateBoardDelivery(mockUsecase)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/boards/1/cards", bytes.NewReader([]byte(`{"title":"New Card"}`)))

	boardDelivery.CreateNewCard(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %v, got %v", http.StatusUnauthorized, w.Code)
	}
}

func TestUpdateCard_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockBoardUsecase(ctrl)
	boardDelivery := CreateBoardDelivery(mockUsecase)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/cards/card_1", bytes.NewReader([]byte(`{"title":"Updated Card"}`)))

	boardDelivery.UpdateCard(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %v, got %v", http.StatusUnauthorized, w.Code)
	}
}
