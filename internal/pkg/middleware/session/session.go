package session

import (
	"RPO_back/internal/pkg/auth/repository"
	"context"
	"net/http"
)

func SessionMiddleware(repo *repository.AuthRepository, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		userID, err := repo.RetrieveUserIdFromSessionId(cookie.Value)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:   "session",
				MaxAge: -1,
			})
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), repository.UserIDContextKey, userID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
