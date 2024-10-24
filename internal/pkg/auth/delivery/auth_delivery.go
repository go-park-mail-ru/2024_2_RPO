package delivery

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/auth/usecase"
	"RPO_back/internal/pkg/utils/requests"
	"RPO_back/internal/pkg/utils/responses"
	"net/http"
)

type AuthDelivery struct {
	authUsecase *usecase.AuthUsecase
}

func CreateAuthDelivery(uc *usecase.AuthUsecase) *AuthDelivery {
	return &AuthDelivery{
		authUsecase: uc,
	}
}

func (this *AuthDelivery) LoginUser(w http.ResponseWriter, r *http.Request) {
	// Получить данные из запроса
	var loginRequest models.LoginRequest
	err := requests.GetRequestData(r, &loginRequest)
	if err != nil {
		responses.DoBadResponse(w, 400, "Invalid request")
		return
	}

	// Получить ID сессии
	sessionId, err := this.authUsecase.LoginUser(loginRequest.Email, loginRequest.Password)

	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionId,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   10000,
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Session cookie is set"))
	http.Redirect(w, r, "/app", http.StatusFound)
}
