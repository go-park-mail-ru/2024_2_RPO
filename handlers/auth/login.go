package auth

import (
	"RPO_back/auth"
	"RPO_back/database"
	"RPO_back/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginRequest models.LoginRequest
	var user models.User
	var hashedPassword string

	sessionID := auth.GenerateSessionID()

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	db, err := database.GetDbConnection()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	err2 := db.QueryRow(context.Background(), "SELECT u_id, nickname, email, description, joined_at, updated_at, password_hash FROM \"User\" WHERE email=$1", loginRequest.Email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Description,
		&user.JoinedAt,
		&user.UpdatedAt,
		&hashedPassword,
	)

	if err2 != nil {
		if err2 == sql.ErrNoRows {
			http.Error(w, "Email not found", http.StatusUnauthorized)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginRequest.Password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   10000,
	}
	http.SetCookie(w, &cookie)
	auth.RegisterSessionRedis(sessionID, user.ID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Session cookie is set"))
	http.Redirect(w, r, "/app", http.StatusFound)
}
