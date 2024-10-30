package delivery

import (
	"RPO_back/internal/pkg/user/usecase"
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
func (this *UserDelivery) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	panic("Not implemented")
}

// UpdateMyProfile обновляет профиль пользователя и возвращает обновлённый профиль
func (this *UserDelivery) UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	panic("Not implemented")
}

// SetMyAvatar принимает у пользователя файл изображения, сохраняет его,
// устанавливает как аватарку и возвращает обновлённый профиль
// Самый низкий приоритет
func (this *UserDelivery) SetMyAvatar(w http.ResponseWriter, r *http.Request) {
	responses.DoBadResponse(w, 501, "not implemented")
}
