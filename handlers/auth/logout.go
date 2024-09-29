package auth

import (
	"RPO_back/auth"
	"context"
	"fmt"
	"net/http"
)

var ctx = context.Background()

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "You are not logged in", http.StatusBadRequest)
        return
	}

	sessionID := cookie.Value

	rdb := auth.GetRedisConnection()

	if err := rdb.Del(ctx, sessionID).Err(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to log out: %v", err), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	w.WriteHeader(http.StatusOK)
}
