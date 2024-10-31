package delivery

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/user/usecase"
	"RPO_back/internal/pkg/utils/requests"
	"RPO_back/internal/pkg/utils/responses"
	"net/http"
)

type UserDelivery struct {
	userUC *usecase.UserUsecase
}

func CreateUserDelivery(userUC *usecase.UserUsecase) *UserDelivery {
	return &UserDelivery{userUC: userUC}
}

// GetMyProfile возвращает пользователю его профиль
func (d *UserDelivery) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := requests.GetUserIDOrFail(w, r, "GetMyProfile")
	if !ok {
		return
	}
	profile, err := d.userUC.GetMyProfile(userID)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, "GetMyProfile")
		return
	}
	responses.DoJSONResponce(w, profile, 200)
}

// UpdateMyProfile обновляет профиль пользователя и возвращает обновлённый профиль
func (d *UserDelivery) UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := requests.GetUserIDOrFail(w, r, "UpdateMyProfile")
	if !ok {
		return
	}
	data := models.UserProfileUpdate{}
	err := requests.GetRequestData(r, &data)
	if err != nil {
		responses.DoBadResponse(w, 400, "bad request")
	}
	newProfile, err := d.userUC.UpdateMyProfile(userID, &data)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, "UpdateMyProfile")
		return
	}
	responses.DoJSONResponce(w, newProfile, 200)
}

// SetMyAvatar принимает у пользователя файл изображения, сохраняет его,
// устанавливает как аватарку и возвращает обновлённый профиль
// Самый низкий приоритет
func (d *UserDelivery) SetMyAvatar(w http.ResponseWriter, r *http.Request) {
	responses.DoBadResponse(w, 501, "not implemented")
}
