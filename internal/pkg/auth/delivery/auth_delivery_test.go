package delivery

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/auth"
	mocks "RPO_back/internal/pkg/auth/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLoginUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUsecase := mocks.NewMockAuthUsecase(ctrl)
	authDelivery := AuthDelivery{authUsecase: mockAuthUsecase}

	loginRequest := models.LoginRequest{
		Email:    "user@example.com",
		Password: "password",
	}
	sessionID := "session123"
	requestBody, _ := json.Marshal(loginRequest)

	mockAuthUsecase.EXPECT().LoginUser(gomock.Any(), loginRequest.Email, loginRequest.Password).Return(sessionID, nil)

	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(requestBody))
	w := httptest.NewRecorder()

	authDelivery.LoginUser(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "success", assertResponseContains(res, "success"))
	cookies := res.Cookies()
	assert.NotEmpty(t, cookies)
	assert.Equal(t, auth.SessionCookieName, cookies[0].Name)
	assert.Equal(t, sessionID, cookies[0].Value)
}

func TestLoginUser_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authDelivery := AuthDelivery{}

	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte("{invalid json}")))
	w := httptest.NewRecorder()

	authDelivery.LoginUser(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, "Invalid request", assertResponseContains(res, "Invalid request"))
}

func TestLoginUser_WrongCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUsecase := mocks.NewMockAuthUsecase(ctrl)

	authDelivery := AuthDelivery{authUsecase: mockAuthUsecase}

	loginRequest := models.LoginRequest{
		Email:    "user@example.com",
		Password: "wrongpassword",
	}
	requestBody, _ := json.Marshal(loginRequest)

	mockAuthUsecase.EXPECT().LoginUser(gomock.Any(), loginRequest.Email, loginRequest.Password).Return("", errs.ErrWrongCredentials)

	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(requestBody))
	w := httptest.NewRecorder()

	authDelivery.LoginUser(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	assert.Equal(t, "Wrong credentials", assertResponseContains(res, "Wrong credentials"))
}

func TestLoginUser_InternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUsecase := mocks.NewMockAuthUsecase(ctrl)

	authDelivery := AuthDelivery{authUsecase: mockAuthUsecase}

	loginRequest := models.LoginRequest{
		Email:    "user@example.com",
		Password: "password",
	}
	requestBody, _ := json.Marshal(loginRequest)

	mockAuthUsecase.EXPECT().LoginUser(gomock.Any(), loginRequest.Email, loginRequest.Password).Return("", errors.New("unexpected error"))

	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(requestBody))
	w := httptest.NewRecorder()

	authDelivery.LoginUser(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Equal(t, "Internal Server Error", assertResponseContains(res, "Internal Server Error"))
}

func assertResponseContains(res *http.Response, expected string) string {
	var jsonResponse map[string]interface{}
	json.NewDecoder(res.Body).Decode(&jsonResponse)
	text, _ := jsonResponse["text"].(string)
	return text
}

func TestRegisterUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUsecase := mocks.NewMockAuthUsecase(ctrl)
	authDelivery := AuthDelivery{authUsecase: mockAuthUsecase}

	user := models.UserRegistration{
		Email:    "user@example.com",
		Name:     "nickname",
		Password: "password",
	}
	sessionID := "session123"
	requestBody, _ := json.Marshal(user)

	mockAuthUsecase.EXPECT().RegisterUser(gomock.Any(), &user).Return(sessionID, nil)

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(requestBody))
	w := httptest.NewRecorder()

	authDelivery.RegisterUser(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "success", assertResponseContains(res, "success"))
	cookies := res.Cookies()
	assert.NotEmpty(t, cookies)
	assert.Equal(t, auth.SessionCookieName, cookies[0].Name)
	assert.Equal(t, sessionID, cookies[0].Value)
}

func TestRegisterUser_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authDelivery := AuthDelivery{}

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte("{invalid json}")))
	w := httptest.NewRecorder()

	authDelivery.RegisterUser(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, "Bad request", assertResponseContains(res, "Bad request"))
}

func TestRegisterUser_EmailBusy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUsecase := mocks.NewMockAuthUsecase(ctrl)
	authDelivery := AuthDelivery{authUsecase: mockAuthUsecase}

	user := models.UserRegistration{
		Email:    "user@example.com",
		Name:     "nickname",
		Password: "password",
	}
	requestBody, _ := json.Marshal(user)

	mockAuthUsecase.EXPECT().RegisterUser(gomock.Any(), &user).Return("", errs.ErrBusyEmail)

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(requestBody))
	w := httptest.NewRecorder()

	authDelivery.RegisterUser(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusConflict, res.StatusCode)
	assert.Equal(t, "Email is busy", assertResponseContains(res, "Email is busy"))
}

func TestRegisterUser_NicknameBusy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUsecase := mocks.NewMockAuthUsecase(ctrl)
	authDelivery := AuthDelivery{authUsecase: mockAuthUsecase}

	user := models.UserRegistration{
		Email:    "user@example.com",
		Name:     "nickname",
		Password: "password",
	}
	requestBody, _ := json.Marshal(user)

	mockAuthUsecase.EXPECT().RegisterUser(gomock.Any(), &user).Return("", errs.ErrBusyNickname)

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(requestBody))
	w := httptest.NewRecorder()

	authDelivery.RegisterUser(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusConflict, res.StatusCode)
	assert.Equal(t, "Nickname is busy", assertResponseContains(res, "Nickname is busy"))
}

func TestRegisterUser_InternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUsecase := mocks.NewMockAuthUsecase(ctrl)
	authDelivery := AuthDelivery{authUsecase: mockAuthUsecase}

	user := models.UserRegistration{
		Email:    "user@example.com",
		Name:     "nickname",
		Password: "password",
	}
	requestBody, _ := json.Marshal(user)

	mockAuthUsecase.EXPECT().RegisterUser(gomock.Any(), &user).Return("", errors.New("unexpected error"))

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(requestBody))
	w := httptest.NewRecorder()

	authDelivery.RegisterUser(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Equal(t, "Internal server error", assertResponseContains(res, "Internal server error"))
}
