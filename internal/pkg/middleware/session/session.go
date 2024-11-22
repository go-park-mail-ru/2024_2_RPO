package session

import (
	auth "RPO_back/internal/pkg/auth"
	AuthGRPC "RPO_back/internal/pkg/auth/delivery/grpc/gen"
	"context"
	"net/http"
)

type contextKey string

const (
	UserIDContextKey contextKey = "userId"
)

type SessionMiddleware struct {
	authGRPC AuthGRPC.AuthClient
}

func CreateSessionMiddleware(authGRPC AuthGRPC.AuthClient) *SessionMiddleware {
	return &SessionMiddleware{
		authGRPC: authGRPC,
	}
}

func (mw *SessionMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(auth.SessionCookieName)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		responce, err := mw.authGRPC.CheckSession(r.Context(), &AuthGRPC.CheckSessionRequest{SessionID: cookie.Value})
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:   auth.SessionCookieName,
				MaxAge: -1,
			})
			next.ServeHTTP(w, r)
			return
		}

		userID := responce.GetUserID()

		ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// UserIDFromContext получает userID из контекста запроса
func UserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDContextKey).(int64)

	return userID, ok
}
