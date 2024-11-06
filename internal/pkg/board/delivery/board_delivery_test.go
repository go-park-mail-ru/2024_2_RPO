package delivery_test

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/board/delivery"
	BoardDelivery "RPO_back/internal/pkg/board/delivery"
	mocks "RPO_back/internal/pkg/board/mocks"
	"RPO_back/internal/pkg/middleware/session"
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
		reqData := models.CreateBoardRequest{Name: "New Board"}
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
		reqData := models.CreateBoardRequest{Name: "New Board"}

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
		userID := 1
		boardID := 1
		reqData := models.BoardPutRequest{NewName: "Updated Board"}
		expectedBoard := models.Board{ID: boardID, Name: reqData.NewName, Description: "", BackgroundImageURL: "", CreatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), UpdatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)}

		mockBoardUsecase.EXPECT().UpdateBoard(gomock.Any(), userID, boardID, reqData).Return(&expectedBoard, nil)

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		reqBody, err := json.Marshal(reqData)
		assert.NoError(t, err)

		req := httptest.NewRequestWithContext(ctx, "PUT", "/boards/board_1", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/boards/{boardId}", boardDelivery.UpdateBoard).Methods("PUT")

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
		r.HandleFunc("/boards/{boardId}", boardDelivery.UpdateBoard).Methods("PUT")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase returns error", func(t *testing.T) {
		userID := 1
		boardID := 1
		reqData := models.BoardPutRequest{NewName: "Updated Board"}

		mockBoardUsecase.EXPECT().UpdateBoard(gomock.Any(), userID, boardID, reqData).Return(nil, errors.New("usecase error"))

		ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

		reqBody, err := json.Marshal(reqData)
		assert.NoError(t, err)

		req := httptest.NewRequestWithContext(ctx, "PUT", "/boards/board_1", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/boards/{boardId}", boardDelivery.UpdateBoard).Methods("PUT")

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
		r.HandleFunc("/boards/{boardId}", boardDelivery.DeleteBoard).Methods("DELETE")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetUserIDOrFail fails", func(t *testing.T) {
		ctx := context.Background()
		req := httptest.NewRequestWithContext(ctx, "DELETE", "/boards/board_1", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/boards/{boardId}", boardDelivery.DeleteBoard).Methods("DELETE")

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
		r.HandleFunc("/boards/{boardId}", boardDelivery.DeleteBoard).Methods("DELETE")

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
		r.HandleFunc("/boards/{boardId}", boardDelivery.DeleteBoard).Methods("DELETE")

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

	// t.Run("successful retrieval of members' permissions", func(t *testing.T) {
	// 	userID := 1
	// 	boardID := 1

	// 	expectedPermissions := []models.MemberWithPermissions{
	// 		{User: &models.UserProfile{ID: 1}, Role: "admin", AddedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), UpdatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)},
	// 		{User: &models.UserProfile{ID: 2}, Role: "editor", AddedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), UpdatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)},
	// 	}

	// 	mockBoardUsecase.EXPECT().GetMembersPermissions(gomock.Any(), gomock.Eq(userID), gomock.Eq(boardID)).Return(expectedPermissions, nil)

	// 	ctx := context.WithValue(context.Background(), session.UserIDContextKey, userID)

	// 	req := httptest.NewRequest("GET", "/userPermissions/1", nil).WithContext(ctx)
	// 	req = mux.SetURLVars(req, map[string]string{"boardId": "1"})

	// 	w := httptest.NewRecorder()

	// 	r := mux.NewRouter()
	// 	r.HandleFunc("/userPermissions/{boardId}", boardDelivery.GetMembersPermissions).Methods("GET")

	// 	r.ServeHTTP(w, req)

	// 	assert.Equal(t, http.StatusOK, w.Code)

	// 	var gotPermissions []models.MemberWithPermissions
	// 	err := json.NewDecoder(w.Body).Decode(&gotPermissions)
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, expectedPermissions, gotPermissions)
	// })

	t.Run("GetUserIDOrFail fails", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			boardDelivery.GetMembersPermissions(w, r)
		})

		req := httptest.NewRequest("GET", "/userPermissions/1", nil)
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

		req := httptest.NewRequest("GET", "/userPermissions/1", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// t.Run("usecase returns error", func(t *testing.T) {
	// 	userID := 1
	// 	boardID := 1

	// 	mockBoardUsecase.EXPECT().GetMembersPermissions(gomock.Any(), userID, boardID).Return(nil, errors.New("usecase error"))

	// 	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		ctx := context.WithValue(r.Context(), session.UserIDContextKey, userID)
	// 		r = r.WithContext(ctx)
	// 		boardDelivery.GetMembersPermissions(w, r)
	// 	})

	// 	req := httptest.NewRequest("GET", "/userPermissions/1", nil)
	// 	w := httptest.NewRecorder()

	// 	handler.ServeHTTP(w, req)

	// 	assert.Equal(t, http.StatusBadRequest, w.Code)
	// })
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

		req := httptest.NewRequest("POST", "/boards/board_1/members", bytes.NewBuffer(reqBody))
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r := mux.NewRouter()
		r.HandleFunc("/boards/{boardId}/members", boardDelivery.AddMember).Methods("POST")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var gotMember models.MemberWithPermissions
		err = json.NewDecoder(w.Body).Decode(&gotMember)
		assert.NoError(t, err)
		assert.Equal(t, expectedMember, gotMember)
	})

	t.Run("нет авторизации", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			boardDelivery.AddMember(w, r)
		})

		req := httptest.NewRequest("POST", "/boards/board_1/members", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("четырёхсотка", func(t *testing.T) {
		userID := 1

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), session.UserIDContextKey, userID)
			r = r.WithContext(ctx)
			boardDelivery.AddMember(w, r)
		})

		req := httptest.NewRequest("POST", "/boards/4/members", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("демаршалинг", func(t *testing.T) {
		userID := 1

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), session.UserIDContextKey, userID)
			r = r.WithContext(ctx)
			boardDelivery.AddMember(w, r)
		})

		req := httptest.NewRequest("POST", "/boards/board_1/members", bytes.NewBuffer([]byte("Кто прочитал и вернул 500 тот проиграл")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
