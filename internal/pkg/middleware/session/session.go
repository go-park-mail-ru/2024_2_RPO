package session

import (
	"RPO_back/internal/pkg/auth/repository"
	"context"
	"net/http"
)

type contextKey string

const (
	UserIDContextKey contextKey = "userId"
)

type SessionMiddleware struct {
	authRepo *repository.AuthRepository
}

func CreateSessionMiddleware(authRepo *repository.AuthRepository) *SessionMiddleware {
	return &SessionMiddleware{
		authRepo: authRepo,
	}
}

func (mw *SessionMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		userID, err := mw.authRepo.RetrieveUserIdFromSessionId(cookie.Value)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:   "session_id",
				MaxAge: -1,
			})
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// UserIDFromContext получает userID из контекста запроса
func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDContextKey).(string)
	return userID, ok
}
