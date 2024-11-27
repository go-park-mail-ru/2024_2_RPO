package delivery

import (
	mocks "RPO_back/internal/pkg/user/mocks"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestGetMyProfile_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockUserUsecase(ctrl)
	userDelivery := CreateUserDelivery(mockUsecase)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/profile", nil)

	userDelivery.GetMyProfile(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %v, got %v", http.StatusUnauthorized, w.Code)
	}
}

func TestUpdateMyProfile_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockUserUsecase(ctrl)
	userDelivery := CreateUserDelivery(mockUsecase)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/profile", bytes.NewReader([]byte(`{"nickname":"new_nick"}`)))

	userDelivery.UpdateMyProfile(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %v, got %v", http.StatusUnauthorized, w.Code)
	}
}

func TestSetMyAvatar_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockUserUsecase(ctrl)
	userDelivery := CreateUserDelivery(mockUsecase)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/profile/avatar", nil)

	userDelivery.SetMyAvatar(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %v, got %v", http.StatusUnauthorized, w.Code)
	}
}

func TestLoginUser_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockUserUsecase(ctrl)
	userDelivery := CreateUserDelivery(mockUsecase)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte(`invalid-json`)))

	userDelivery.LoginUser(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %v, got %v", http.StatusBadRequest, w.Code)
	}
}

func TestLogoutUser_NotLoggedIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockUserUsecase(ctrl)
	userDelivery := CreateUserDelivery(mockUsecase)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/logout", nil)

	userDelivery.LogoutUser(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %v, got %v", http.StatusBadRequest, w.Code)
	}
}

func TestChangePassword_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockUserUsecase(ctrl)
	userDelivery := CreateUserDelivery(mockUsecase)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/password/change", bytes.NewReader([]byte(`{"oldPassword":"123","newPassword":"456"}`)))

	userDelivery.ChangePassword(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %v, got %v", http.StatusUnauthorized, w.Code)
	}
}

func TestSubmitPoll_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockUserUsecase(ctrl)
	userDelivery := CreateUserDelivery(mockUsecase)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/poll/submit", bytes.NewReader([]byte(`{"pollID":1,"answers":[2,3]}`)))

	userDelivery.SubmitPoll(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %v, got %v", http.StatusUnauthorized, w.Code)
	}
}

func TestGetPollResults_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockUserUsecase(ctrl)
	userDelivery := CreateUserDelivery(mockUsecase)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/poll/results", nil)

	userDelivery.GetPollResults(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %v, got %v", http.StatusUnauthorized, w.Code)
	}
}
