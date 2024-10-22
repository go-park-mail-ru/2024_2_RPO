package users

import (
	"RPO_back/internal/pkg/auth/repository"
	"encoding/json"
	"net/http"
)

func GetMe(w http.ResponseWriter, r *http.Request) {
	userId, err := repository.GetUserId(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	user, err2 := repository.GetUserByID(userId)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Error while serializing JSON", http.StatusInternalServerError)
		return
	}
}
