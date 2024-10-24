package main

import (
	auth_handlers "RPO_back/handlers/auth"
	boards_handlers "RPO_back/handlers/boards"
	user_handlers "RPO_back/handlers/users"
	"RPO_back/internal/pkg/middleware/cors"
	"RPO_back/internal/pkg/middleware/logging_middleware"
	"RPO_back/internal/pkg/utils/environment"
	"RPO_back/internal/pkg/utils/logging"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

// Инициализировать приложение
// Возвращает мультиплексор, который можно тестировать, а можно запрячь для обработки запросов
func initializeApp() http.Handler {

	// Создаём новый маршрутизатор
	r := mux.NewRouter()

	// Применяем middleware
	r.Use(cors.CorsMiddleware)
	r.Use(logging_middleware.LoggingMiddleware)

	// Регистрируем обработчики
	r.HandleFunc("/auth/register", auth_handlers.RegisterUser).Methods("POST", "OPTIONS")
	r.HandleFunc("/users/me", user_handlers.GetMe).Methods("GET", "OPTIONS")
	r.HandleFunc("/boards/my", boards_handlers.GetMyBoardsHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/boards", boards_handlers.CreateBoardHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/boards/{boardId}", boards_handlers.DeleteBoardHandler).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/auth/login", auth_handlers.LoginUser).Methods("POST", "OPTIONS")
	r.HandleFunc("/auth/logout", auth_handlers.LogoutUser).Methods("POST", "OPTIONS")

	return r
}

func main() {
	// Настройка движка логов
	if _, exists := os.LookupEnv("LOGS_FILE"); exists == false {
		fmt.Printf("You should provide log file env variable: LOGS_FILE\n")
		return
	}
	logsFile, err := os.OpenFile(os.Getenv("LOGS_FILE"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Error while opening log file %s: %s\n", os.Getenv("LOGS_FILE"), err.Error())
		return
	}
	defer logsFile.Close()
	logging.SetupLogger(logsFile)

	// Загрузка переменных окружения
	err = godotenv.Load(".env")
	if err != nil {
		log.Warn("warning: no .env file loaded", err.Error())
		fmt.Print()
	} else {
		log.Info(".env file loaded")
	}

	// Проверка переменных окружения
	err = environment.ValidateEnv()
	if err != nil {
		log.Fatalf("environment configuration is invalid: %s", err.Error())
		return
	}

	// Подключение к PostgreSQL
	postgresDb, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Error("error connecting to postgres: ", err)
		return
	}
	defer postgresDb.Close()

	// Проверка подключения к PostgreSQL
	if err = postgresDb.Ping(context.Background()); err != nil {
		log.Fatal("error while pinging PostgreSQL: ", err)
	}

	//Подключение к Redis
	redisOpts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal("error connecting to Redis: ", err)
		return
	}
	redisDb := redis.NewClient(redisOpts)
	defer redisDb.Close()

	// Проверка подключения к Redis

	if pingStatus := redisDb.Ping(redisDb.Context()); pingStatus == nil || pingStatus.Err() != nil {
		if pingStatus != nil {
			log.Fatal("error while pinging Redis: ", pingStatus.Err())
		} else {
			log.Fatal("unknown error while pinging Redis")
		}
		return
	}

	router := initializeApp()

	// Определяем адрес и порт для сервера
	addr := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))
	log.Infof("server started at http://localhost%s", addr)

	// Запускаем сервер
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("error whie starting server: %v", err)
	}
}
