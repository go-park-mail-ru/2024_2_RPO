package delivery

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/auth"
	"RPO_back/internal/pkg/auth/usecase"
	"RPO_back/internal/pkg/utils/requests"
	"RPO_back/internal/pkg/utils/responses"
	"RPO_back/internal/pkg/utils/validate"
	"errors"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

const sessionIdCookieName string = "session_id"

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
		responses.DoBadResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	sessionId, err := this.authUsecase.LoginUser(loginRequest.Email, loginRequest.Password)
	if err != nil {
		if errors.Is(err, auth.ErrWrongCredentials) {
			responses.DoBadResponse(w, 401, "Wrong credentials")
			return
		}
		responses.DoBadResponse(w, 500, "Internal Server Error")
	}

	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionId,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int((7 * 24 * time.Hour).Seconds()),
	}
	http.SetCookie(w, &cookie)

	responses.DoEmptyOkResponce(w)
}

func (this *AuthDelivery) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.UserRegistration
	err := requests.GetRequestData(r, &user)
	if err != nil {
		log.Error("AuthDelivery Parsing JSON: ", err)
		responses.DoBadResponse(w, http.StatusBadRequest, "Bad request")
		return
	}

	err = validate.Validate(user)
	if err != nil {
		log.Error("AuthDelivery Validating: ", err)
		responses.DoBadResponse(w, http.StatusBadRequest, "Validation error")
		return
	}

	sessionId, err := this.authUsecase.RegisterUser(&user)
	if err != nil {
		log.Error("Auth: ", err)
		if errors.Is(err, auth.ErrBusyEmail) {
			responses.DoBadResponse(w, http.StatusConflict, "Email is busy")
		} else if errors.Is(err, auth.ErrBusyNickname) {
			responses.DoBadResponse(w, http.StatusConflict, "Nickname is busy")
		} else {
			responses.DoBadResponse(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	cookie := http.Cookie{
		Name:     sessionIdCookieName,
		Value:    sessionId,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   10000,
	}
	http.SetCookie(w, &cookie)

	responses.DoEmptyOkResponce(w)
}

func (this *AuthDelivery) LogoutUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionIdCookieName)
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "You are not logged in")
		return
	}

	sessionID := cookie.Value

	err = this.authUsecase.LogoutUser(sessionID)

	http.SetCookie(w, &http.Cookie{
		Name:   sessionIdCookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	if err != nil {
		responses.DoBadResponse(w, 500, "Internal server error")
	} else {
		responses.DoEmptyOkResponce(w)
	}
}
