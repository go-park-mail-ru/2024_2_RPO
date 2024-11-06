package delivery

import (
	"RPO_back/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"RPO_back/internal/pkg/middleware/session"
	mock_user "RPO_back/internal/pkg/user/mocks"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

func TestGetMyProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUserUC := mock_user.NewMockUserUsecase(ctrl)
	handler := CreateUserDelivery(mockUserUC)

	profile := models.UserProfile{ID: 1, Name: "John Doe"}
	mockUserUC.EXPECT().GetMyProfile(gomock.Any(), 1).Return(&profile, nil)

	ctx := context.WithValue(context.Background(), session.UserIDContextKey, 1)
	req, _ := http.NewRequestWithContext(ctx, "GET", "/users/me", nil)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/me", handler.GetMyProfile).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %v, got %v", http.StatusOK, status)
	}

	var responseProfile models.UserProfile
	json.NewDecoder(rr.Body).Decode(&responseProfile)

	if responseProfile != profile {
		t.Errorf("expected response profile %v, got %v", profile, responseProfile)
	}
}

func TestUpdateMyProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUserUC := mock_user.NewMockUserUsecase(ctrl)
	handler := CreateUserDelivery(mockUserUC)

	oldProfile := models.UserProfile{ID: 1, Name: "John Smith"}
	updateData := models.UserProfileUpdate{NewName: "Romanov Vasily", Email: "rvasily@google.com"}

	mockUserUC.EXPECT().UpdateMyProfile(gomock.Any(), 1, &updateData).Return(&oldProfile, nil)

	ctx := context.WithValue(context.Background(), session.UserIDContextKey, 1)
	updateDataJSON, _ := json.Marshal(updateData)
	req, _ := http.NewRequestWithContext(ctx, "PUT", "/users/me", bytes.NewBuffer(updateDataJSON))
	req = mux.SetURLVars(req, map[string]string{"userID": "1"})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/me", handler.UpdateMyProfile).Methods("PUT")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %v, got %v", http.StatusOK, status)
	}

	var responseProfile models.UserProfile
	json.NewDecoder(rr.Body).Decode(&responseProfile)

	if responseProfile != oldProfile {
		t.Errorf("expected response profile %v, got %v", oldProfile, responseProfile)
	}
}

func TestSetMyAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUserUC := mock_user.NewMockUserUsecase(ctrl)
	handler := CreateUserDelivery(mockUserUC)

	updatedProfile := models.UserProfile{ID: 1, Name: "John Doe", AvatarImageURL: "http://example.com/avatar.jpg"}

	mockUserUC.EXPECT().SetMyAvatar(gomock.Any(), 1, gomock.Any(), gomock.Any()).Return(&updatedProfile, nil)

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	part, err := writer.CreateFormFile("file", "avatar.jpg")
	if err != nil {
		t.Fatal(err)
	}
	fileContent := strings.NewReader("fake image content")
	if _, err := io.Copy(part, fileContent); err != nil {
		t.Fatal(err)
	}
	writer.Close()

	ctx := context.WithValue(context.Background(), session.UserIDContextKey, 1)
	req, _ := http.NewRequestWithContext(ctx, "PUT", "/users/me/avatarImage", &buffer)
	req = mux.SetURLVars(req, map[string]string{"userID": "1"})
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/me/avatarImage", handler.SetMyAvatar).Methods("PUT")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %v, got %v", http.StatusOK, status)
	}

	var responseProfile models.UserProfile
	json.NewDecoder(rr.Body).Decode(&responseProfile)

	if responseProfile != updatedProfile {
		t.Errorf("expected response profile %v, got %v", updatedProfile, responseProfile)
	}
}
