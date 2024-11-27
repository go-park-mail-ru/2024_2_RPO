package session

import (
	auth "RPO_back/internal/pkg/auth"
	"net/http"
	"net/http/httptest"
	"testing"

	AuthGRPC "RPO_back/internal/pkg/auth/delivery/grpc/gen"
	mocks "RPO_back/internal/pkg/auth/delivery/grpc/mocks"

	gomock "github.com/golang/mock/gomock"
)

func TestSessionMiddleware_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAuth := mocks.NewMockAuthClient(ctrl)
	mockAuth.EXPECT().CheckSession(gomock.Any(), &AuthGRPC.CheckSessionRequest{SessionID: "valid_session"}).
		Return(&AuthGRPC.UserDataResponse{UserID: 123}, nil)

	mw := CreateSessionMiddleware(mockAuth)

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.AddCookie(&http.Cookie{Name: auth.SessionCookieName, Value: "valid_session"})
	w := httptest.NewRecorder()

	mw.Middleware(next).ServeHTTP(w, req)

	if !nextCalled {
		t.Error("next handler was not called")
	}
}
