package usecase

import (
	"RPO_back/internal/pkg/auth/repository"
	"RPO_back/internal/pkg/utils/encrypt"
	"database/sql"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	authRepo *repository.AuthRepository
}

func CreateAuthUsecase(repo *repository.AuthRepository) *AuthUsecase {
	return &AuthUsecase{
		authRepo: repo,
	}
}

func (this *AuthUsecase) LoginUser(email string, password string) (sessionId string, err error) {
	sessionID := encrypt.GenerateSessionID()

	user, err := this.authRepo.GetUserByEmail(email)

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
	auth.RegisterSessionRedis(sessionID, user.Id)
}
