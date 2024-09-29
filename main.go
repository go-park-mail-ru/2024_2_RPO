package main

import (
	"RPO_back/auth"
	"RPO_back/database"
	auth_handlers "RPO_back/handlers/auth"
	boards_handlers "RPO_back/handlers/boards"
	user_handlers "RPO_back/handlers/users"
	"RPO_back/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Создаём логгер
	logger := log.Default()

	// Обрабатываем файл .env
	serverConfig, err := utils.LoadDotEnv()
	if err != nil {
		log.Fatalf(err.Error())
	}
	logger.Printf("Server config: %#v", serverConfig)

	// Подключаемся к базе
	err2 := database.InitDBConnection(serverConfig.DbPort, serverConfig.DbUser, serverConfig.DbPasswd)
	if err2 != nil {
		log.Fatal(err2.Error())
	}

	// Подключаемся к Redis
	err3 := auth.ConnectToRedis(serverConfig.RedisPort, serverConfig.RedisUser, serverConfig.RedisPasswd)
	if err3 != nil {
		log.Fatal(err3.Error())
	}

	// Создаём новый маршрутизатор
	r := mux.NewRouter()

	r.Use(loggingMiddleware)
	r.Use(corsMiddleware)

	// Регистрируем обработчики
	r.HandleFunc("/auth/register", auth_handlers.RegisterUser).Methods("POST", "OPTIONS")
	r.HandleFunc("/users/me", user_handlers.GetMe).Methods("GET", "OPTIONS")
	r.HandleFunc("/boards/my", boards_handlers.GetMyBoardsHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/boards", boards_handlers.CreateBoardHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/boards/{boardId}", boards_handlers.DeleteBoardHandler).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/auth/login", auth_handlers.LoginUser).Methods("POST", "OPTIONS")
	r.HandleFunc("/auth/logout", auth_handlers.LogoutUser).Methods("POST", "OPTIONS")

	// Определяем адрес и порт для сервера
	addr := fmt.Sprintf(":%d", serverConfig.ServerPort)
	fmt.Printf("Сервер запущен на http://localhost%s\n", addr)

	// Запускаем сервер
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Запрос: %s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

// Middleware для настройки CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("Cors mware")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Для префлайт-запросов
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
