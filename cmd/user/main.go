package main

import (
	"RPO_back/internal/pkg/config"
	"RPO_back/internal/pkg/middleware/cors"
	"RPO_back/internal/pkg/middleware/csrf"
	"RPO_back/internal/pkg/middleware/logging_middleware"
	"RPO_back/internal/pkg/middleware/no_panic"
	"RPO_back/internal/pkg/middleware/performance"
	"RPO_back/internal/pkg/middleware/session"
	UserDelivery "RPO_back/internal/pkg/user/delivery"
	UserRepository "RPO_back/internal/pkg/user/repository"
	UserUsecase "RPO_back/internal/pkg/user/usecase"
	"RPO_back/internal/pkg/utils/logging"
	"RPO_back/internal/pkg/utils/misc"
	"net/http"
	"time"

	"context"
	"fmt"
	"os"

	AuthGRPC "RPO_back/internal/pkg/auth/delivery/grpc/gen"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Костыль
	log.Info("Sleeping 10 seconds waiting Postgres to start...")
	time.Sleep(10 * time.Second)

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
	postgresDB, err := misc.ConnectToPgx(config.CurrentConfig.User.PostgresPoolSize)
	if err != nil {
		log.Fatal("error connecting to PostgreSQL: ", err)
		return
	}
	defer postgresDB.Close()

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
	userRepository := UserRepository.CreateUserRepository(postgresDB)
	userUsecase := UserUsecase.CreateUserUsecase(userRepository, authGRPC)
	userDelivery := UserDelivery.CreateUserDelivery(userUsecase)

	// Создаём новый маршрутизатор
	router := mux.NewRouter()

	// Применяем middleware
	router.Use(no_panic.PanicMiddleware)
	pm, err := performance.CreateHTTPPerformanceMiddleware("user")
	if err != nil {
		log.Fatal("create HTTP middleware: ", err)
	}
	router.Use(pm.Middleware)
	router.Use(logging_middleware.LoggingMiddleware)
	router.Use(cors.CorsMiddleware)
	router.Use(csrf.CSRFMiddleware)
	sm := session.CreateSessionMiddleware(authGRPC)
	router.Use(sm.Middleware)

	// Регистрируем обработчики
	router.HandleFunc("/prometheus/metrics", promhttp.Handler().ServeHTTP)
	router.HandleFunc("/auth/register", userDelivery.RegisterUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/auth/login", userDelivery.LoginUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/auth/logout", userDelivery.LogoutUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/auth/changePassword", userDelivery.ChangePassword).Methods("POST", "OPTIONS")
	router.HandleFunc("/users/me", userDelivery.GetMyProfile).Methods("GET", "OPTIONS")
	router.HandleFunc("/users/me", userDelivery.UpdateMyProfile).Methods("PUT", "OPTIONS")
	router.HandleFunc("/users/me/avatar", userDelivery.SetMyAvatar).Methods("PUT", "OPTIONS")

	// Регистрируем обработчик Prometheus
	metricsRouter := mux.NewRouter()
	metricsRouter.Handle("/prometheus/metrics", promhttp.Handler())

	// Объявляем серверы
	mainAddr := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))
	promAddr := ":8087"
	mainServer := &http.Server{
		Addr:    mainAddr,
		Handler: router,
	}
	promServer := &http.Server{
		Addr:    promAddr,
		Handler: metricsRouter,
	}

	misc.StartServers(mainServer, promServer)
}
