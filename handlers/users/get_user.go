package users

import (
	"RPO_back/auth"
	"RPO_back/database"
	"encoding/json"
	"net/http"
)

func GetMe(w http.ResponseWriter, r *http.Request) {
	// Получаем cookie сессии
	sessionCookie, err := r.Cookie("session_id")
	if err != nil || sessionCookie.Value == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}
	userId, err := auth.RetrieveUserIdFromSessionId(sessionCookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
	}
	user, err2 := database.GetUserByID(userId)
	if err2 != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Error while serializing JSON", http.StatusInternalServerError)
		return
	}
}
