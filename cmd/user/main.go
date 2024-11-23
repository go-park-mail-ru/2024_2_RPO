package main

import (
	"RPO_back/internal/pkg/config"
	"RPO_back/internal/pkg/middleware/cors"
	"RPO_back/internal/pkg/middleware/csrf"
	"RPO_back/internal/pkg/middleware/logging_middleware"
	"RPO_back/internal/pkg/middleware/no_panic"
	"RPO_back/internal/pkg/middleware/session"
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
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Формирование конфига
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("environment configuration is invalid: %v", err)
		return
	}

	// Настройка движка логов
	logsFile, err := os.OpenFile(config.CurrentConfig.User.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Error while opening log file %s: %s\n", config.CurrentConfig.User.LogFile, err.Error())
		return
	}
	defer logsFile.Close()
	logging.SetupLogger(logsFile)
	log.Info("Config: ", fmt.Sprintf("%#v", config.CurrentConfig))

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
	grpcAddr := config.CurrentConfig.AuthURL
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("error connecting to GRPC: ", err)
	}
	authGRPC := AuthGRPC.NewAuthClient(conn)

	// Проверка подключения к GRPC
	sess := &AuthGRPC.CheckSessionRequest{SessionID: "12345678"}
	_, err = authGRPC.CheckSession(context.Background(), sess)
	if err != nil {
		log.Fatal("error while pinging GRPC: ", err)
	}

	// User
	userRepository := UserRepository.CreateUserRepository(postgresDb)
	userUsecase := UserUsecase.CreateUserUsecase(userRepository, authGRPC)
	userDelivery := UserDelivery.CreateUserDelivery(userUsecase)

	// Создаём новый маршрутизатор
	router := mux.NewRouter()

	// Применяем middleware
	router.Use(no_panic.PanicMiddleware)
	router.Use(logging_middleware.LoggingMiddleware)
	router.Use(cors.CorsMiddleware)
	router.Use(csrf.CSRFMiddleware)
	sm := session.CreateSessionMiddleware(authGRPC)
	router.Use(sm.Middleware)

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
