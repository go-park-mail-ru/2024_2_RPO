package auth

import (
	"RPO_back/auth"
	"RPO_back/database"
	"RPO_back/logs"
	"RPO_back/models"
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.UserRegistration
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if len(user.Password) < 8 {
		http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	db, err := database.GetDbConnection()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var userID int
	query := `INSERT INTO "User" (nickname, email, password_hash, description, joined_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING u_id`
	err = db.QueryRow(query, user.Name, user.Email, string(hashedPassword), "", time.Now(), time.Now()).Scan(&userID)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		logs.GetLogger().Printf("Error in RegisterUser: %s", err)
		return
	}

	sessionID := auth.GenerateSessionID()

	err = auth.RegisterSessionRedis(sessionID, userID)
	if err != nil {
		http.Error(w, "Failed to register session", http.StatusInternalServerError)
		logs.GetLogger().Printf("Error in RegisterSessionRedis: %s", err)
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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Session cookie is set"))

	// 8. Редирект на /index.html
	http.Redirect(w, r, "/index.html", http.StatusFound)
}
