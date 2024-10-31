package main

import (
	AuthDelivery "RPO_back/internal/pkg/auth/delivery"
	AuthRepository "RPO_back/internal/pkg/auth/repository"
	AuthUsecase "RPO_back/internal/pkg/auth/usecase"
	BoardDelivery "RPO_back/internal/pkg/board/delivery"
	BoardRepository "RPO_back/internal/pkg/board/repository"
	BoardUsecase "RPO_back/internal/pkg/board/usecase"
	"RPO_back/internal/pkg/middleware/cors"
	"RPO_back/internal/pkg/middleware/csrf"
	"RPO_back/internal/pkg/middleware/logging_middleware"
	"RPO_back/internal/pkg/middleware/no_panic"
	sessionMiddleware "RPO_back/internal/pkg/middleware/session"
	UserDelivery "RPO_back/internal/pkg/user/delivery"
	UserRepository "RPO_back/internal/pkg/user/repository"
	UserUsecase "RPO_back/internal/pkg/user/usecase"
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

	//Составление URL подключения
	os.Setenv("DATABASE_URL", fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_SSLMODE"),
	))
	os.Setenv("REDIS_URL", fmt.Sprintf("redis://:%s@%s:%s",
		os.Getenv("REDIS_PASSWORD"),
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
	))

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

	// Auth
	authRepository := AuthRepository.CreateAuthRepository(postgresDb, redisDb)
	authUsecase := AuthUsecase.CreateAuthUsecase(authRepository)
	authDelivery := AuthDelivery.CreateAuthDelivery(authUsecase)

	//Board
	boardRepository := BoardRepository.CreateBoardRepository(postgresDb)
	boardUsecase := BoardUsecase.CreateBoardUsecase(boardRepository)
	boardDelivery := BoardDelivery.CreateBoardDelivery(boardUsecase)

	//User
	userRepository := UserRepository.CreateUserRepository(postgresDb)
	userUsecase := UserUsecase.CreateUserUsecase(userRepository)
	userDelivery := UserDelivery.CreateUserDelivery(userUsecase)

	// Создаём новый маршрутизатор
	router := mux.NewRouter()

	// Применяем middleware
	router.Use(no_panic.PanicMiddleware)
	router.Use(logging_middleware.LoggingMiddleware)
	router.Use(csrf.CSRFMiddleware)
	router.Use(cors.CorsMiddleware)
	sessionMWare := sessionMiddleware.CreateSessionMiddleware(authRepository)
	router.Use(sessionMWare.Middleware)

	// Регистрируем обработчики
	router.HandleFunc("/auth/register", authDelivery.RegisterUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/auth/login", authDelivery.LoginUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/auth/logout", authDelivery.LogoutUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/auth/changePassword", authDelivery.ChangePassword).Methods("POST", "OPTIONS")
	router.HandleFunc("/users/me", userDelivery.GetMyProfile).Methods("GET", "OPTIONS")
	router.HandleFunc("/users/me", userDelivery.UpdateMyProfile).Methods("PUT", "OPTIONS")
	router.HandleFunc("/users/me/avatar", userDelivery.SetMyAvatar).Methods("PUT", "OPTIONS")
	router.HandleFunc("/boards/my", boardDelivery.GetMyBoards).Methods("GET", "OPTIONS")
	router.HandleFunc("/boards", boardDelivery.CreateNewBoard).Methods("POST", "OPTIONS")
	router.HandleFunc("/boards/{boardId}", boardDelivery.DeleteBoard).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/userPermissions/{boardId}", boardDelivery.GetMembersPermissions).Methods("GET", "OPTIONS")
	router.HandleFunc("/userPermissions/{boardId}", boardDelivery.AddMember).Methods("POST", "OPTIONS")
	router.HandleFunc("/userPermissions/{boardId}/{userId}", boardDelivery.UpdateMemberRole).Methods("PUT", "OPTIONS")
	router.HandleFunc("/userPermissions/{boardId}/{userId}", boardDelivery.RemoveMember).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/cards/{boardId}/allContent", boardDelivery.GetBoardContent).Methods("GET", "OPTIONS")
	router.HandleFunc("/cards/{boardId}", boardDelivery.CreateNewCard).Methods("POST", "OPTIONS")
	router.HandleFunc("/cards/{boardId}/{cardId}", boardDelivery.UpdateCard).Methods("PATCH", "OPTIONS")
	router.HandleFunc("/cards/{boardId}/{cardId}", boardDelivery.DeleteCard).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/columns/{boardId}", boardDelivery.CreateColumn).Methods("POST", "OPTIONS")
	router.HandleFunc("/columns/{boardId}/{columnId}", boardDelivery.UpdateColumn).Methods("PUT", "OPTIONS")
	router.HandleFunc("/columns/{boardId}/{columnId}", boardDelivery.DeleteColumn).Methods("DELETE", "OPTIONS")

	// Запускаем сервер
	addr := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))
	log.Infof("server started at http://localhost%s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("error while starting server: %v", err)
	}
}
