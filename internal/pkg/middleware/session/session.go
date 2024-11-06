package session

import (
	auth "RPO_back/internal/pkg/auth"
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
		cookie, err := r.Cookie(auth.SessionCookieName)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		userID, err := mw.authRepo.RetrieveUserIdFromSessionId(r.Context(), cookie.Value)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:   auth.SessionCookieName,
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
func UserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserIDContextKey).(int)

	return userID, ok
}
