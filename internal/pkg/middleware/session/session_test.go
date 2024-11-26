package session

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	auth "RPO_back/internal/pkg/auth"
	mocks "RPO_back/internal/pkg/auth/delivery/grpc/mocks"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware_NoCookie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthRepo := mocks.NewMockAuthClient(ctrl)

	mw := CreateSessionMiddleware(mockAuthRepo)

	handler := mw.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userIDFromContext, _ := UserIDFromContext(r.Context())
		assert.Equal(t, 0, userIDFromContext)
	}))

	req, _ := http.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
}

func TestMiddleware_InvalidCookie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthRepo := mocks.NewMockAuthClient(ctrl)
	mockAuthRepo.EXPECT().CheckSession(gomock.Any(), "invalid-session-id").Return(0, errors.New("Invalid session id"))

	mw := CreateSessionMiddleware(mockAuthRepo)

	handler := mw.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDFromContext, _ := UserIDFromContext(r.Context())
		assert.Equal(t, 0, userIDFromContext)
	}))

	req, _ := http.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: auth.SessionCookieName, Value: "invalid-session-id"})
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
}

func TestMiddleware_ValidCookie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthRepo := mocks.NewMockAuthClient(ctrl)
	userID := 123
	mockAuthRepo.EXPECT().CheckSession(gomock.Any(), "valid-session-id").Return(userID, nil)

	mw := CreateSessionMiddleware(mockAuthRepo)

	handler := mw.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := UserIDFromContext(r.Context())
		assert.True(t, ok)
		assert.Equal(t, userID, id)
	}))

	req, _ := http.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: auth.SessionCookieName, Value: "valid-session-id"})
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
}
