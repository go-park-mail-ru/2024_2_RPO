package main

import (
	"RPO_back/auth"
	auth_handlers "RPO_back/handlers/auth"
	boards_handlers "RPO_back/handlers/boards"
	user_handlers "RPO_back/handlers/users"
	"RPO_back/internal/pkg/auth/repository"
	"RPO_back/internal/pkg/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Инициализировать приложение
// Возвращает мультиплексор, который можно тестировать, а можно запрячь для обработки запросов
func initializeApp() (*mux.Router, error) {
	// Создаём логгер
	logger := log.Default()

	// Обрабатываем файл .env
	serverConfig, err := utils.LoadDotEnv()
	if err != nil {
		return nil, fmt.Errorf("error while load .env file: %w", err)
	}
	logger.Printf("Server config: %#v", serverConfig)

	// Подключаемся к базе
	if err := repository.InitDBConnection(serverConfig.DbUrl); err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	// Подключаемся к Redis
	if err := auth.ConnectToRedis(serverConfig.RedisUrl); err != nil {
		return nil, fmt.Errorf("ошибка подключения к Redis: %w", err)
	}

	// Создаём новый маршрутизатор
	r := mux.NewRouter()

	// Применяем middleware
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

	return r, nil
}

func main() {
	router, err := initializeApp()
	if err != nil {
		log.Fatalf("Ошибка инициализации приложения: %v", err)
	}

	// Определяем адрес и порт для сервера
	serverConfig, _ := utils.LoadDotEnv() // Предполагается, что ошибка уже проверена в initializeApp
	addr := fmt.Sprintf(":%d", serverConfig.ServerPort)
	fmt.Printf("Сервер запущен на http://localhost%s\n", addr)

	// Запускаем сервер
	if err := http.ListenAndServe(addr, router); err != nil {
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
