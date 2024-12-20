package main

import (
	AuthDelivery "RPO_back/internal/pkg/auth/delivery"
	"RPO_back/internal/pkg/auth/delivery/grpc/gen"
	AuthRepository "RPO_back/internal/pkg/auth/repository"
	AuthUsecase "RPO_back/internal/pkg/auth/usecase"
	"RPO_back/internal/pkg/config"
	"RPO_back/internal/pkg/middleware/logging_middleware"
	"RPO_back/internal/pkg/middleware/performance"
	"RPO_back/internal/pkg/utils/logging"
	"RPO_back/internal/pkg/utils/misc"
	"net"
	"net/http"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Формирование конфига
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("environment configuration is invalid: %v", err)
		return
	}

	// Настройка движка логов
	logsFile, err := os.OpenFile(config.CurrentConfig.Auth.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Error while opening log file %s: %s\n", config.CurrentConfig.Auth.LogFile, err.Error())
		return
	}
	defer logsFile.Close()
	logging.SetupLogger(logsFile)
	log.Info("Config: ", fmt.Sprintf("%#v", config.CurrentConfig))

	// Подключение к PostgreSQL
	postgresDB, err := misc.ConnectToPgx(config.CurrentConfig.Auth.PostgresPoolSize)
	if err != nil {
		log.Error("error connecting to PostgreSQL: ", err)
		return
	}
	defer postgresDB.Close()

	//Подключение к Redis
	redisOpts, err := redis.ParseURL(config.CurrentConfig.RedisDSN)
	if err != nil {
		log.Fatal("error connecting to Redis: ", err)
		return
	}
	redisDB := redis.NewClient(redisOpts)
	defer redisDB.Close()

	// Проверка подключения к Redis
	if pingStatus := redisDB.Ping(redisDB.Context()); pingStatus == nil || pingStatus.Err() != nil {
		if pingStatus != nil {
			log.Fatal("error while pinging Redis: ", pingStatus.Err())
		} else {
			log.Fatal("unknown error while pinging Redis")
		}
		return
	}

	grpcMetricsMiddleware, err := performance.CreateGRPCPerformanceMiddleware("auth")
	if err != nil {
		log.Error("error creating gRPC performance middleware: ", err)
	}

	// Auth
	authRepository := AuthRepository.CreateAuthRepository(postgresDB, redisDB)
	authUsecase := AuthUsecase.CreateAuthUsecase(authRepository)
	authDelivery := AuthDelivery.CreateAuthServer(authUsecase)

	if authDelivery == nil {
		panic("authDelivery is nil")
	}

	LogMiddleware := logging_middleware.CreateGrpcLogMiddleware(log.StandardLogger())

	grpcServer := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(grpcMetricsMiddleware.GRPCMetricsInterceptor, LogMiddleware.InterceptorLogger),
	)
	gen.RegisterAuthServer(grpcServer, authDelivery)

	router := mux.NewRouter()
	router.Use(logging_middleware.LoggingMiddleware)
	router.HandleFunc("/prometheus/metrics", promhttp.Handler().ServeHTTP)

	go func() {
		if err := http.ListenAndServe(":8087", router); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		listener, err := net.Listen("tcp4", ":"+config.CurrentConfig.ServerPort)
		if err != nil {
			log.Fatalf("failed to listen on port %s: %v", config.CurrentConfig.ServerPort, err)
		}
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve gRPC server: %v", err)
		}
		log.Infof("gRPC server is listening on port %s", config.CurrentConfig.ServerPort)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	os.Exit(0)
}
