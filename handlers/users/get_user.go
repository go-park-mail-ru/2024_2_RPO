package users

import (
	"RPO_back/database"
	"RPO_back/logs"
	"RPO_back/models"
	"encoding/json"
	"net/http"
)

func GetMe(w http.ResponseWriter, r *http.Request) {
	// Получаем cookie сессии
	sessionCookie, err := r.Cookie("session")
	if err != nil || sessionCookie.Value == "" {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	// Получаем подключение к базе данных
	db, err := database.GetDbConnection()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		logs.GetLogger().Printf("Ошибка подключения к базе данных: %s", err)
		return
	}
	defer db.Close()

	// SQL-запрос для получения информации о пользователе
// 	var user models.User
// 	query := "SELECT id, username, email FROM users WHERE session_id = $1"
// 	err = db.QueryRow(query, sessionCookie.Value).Scan(&user.ID, &user.Username, &user.Email)
// 	if err != nil {
// 		http.Error(w, "Пользователь не найден", http.StatusNotFound)
// 		logs.GetLogger().Printf("Ошибка получения пользователя: %s", err)
// 		return
// 	}

// 	// Устанавливаем заголовки и отправляем ответ
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(user)
// }
