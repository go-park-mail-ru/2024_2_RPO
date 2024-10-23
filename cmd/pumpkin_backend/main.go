package main

import (
	auth_handlers "RPO_back/handlers/auth"
	boards_handlers "RPO_back/handlers/boards"
	user_handlers "RPO_back/handlers/users"
	"RPO_back/internal/pkg/middleware"
	"RPO_back/internal/pkg/utils/environment"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
)

type FormatterHook struct {
	Writer    io.Writer
	LogLevels []logrus.Level
	Formatter logrus.Formatter
}

// Инициализировать приложение
// Возвращает мультиплексор, который можно тестировать, а можно запрячь для обработки запросов
func initializeApp() (*mux.Router, error) {

	// Создаём новый маршрутизатор
	r := mux.NewRouter()

	// Применяем middleware
	r.Use(loggingMiddleware)
	r.Use(middleware.CorsMiddleware)

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
	environment.ValidateEnv()
	var logger = logrus.New()
	logger.SetOutput(io.Discard)
	logger.AddHook(&writer.Hook{
		Writer:    os.Stdout, // Цветной вывод на экран
		LogLevels: logrus.AllLevels,
	})

	file, err := os.OpenFile("log.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.AddHook(&FormatterHook{
			Writer:    file,
			LogLevels: logrus.AllLevels,
			Formatter: &logrus.JSONFormatter{}, // Запись в файл в формате JSON
		})
	} else {
		logger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to log to file, using default stderr")
	}
	err := godotenv.Load(".env")
	if err != nil {
		logger.Warn("Warning: no .env file loaded", err.Error())
		fmt.Print()
	} else {
		logger.Info(".env file loaded")
	}
	router, err := initializeApp()
	if err != nil {
		logger.Fatalf("Ошибка инициализации приложения: %v", err)
	}

	// Определяем адрес и порт для сервера
	addr := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))
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
