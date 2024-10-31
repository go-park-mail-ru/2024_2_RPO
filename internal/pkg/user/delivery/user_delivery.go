package delivery

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/user/usecase"
	"RPO_back/internal/pkg/utils/requests"
	"RPO_back/internal/pkg/utils/responses"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
func (d *UserDelivery) SetMyAvatar(w http.ResponseWriter, r *http.Request) {
	// Ограничение размера 10 МБ
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Не удалось получить файл", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Создание директории для сохранения файлов, если её нет
	uploadDir := os.Getenv("USER_UPLOADS_DIR")
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, os.ModePerm)
	}

	// Генерация уникального имени файла (опционально)
	filename := filepath.Base(handler.Filename) // Можно добавить префикс или использовать UUID

	// Создание файла на сервере
	filePath := filepath.Join(uploadDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Не удалось создать файл на сервере", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Копирование содержимого загруженного файла в созданный файл на сервере
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Ошибка при сохранении файла", http.StatusInternalServerError)
		return
	}

	// Отправка успешного ответа клиенту
	fmt.Fprintf(w, "Файл успешно загружен: %s\n", filename)
}
