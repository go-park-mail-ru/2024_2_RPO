package delivery_test

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/board/delivery"
	BoardDelivery "RPO_back/internal/pkg/board/delivery"
	mocks "RPO_back/internal/pkg/board/mocks"
	"RPO_back/internal/pkg/middleware/session"
	"RPO_back/internal/pkg/utils/misc"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewBoard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBoardUsecase := mocks.NewMockBoardUsecase(ctrl)
	boardDelivery := BoardDelivery.CreateBoardDelivery(mockBoardUsecase)

	t.Run("successful board creation", func(t *testing.T) {
		userID := 1
		reqData := models.BoardRequest{NewName: "New Board"}
		expectedBoard := models.Board{ID: 1, Name: "New Board"}

		mockBoardUsecase.EXPECT().CreateNewBoard(gomock.Any(), userID, reqData).Return(&expectedBoard, nil)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), session.UserIDContextKey, userID)
			r = r.WithContext(ctx)
			boardDelivery.CreateNewBoard(w, r)
		})

		reqBody, err := json.Marshal(reqData)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/boards", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var gotBoard models.Board
		err = json.NewDecoder(w.Body).Decode(&gotBoard)
		assert.NoError(t, err)
		assert.Equal(t, expectedBoard, gotBoard)
	})

	t.Run("invalid request data", func(t *testing.T) {
		userID := 1

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), session.UserIDContextKey, userID)
			r = r.WithContext(ctx)
			boardDelivery.CreateNewBoard(w, r)
		})

		req := httptest.NewRequest("POST", "/boards", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase returns error", func(t *testing.T) {
		userID := 1
		reqData := models.BoardRequest{NewName: "New Board"}

		mockBoardUsecase.EXPECT().CreateNewBoard(gomock.Any(), userID, reqData).Return(nil, errors.New("usecase error"))

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), session.UserIDContextKey, userID)
			r = r.WithContext(ctx)
			boardDelivery.CreateNewBoard(w, r)
		})

		reqBody, err := json.Marshal(reqData)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/boards", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("GetUserIDOrFail fails", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			boardDelivery.CreateNewBoard(w, r)
		})

		req := httptest.NewRequest("POST", "/boards", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestUpdateBoard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBoardUsecase := mocks.NewMockBoardUsecase(ctrl)
	boardDelivery := BoardDelivery.CreateBoardDelivery(mockBoardUsecase)

	t.Run("successful board update", func(t *testing.T) {
		userID := int64(1)
		boardID := int64(1)
		reqData := models.BoardRequest{NewName: "Updated Board"}
		expectedBoard := models.Board{ID: boardID, Name: reqData.NewName, BackgroundImageURL: "", CreatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), UpdatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)}

		mockBoardUsecase.EXPECT().UpdateBoard(gomock.Any(), userID, boardID, reqData).Return(&expectedBoard, nil)

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		reqBody, err := json.Marshal(reqData)
		assert.NoError(t, err)

		req := httptest.NewRequestWithContext(ctx, "PUT", "/boards/board_1", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/boards/{boardID}", boardDelivery.UpdateBoard).Methods("PUT")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var gotBoard models.Board
		err = json.NewDecoder(w.Body).Decode(&gotBoard)
		assert.NoError(t, err)
		assert.Equal(t, expectedBoard, gotBoard)
	})

	t.Run("invalid request data", func(t *testing.T) {
		userID := 1

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		req := httptest.NewRequestWithContext(ctx, "PUT", "/boards/1", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/boards/{boardID}", boardDelivery.UpdateBoard).Methods("PUT")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase returns error", func(t *testing.T) {
		userID := 1
		boardID := 1
		reqData := models.BoardRequest{NewName: "Updated Board"}

		mockBoardUsecase.EXPECT().UpdateBoard(gomock.Any(), userID, boardID, reqData).Return(nil, errors.New("usecase error"))

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		reqBody, err := json.Marshal(reqData)
		assert.NoError(t, err)

		req := httptest.NewRequestWithContext(ctx, "PUT", "/boards/board_1", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/boards/{boardID}", boardDelivery.UpdateBoard).Methods("PUT")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("GetUserIDOrFail fails", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			boardDelivery.UpdateBoard(w, r)
		})

		req := httptest.NewRequest("PUT", "/boards/board_1", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("GetIDFromRequest fails", func(t *testing.T) {
		userID := 1

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), session.UserIDContextKey, userID)
			r = r.WithContext(ctx)
			boardDelivery.UpdateBoard(w, r)
		})

		req := httptest.NewRequest("PUT", "/boards/invalid", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDeleteBoard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBoardUsecase := mocks.NewMockBoardUsecase(ctrl)
	boardDelivery := delivery.CreateBoardDelivery(mockBoardUsecase)

	t.Run("successful board deletion", func(t *testing.T) {
		userID := 1
		boardID := 1

		mockBoardUsecase.EXPECT().DeleteBoard(gomock.Any(), userID, boardID).Return(nil)

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		req := httptest.NewRequestWithContext(ctx, "DELETE", "/boards/board_1", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/boards/{boardID}", boardDelivery.DeleteBoard).Methods("DELETE")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetUserIDOrFail fails", func(t *testing.T) {
		ctx := context.Background()
		req := httptest.NewRequestWithContext(ctx, "DELETE", "/boards/board_1", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/boards/{boardID}", boardDelivery.DeleteBoard).Methods("DELETE")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("GetIDFromRequest fails", func(t *testing.T) {
		userID := 1
		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		req := httptest.NewRequestWithContext(ctx, "DELETE", "/boards/invalid", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/boards/{boardID}", boardDelivery.DeleteBoard).Methods("DELETE")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase returns error", func(t *testing.T) {
		userID := 1
		boardID := 1

		mockBoardUsecase.EXPECT().DeleteBoard(gomock.Any(), userID, boardID).Return(errors.New("usecase error"))

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		req := httptest.NewRequestWithContext(ctx, "DELETE", "/boards/board_1", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/boards/{boardID}", boardDelivery.DeleteBoard).Methods("DELETE")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetMyBoards(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBoardUsecase := mocks.NewMockBoardUsecase(ctrl)
	boardDelivery := BoardDelivery.CreateBoardDelivery(mockBoardUsecase)

	t.Run("successful retrieval of boards", func(t *testing.T) {
		userID := 1
		expectedBoards := []models.Board{
			{ID: 1, Name: "Board 1"},
			{ID: 2, Name: "Board 2"},
		}

		mockBoardUsecase.EXPECT().GetMyBoards(gomock.Any(), userID).Return(expectedBoards, nil)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), session.UserIDContextKey, userID)
			r = r.WithContext(ctx)
			boardDelivery.GetMyBoards(w, r)
		})

		req := httptest.NewRequest("GET", "/boards/my", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var gotBoards []models.Board
		err := json.NewDecoder(w.Body).Decode(&gotBoards)
		assert.NoError(t, err)
		assert.Equal(t, expectedBoards, gotBoards)
	})

	t.Run("GetUserIDOrFail fails", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			boardDelivery.GetMyBoards(w, r)
		})

		req := httptest.NewRequest("GET", "/boards", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("usecase returns error", func(t *testing.T) {
		userID := 1

		mockBoardUsecase.EXPECT().GetMyBoards(gomock.Any(), userID).Return(nil, errors.New("usecase error"))

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), session.UserIDContextKey, userID)
			r = r.WithContext(ctx)
			boardDelivery.GetMyBoards(w, r)
		})

		req := httptest.NewRequest("GET", "/boards", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetMembersPermissions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBoardUsecase := mocks.NewMockBoardUsecase(ctrl)
	boardDelivery := BoardDelivery.CreateBoardDelivery(mockBoardUsecase)

	t.Run("successful retrieval of members' permissions", func(t *testing.T) {
		userID := 1
		boardID := 1

		expectedPermissions := []models.MemberWithPermissions{
			{User: &models.UserProfile{ID: 1}, Role: "admin", AddedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), UpdatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)},
			{User: &models.UserProfile{ID: 2}, Role: "editor", AddedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), UpdatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)},
		}

		mockBoardUsecase.EXPECT().GetMembersPermissions(gomock.Any(), gomock.Eq(userID), gomock.Eq(boardID)).Return(expectedPermissions, nil)

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		req := httptest.NewRequest("GET", "/userPermissions/board_1", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"boardID": "1"})

		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/userPermissions/{boardID}", boardDelivery.GetMembersPermissions).Methods("GET")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var gotPermissions []models.MemberWithPermissions
		err := json.NewDecoder(w.Body).Decode(&gotPermissions)
		assert.NoError(t, err)
		assert.Equal(t, expectedPermissions, gotPermissions)
	})

	t.Run("GetUserIDOrFail fails", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			boardDelivery.GetMembersPermissions(w, r)
		})

		req := httptest.NewRequest("GET", "/userPermissions/board_1", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("GetIDFromRequest fails", func(t *testing.T) {
		userID := 1

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), session.UserIDContextKey, userID)
			r = r.WithContext(ctx)
			boardDelivery.GetMembersPermissions(w, r)
		})

		req := httptest.NewRequest("GET", "/userPermissions/board_1", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase returns error", func(t *testing.T) {
		userID := 1
		boardID := 1

		mockBoardUsecase.EXPECT().GetMembersPermissions(gomock.Any(), userID, boardID).Return(nil, errors.New("usecase error"))

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		req := httptest.NewRequest("GET", "/userPermissions/board_1", nil).WithContext(ctx)
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/userPermissions/{boardID}", boardDelivery.GetMembersPermissions).Methods("GET")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestAddMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBoardUsecase := mocks.NewMockBoardUsecase(ctrl)
	boardDelivery := BoardDelivery.CreateBoardDelivery(mockBoardUsecase)

	t.Run("successful addition of a member", func(t *testing.T) {
		userID := 1
		boardID := 1
		reqData := models.AddMemberRequest{MemberNickname: "user123"}
		expectedMember := models.MemberWithPermissions{User: &models.UserProfile{ID: 2}, Role: "viewer"}

		mockBoardUsecase.EXPECT().AddMember(gomock.Any(), userID, boardID, &reqData).Return(&expectedMember, nil)

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		reqBody, err := json.Marshal(reqData)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/userPermissions/board_1", bytes.NewBuffer(reqBody))
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/userPermissions/{boardID}", boardDelivery.AddMember).Methods("POST")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var gotMember models.MemberWithPermissions
		err = json.NewDecoder(w.Body).Decode(&gotMember)
		assert.NoError(t, err)
		assert.Equal(t, expectedMember, gotMember)
	})

	t.Run("GetUserIDOrFail fails", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			boardDelivery.AddMember(w, r)
		})

		req := httptest.NewRequest("POST", "/userPermissions/board_1", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("GetIDFromRequest fails", func(t *testing.T) {
		userID := 1

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), session.UserIDContextKey, userID)
			r = r.WithContext(ctx)
			boardDelivery.AddMember(w, r)
		})

		req := httptest.NewRequest("POST", "/userPermissions/invalid", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid request data", func(t *testing.T) {
		userID := 1

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), session.UserIDContextKey, userID)
			r = r.WithContext(ctx)
			boardDelivery.AddMember(w, r)
		})

		req := httptest.NewRequest("POST", "/userPermissions/board_1", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase returns error", func(t *testing.T) {
		userID := 1
		boardID := 1
		reqData := models.AddMemberRequest{MemberNickname: "user123"}

		mockBoardUsecase.EXPECT().AddMember(gomock.Any(), userID, boardID, &reqData).Return(nil, errors.New("usecase error"))

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		marshal, _ := json.Marshal(reqData)

		req := httptest.NewRequest("POST", fmt.Sprintf("/userPermissions/board_%d", boardID), bytes.NewBuffer(marshal))
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/userPermissions/{boardID}", boardDelivery.AddMember).Methods("POST")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetBoardContent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBoardUsecase := mocks.NewMockBoardUsecase(ctrl)
	boardDelivery := BoardDelivery.CreateBoardDelivery(mockBoardUsecase)

	t.Run("successful retrieval of board content", func(t *testing.T) {
		userID := 1
		boardID := 1

		mockBoardUsecase.EXPECT().GetBoardContent(gomock.Any(), gomock.Eq(userID), gomock.Eq(boardID)).Return(nil, nil)

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		req := httptest.NewRequest("GET", "/cards/board_1/allContent", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"boardID": "1"})
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/cards/{boardID}/allContent", boardDelivery.GetBoardContent).Methods("GET")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var gotContent models.BoardContent
		err := json.NewDecoder(w.Body).Decode(&gotContent)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code, gotContent)
	})

	t.Run("GetUserIDOrFail fails", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/cards/board_1/allContent", nil)
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/cards/{boardID}/allContent", boardDelivery.GetBoardContent).Methods("GET")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("GetIDFromRequest fails", func(t *testing.T) {
		userID := 1

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		req := httptest.NewRequest("GET", "/cards/invalid_board_id/allContent", nil).WithContext(ctx)
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/cards/{boardID}/allContent", boardDelivery.GetBoardContent).Methods("GET")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase returns error", func(t *testing.T) {
		userID := 1
		boardID := 1

		mockBoardUsecase.EXPECT().GetBoardContent(gomock.Any(), userID, boardID).Return(&models.BoardContent{}, errors.New("usecase error"))

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		req := httptest.NewRequest("GET", "/cards/board_1/allContent", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"boardID": "1"})
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/cards/{boardID}/allContent", boardDelivery.GetBoardContent).Methods("GET")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestCreateNewCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBoardUsecase := mocks.NewMockBoardUsecase(ctrl)
	boardDelivery := BoardDelivery.CreateBoardDelivery(mockBoardUsecase)

	t.Run("successful creation of new card", func(t *testing.T) {
		userID := 1
		boardID := 1
		requestData := models.CardPatchRequest{NewTitle: misc.StringPtr("Title")}

		mockBoardUsecase.EXPECT().CreateNewCard(gomock.Any(), gomock.Eq(userID), gomock.Eq(boardID), gomock.Eq(&requestData)).Return(&models.Card{ID: 1, Title: "", ColumnID: 123}, nil)

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		body, _ := json.Marshal(requestData)
		req := httptest.NewRequest("POST", "/cards/board_1", bytes.NewReader(body)).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"boardID": "1"})
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/cards/{boardID}", boardDelivery.CreateNewCard).Methods("POST")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var gotCard models.Card
		err := json.NewDecoder(w.Body).Decode(&gotCard)
		assert.NoError(t, err)
	})

	t.Run("GetUserIDOrFail fails", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/cards/board_1", nil)
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/cards/{boardID}", boardDelivery.CreateNewCard).Methods("POST")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("GetIDFromRequest fails", func(t *testing.T) {
		userID := 1

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		req := httptest.NewRequest("POST", "/cards/invalid", nil).WithContext(ctx)
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/cards/{boardID}", boardDelivery.CreateNewCard).Methods("POST")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid request data", func(t *testing.T) {
		userID := 1

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		req := httptest.NewRequest("POST", "/cards/board_1", bytes.NewReader([]byte("invalid json"))).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"boardID": "1"})
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/cards/{boardID}", boardDelivery.CreateNewCard).Methods("POST")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase returns error", func(t *testing.T) {
		userID := 1
		boardID := 1
		requestData := models.CardPatchRequest{NewTitle: misc.StringPtr("New Task")}

		mockBoardUsecase.EXPECT().CreateNewCard(gomock.Any(), userID, boardID, &requestData).Return(nil, errors.New("usecase error"))

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		body, _ := json.Marshal(requestData)
		req := httptest.NewRequest("POST", "/cards/board_1", bytes.NewReader(body)).WithContext(ctx)
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/cards/{boardID}", boardDelivery.CreateNewCard).Methods("POST")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUpdateCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBoardUsecase := mocks.NewMockBoardUsecase(ctrl)
	boardDelivery := BoardDelivery.CreateBoardDelivery(mockBoardUsecase)
	t.Run("successful update of card", func(t *testing.T) {
		userID := int64(1)
		cardID := int64(1)
		requestData := models.CardPatchRequest{NewTitle: misc.StringPtr("Updated Task")}

		mockBoardUsecase.EXPECT().UpdateCard(gomock.Any(), userID, cardID, &requestData).Return(&models.Card{ID: cardID, Title: *requestData.NewTitle}, nil)

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		body, _ := json.Marshal(requestData)
		req := httptest.NewRequest("PUT", "/cards/board_1/card_1", bytes.NewReader(body)).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"boardID": "1", "cardId": "1"})
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/cards/{boardID}/{cardId}", boardDelivery.UpdateCard).Methods("PUT")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var gotCard models.Card
		err := json.NewDecoder(w.Body).Decode(&gotCard)
		assert.NoError(t, err)
	})

	t.Run("GetUserIDOrFail fails", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/cards/board_1/card_1", nil)
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/cards/{boardID}/{cardId}", boardDelivery.UpdateCard).Methods("PUT")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("GetIDFromRequest fails for boardID", func(t *testing.T) {
		userID := 1

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		req := httptest.NewRequest("PUT", "/cards/invalid/card_1", nil).WithContext(ctx)
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/cards/{boardID}/{cardId}", boardDelivery.UpdateCard).Methods("PUT")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GetIDFromRequest fails for cardId", func(t *testing.T) {
		userID := 1

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		req := httptest.NewRequest("PUT", "/cards/board_1/invalid", nil).WithContext(ctx)
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/cards/{boardID}/{cardId}", boardDelivery.UpdateCard).Methods("PUT")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid request data", func(t *testing.T) {
		userID := 1

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		req := httptest.NewRequest("PUT", "/cards/board_1/card_1", bytes.NewReader([]byte("invalid json"))).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"boardID": "1", "cardId": "1"})
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/cards/{boardID}/{cardId}", boardDelivery.UpdateCard).Methods("PUT")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase returns error", func(t *testing.T) {
		userID := 1
		cardID := 1
		requestData := models.CardPatchRequest{NewTitle: misc.StringPtr("Updated Task")}

		mockBoardUsecase.EXPECT().UpdateCard(gomock.Any(), userID, cardID, &requestData).Return(nil, errors.New("usecase error"))

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		body, _ := json.Marshal(requestData)
		req := httptest.NewRequest("PUT", "/cards/board_1/card_1", bytes.NewReader(body)).WithContext(ctx)
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/cards/{boardID}/{cardId}", boardDelivery.UpdateCard).Methods("PUT")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestDeleteCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBoardUsecase := mocks.NewMockBoardUsecase(ctrl)
	boardDelivery := BoardDelivery.CreateBoardDelivery(mockBoardUsecase)

	t.Run("successful deletion of card", func(t *testing.T) {
		userID := int64(1)
		cardID := int64(1)

		mockBoardUsecase.EXPECT().DeleteCard(gomock.Any(), userID, cardID).Return(nil)

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		req := httptest.NewRequest("DELETE", "/cards/board_1/card_1", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"boardID": "1", "cardId": "1"})
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/cards/{boardID}/{cardId}", boardDelivery.DeleteCard).Methods("DELETE")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetUserIDOrFail fails", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/cards/board_1/card_1", nil)
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/cards/{boardID}/{cardId}", boardDelivery.DeleteCard).Methods("DELETE")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("GetIDFromRequest fails for boardID", func(t *testing.T) {
		userID := 1

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		req := httptest.NewRequest("DELETE", "/cards/invalid/card_1", nil).WithContext(ctx)
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/cards/{boardID}/{cardId}", boardDelivery.DeleteCard).Methods("DELETE")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GetIDFromRequest fails for cardId", func(t *testing.T) {
		userID := 1

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		req := httptest.NewRequest("DELETE", "/cards/board_1/invalid", nil).WithContext(ctx)
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/cards/{boardID}/{cardId}", boardDelivery.DeleteCard).Methods("DELETE")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase returns error", func(t *testing.T) {
		userID := 1
		cardID := 1

		mockBoardUsecase.EXPECT().DeleteCard(gomock.Any(), userID, cardID).Return(errors.New("usecase error"))

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		req := httptest.NewRequest("DELETE", "/cards/board_1/card_1", nil).WithContext(ctx)
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/cards/{boardID}/{cardId}", boardDelivery.DeleteCard).Methods("DELETE")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestCreateColumn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBoardUsecase := mocks.NewMockBoardUsecase(ctrl)
	boardDelivery := BoardDelivery.CreateBoardDelivery(mockBoardUsecase)

	t.Run("successful creation of column", func(t *testing.T) {
		userID := 1
		boardID := 1
		requestData := models.ColumnRequest{NewTitle: "New Column"}

		mockBoardUsecase.EXPECT().CreateColumn(gomock.Any(), userID, boardID, &requestData).Return(&models.Column{ID: 1, Title: requestData.NewTitle}, nil)

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		body, _ := json.Marshal(requestData)
		req := httptest.NewRequest("POST", "/columns/board_1", bytes.NewReader(body)).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"boardID": "1"})
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/columns/{boardID}", boardDelivery.CreateColumn).Methods("POST")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var gotColumn models.Column
		err := json.NewDecoder(w.Body).Decode(&gotColumn)
		assert.NoError(t, err)
		expectedColumn := models.Column{ID: 1, Title: requestData.NewTitle}
		assert.Equal(t, expectedColumn, gotColumn)
	})

	t.Run("GetUserIDOrFail fails", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/columns/board_1", nil)
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/columns/{boardID}", boardDelivery.CreateColumn).Methods("POST")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("GetIDFromRequest fails", func(t *testing.T) {
		userID := 1

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		req := httptest.NewRequest("POST", "/columns/invalid", nil).WithContext(ctx)
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/columns/{boardID}", boardDelivery.CreateColumn).Methods("POST")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid request data", func(t *testing.T) {
		userID := 1

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		req := httptest.NewRequest("POST", "/columns/board_1", bytes.NewReader([]byte("invalid json"))).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"boardID": "1"})
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/columns/{boardID}", boardDelivery.CreateColumn).Methods("POST")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase returns error", func(t *testing.T) {
		userID := 1
		boardID := 1
		requestData := models.ColumnRequest{NewTitle: "New Column"}

		mockBoardUsecase.EXPECT().CreateColumn(gomock.Any(), userID, boardID, &requestData).Return(nil, errors.New("usecase error"))

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		body, _ := json.Marshal(requestData)
		req := httptest.NewRequest("POST", "/columns/board_1", bytes.NewReader(body)).WithContext(ctx)
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/columns/{boardID}", boardDelivery.CreateColumn).Methods("POST")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUpdateColumn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBoardUsecase := mocks.NewMockBoardUsecase(ctrl)
	boardDelivery := BoardDelivery.CreateBoardDelivery(mockBoardUsecase)

	t.Run("successful update of column", func(t *testing.T) {
		userID := 1
		columnID := 1
		requestData := models.ColumnRequest{NewTitle: "Updated Column"}

		mockBoardUsecase.EXPECT().UpdateColumn(gomock.Any(), userID, columnID, &requestData).Return(&models.Column{ID: columnID, Title: requestData.NewTitle}, nil)

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		body, _ := json.Marshal(requestData)
		req := httptest.NewRequest("PUT", "/columns/board_1/column_1", bytes.NewReader(body)).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"boardID": "1", "columnID": "1"})
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/columns/{boardID}/{columnID}", boardDelivery.UpdateColumn).Methods("PUT")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var gotColumn models.Column
		err := json.NewDecoder(w.Body).Decode(&gotColumn)
		assert.NoError(t, err)
		expectedColumn := models.Column{ID: columnID, Title: requestData.NewTitle}
		assert.Equal(t, expectedColumn, gotColumn)
	})

	t.Run("GetUserIDOrFail fails", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/columns/board_1/column_1", nil)
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/columns/{boardID}/{columnID}", boardDelivery.UpdateColumn).Methods("PUT")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("GetIDFromRequest fails for boardID", func(t *testing.T) {
		userID := 1

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		req := httptest.NewRequest("PUT", "/columns/invalid/column_1", nil).WithContext(ctx)
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/columns/{boardID}/{columnID}", boardDelivery.UpdateColumn).Methods("PUT")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GetIDFromRequest fails for columnID", func(t *testing.T) {
		userID := 1

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		req := httptest.NewRequest("PUT", "/columns/board_1/invalid", nil).WithContext(ctx)
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/columns/{boardID}/{columnID}", boardDelivery.UpdateColumn).Methods("PUT")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid request data", func(t *testing.T) {
		userID := 1

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		req := httptest.NewRequest("PUT", "/columns/board_1/column_1", bytes.NewReader([]byte("invalid json"))).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"boardID": "1", "columnID": "1"})
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/columns/{boardID}/{columnID}", boardDelivery.UpdateColumn).Methods("PUT")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase returns error", func(t *testing.T) {
		userID := 1
		columnID := 1
		requestData := models.ColumnRequest{NewTitle: "Updated Column"}

		mockBoardUsecase.EXPECT().UpdateColumn(gomock.Any(), userID, columnID, &requestData).Return(nil, errors.New("usecase error"))

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		body, _ := json.Marshal(requestData)
		req := httptest.NewRequest("PUT", "/columns/board_1/column_1", bytes.NewReader(body)).WithContext(ctx)
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/columns/{boardID}/{columnID}", boardDelivery.UpdateColumn).Methods("PUT")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestDeleteColumn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBoardUsecase := mocks.NewMockBoardUsecase(ctrl)
	boardDelivery := BoardDelivery.CreateBoardDelivery(mockBoardUsecase)

	t.Run("successful deletion of column", func(t *testing.T) {
		userID := 1
		columnID := 1

		mockBoardUsecase.EXPECT().DeleteColumn(gomock.Any(), userID, columnID).Return(nil)

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		req := httptest.NewRequest("DELETE", "/columns/board_1/column_1", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"boardID": "1", "columnID": "1"})
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/columns/{boardID}/{columnID}", boardDelivery.DeleteColumn).Methods("DELETE")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetUserIDOrFail fails", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/columns/board_1/column_1", nil)
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/columns/{boardID}/{columnID}", boardDelivery.DeleteColumn).Methods("DELETE")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("GetIDFromRequest fails for boardID", func(t *testing.T) {
		userID := 1

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		req := httptest.NewRequest("DELETE", "/columns/invalid/column_1", nil).WithContext(ctx)
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/columns/{boardID}/{columnID}", boardDelivery.DeleteColumn).Methods("DELETE")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GetIDFromRequest fails for columnID", func(t *testing.T) {
		userID := 1

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		req := httptest.NewRequest("DELETE", "/columns/board_1/invalid", nil).WithContext(ctx)
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/columns/{boardID}/{columnID}", boardDelivery.DeleteColumn).Methods("DELETE")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase returns error", func(t *testing.T) {
		userID := 1
		columnID := 1

		mockBoardUsecase.EXPECT().DeleteColumn(gomock.Any(), userID, columnID).Return(errors.New("usecase error"))

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		req := httptest.NewRequest("DELETE", "/columns/board_1/column_1", nil).WithContext(ctx)
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/columns/{boardID}/{columnID}", boardDelivery.DeleteColumn).Methods("DELETE")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
