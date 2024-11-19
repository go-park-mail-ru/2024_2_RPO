package delivery

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/auth"
	"RPO_back/internal/pkg/user"
	"RPO_back/internal/pkg/utils/requests"
	"RPO_back/internal/pkg/utils/responses"
	"errors"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type UserDelivery struct {
	userUC user.UserUsecase
}

func CreateUserDelivery(userUC user.UserUsecase) *UserDelivery {
	return &UserDelivery{userUC: userUC}
}

// GetMyProfile возвращает пользователю его профиль
func (d *UserDelivery) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	funcName := "GetMyProfile"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}
	profile, err := d.userUC.GetMyProfile(r.Context(), userID)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, funcName)
		return
	}
	responses.DoJSONResponse(w, profile, 200)
}

// UpdateMyProfile обновляет профиль пользователя и возвращает обновлённый профиль
func (d *UserDelivery) UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	funcName := "UpdateMyProfile"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}
	data := models.UserProfileUpdateRequest{}
	err := requests.GetRequestData(r, &data)
	if err != nil {
		responses.DoBadResponse(w, 400, "bad request")
		return
	}
	newProfile, err := d.userUC.UpdateMyProfile(r.Context(), userID, &data)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, funcName)
		return
	}
	responses.DoJSONResponse(w, newProfile, 200)
}

// SetMyAvatar принимает у пользователя файл изображения, сохраняет его,
// устанавливает как аватарку и возвращает обновлённый профиль
func (d *UserDelivery) SetMyAvatar(w http.ResponseWriter, r *http.Request) {
	funcName := "SetMyAvatar"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	// Ограничение размера 10 МБ
	r.ParseMultipartForm(10 << 20)

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		responses.DoBadResponse(w, 400, "bad request")
		log.Warn(funcName, ": ", err)
		return
	}
	defer file.Close()

	updatedProfile, err := d.userUC.SetMyAvatar(r.Context(), userID, &file, fileHeader)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, funcName)
		return
	}

	responses.DoJSONResponse(w, updatedProfile, 200)
}

// LoginUser обеспечивает вход в сервис
func (d *UserDelivery) LoginUser(w http.ResponseWriter, r *http.Request) {
	// Получить данные из запроса
	var loginRequest models.LoginRequest
	err := requests.GetRequestData(r, &loginRequest)
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "Invalid request")
		log.Warn("LoginUser (getting data): ", err)
		return
	}

	sessionID, err := d.userUC.LoginUser(r.Context(), loginRequest.Email, loginRequest.Password)
	if err != nil {
		if errors.Is(err, errs.ErrWrongCredentials) {
			responses.DoBadResponse(w, 401, "Wrong credentials")
			log.Warn("LoginUser (checking credentials): ", err)
			return
		}
		responses.DoBadResponse(w, 500, "Internal Server Error")
		log.Error("LoginUser (checking credentials): ", err)
		return
	}

	cookie := http.Cookie{
		Name:     auth.SessionCookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int((7 * 24 * time.Hour).Seconds()),
	}
	http.SetCookie(w, &cookie)

	responses.DoEmptyOkResponse(w)
}

// RegisterUser регистрирует пользователя
func (d *UserDelivery) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.UserRegistration
	err := requests.GetRequestData(r, &user)
	if err != nil {
		log.Error("UserDelivery Parsing JSON: ", err)
		responses.DoBadResponse(w, http.StatusBadRequest, "Bad request")
		return
	}

	sessionID, err := d.userUC.RegisterUser(r.Context(), &user)
	if err != nil {
		log.Error("Auth: ", err)
		if errors.Is(err, errs.ErrBusyEmail) && errors.Is(err, errs.ErrBusyNickname) {
			responses.DoBadResponse(w, http.StatusConflict, "Email and nickname are busy")
		} else if errors.Is(err, errs.ErrBusyEmail) {
			responses.DoBadResponse(w, http.StatusConflict, "Email is busy")
		} else if errors.Is(err, errs.ErrBusyNickname) {
			responses.DoBadResponse(w, http.StatusConflict, "Nickname is busy")
		} else {
			responses.DoBadResponse(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	cookie := http.Cookie{
		Name:     auth.SessionCookieName,
		Value:    sessionID,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int((7 * 24 * time.Hour).Seconds()),
	}
	http.SetCookie(w, &cookie)

	responses.DoEmptyOkResponse(w)
}

// LogoutUser разлогинивает пользователя
func (d *UserDelivery) LogoutUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(auth.SessionCookieName)
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "You are not logged in")
		return
	}

	sessionID := cookie.Value

	err = d.userUC.LogoutUser(r.Context(), sessionID)

	http.SetCookie(w, &http.Cookie{
		Name:   auth.SessionCookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	if err != nil {
		responses.ResponseErrorAndLog(w, err, "LogoutUser")
	} else {
		responses.DoEmptyOkResponse(w)
	}
}

// ChangePassword отвечает за смену пароля
func (d *UserDelivery) ChangePassword(w http.ResponseWriter, r *http.Request) {
	funcName := "ChangePassword"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}
	data := models.ChangePasswordRequest{}
	err := requests.GetRequestData(r, &data)
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}
	err = d.userUC.ChangePassword(r.Context(), userID, data.OldPassword, data.NewPassword)
	if err != nil {
		responses.DoBadResponse(w, http.StatusInternalServerError, "internal error")
		return
	}
	responses.DoEmptyOkResponse(w)
}
