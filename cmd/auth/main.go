package main

import (
	AuthDelivery "RPO_back/internal/pkg/auth/delivery"
	"RPO_back/internal/pkg/auth/delivery/grpc/gen"
	AuthRepository "RPO_back/internal/pkg/auth/repository"
	AuthUsecase "RPO_back/internal/pkg/auth/usecase"
	"RPO_back/internal/pkg/config"
	"RPO_back/internal/pkg/middleware/logging_middleware"
	"RPO_back/internal/pkg/utils/logging"
	"net"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
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
	postgresDB, err := pgxpool.New(context.Background(), config.CurrentConfig.PostgresDSN)
	if err != nil {
		log.Error("error connecting to PostgreSQL: ", err)
		return
	}
	defer postgresDB.Close()

	// Проверка подключения к PostgreSQL
	if err = postgresDB.Ping(context.Background()); err != nil {
		log.Fatal("error while pinging PostgreSQL: ", err)
	}

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

	// Auth
	authRepository := AuthRepository.CreateAuthRepository(postgresDB, redisDB)
	authUsecase := AuthUsecase.CreateAuthUsecase(authRepository)
	authDelivery := AuthDelivery.CreateAuthServer(authUsecase)

	LogMiddleware := logging_middleware.CreateGrpcLogMiddleware(log.StandardLogger())

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(LogMiddleware.InterceptorLogger),
	)

	gen.RegisterAuthServer(grpcServer, authDelivery)

	listener, err := net.Listen("tcp", ":"+os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", os.Getenv("SERVER_PORT"), err)
	}
	log.Infof("gRPC server is listening on port %s", os.Getenv("SERVER_PORT"))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		log.Info("Shutting down gRPC server...")
		grpcServer.GracefulStop()
	}()

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC server: %v", err)
	}
}
