package main

import (
	"RPO_back/internal/pkg/config"
	"RPO_back/internal/pkg/middleware/cors"
	"RPO_back/internal/pkg/middleware/csrf"
	"RPO_back/internal/pkg/middleware/logging_middleware"
	"RPO_back/internal/pkg/middleware/no_panic"
	UserDelivery "RPO_back/internal/pkg/user/delivery"
	UserRepository "RPO_back/internal/pkg/user/repository"
	UserUsecase "RPO_back/internal/pkg/user/usecase"
	"RPO_back/internal/pkg/utils/logging"
	"net/http"

	"context"
	"fmt"
	"os"

	AuthGRPC "RPO_back/internal/pkg/auth/delivery/grpc/gen"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
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

	// Формирование конфига
	err = config.LoadConfig()
	if err != nil {
		log.Fatalf("environment configuration is invalid: %w", err)
		return
	}

	// Подключение к PostgreSQL
	postgresDb, err := pgxpool.New(context.Background(), config.CurrentConfig.PostgresDSN)
	if err != nil {
		log.Error("error connecting to PostgreSQL: ", err)
		return
	}
	defer postgresDb.Close()

	// Проверка подключения к PostgreSQL
	if err = postgresDb.Ping(context.Background()); err != nil {
		log.Fatal("error while pinging PostgreSQL: ", err)
	}

	// Подключение к GRPC сервису авторизации
	grpcAddr := fmt.Sprintf("%s:%s", config.CurrentConfig.AuthGRPCHost, config.CurrentConfig.AuthGRPCPort)
	conn, err := grpc.NewClient(grpcAddr)
	authGRPC := AuthGRPC.NewAuthClient(conn)
	// Проверка подключения к GRPC
	sess := &AuthGRPC.CheckSessionRequest{SessionID: "12345678"}
	authGRPC.CheckSession(context.Background(), sess)

	// User
	userRepository := UserRepository.CreateUserRepository(postgresDb)
	userUsecase := UserUsecase.CreateUserUsecase(userRepository)
	userDelivery := UserDelivery.CreateUserDelivery(userUsecase)

	// Создаём новый маршрутизатор
	router := mux.NewRouter()

	// Применяем middleware
	router.Use(no_panic.PanicMiddleware)
	router.Use(logging_middleware.LoggingMiddleware)
	router.Use(cors.CorsMiddleware)
	router.Use(csrf.CSRFMiddleware)

	// Регистрируем обработчики
	router.HandleFunc("/auth/register", userDelivery.RegisterUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/auth/login", userDelivery.LoginUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/auth/logout", userDelivery.LogoutUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/auth/changePassword", userDelivery.ChangePassword).Methods("POST", "OPTIONS")
	router.HandleFunc("/users/me", userDelivery.GetMyProfile).Methods("GET", "OPTIONS")
	router.HandleFunc("/users/me", userDelivery.UpdateMyProfile).Methods("PUT", "OPTIONS")
	router.HandleFunc("/users/me/avatar", userDelivery.SetMyAvatar).Methods("PUT", "OPTIONS")

	// Запускаем сервер
	addr := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))
	log.Infof("server started at http://0.0.0.0%s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("error while starting server: %w", err)
	}
}
