package delivery

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/user/usecase"
	"RPO_back/internal/pkg/utils/requests"
	"RPO_back/internal/pkg/utils/responses"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type UserDelivery struct {
	userUC *usecase.UserUsecase
}

func CreateUserDelivery(userUC *usecase.UserUsecase) *UserDelivery {
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
	responses.DoJSONResponce(w, profile, 200)
}

// UpdateMyProfile обновляет профиль пользователя и возвращает обновлённый профиль
func (d *UserDelivery) UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	funcName := "UpdateMyProfile"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}
	data := models.UserProfileUpdate{}
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
	responses.DoJSONResponce(w, newProfile, 200)
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

	responses.DoJSONResponce(w, updatedProfile, 200)
}
