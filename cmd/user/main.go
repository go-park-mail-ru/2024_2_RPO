package main

import (
	"RPO_back/internal/pkg/config"
	UserDelivery "RPO_back/internal/pkg/user/delivery"
	UserRepository "RPO_back/internal/pkg/user/repository"
	UserUsecase "RPO_back/internal/pkg/user/usecase"
	"RPO_back/internal/pkg/utils/logging"

	"context"
	"fmt"
	"os"

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

	// Формирование конфига
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("environment configuration is invalid: %w", err)
		return
	}

	// Подключение к PostgreSQL
	postgresDb, err := pgxpool.New(context.Background(), config.PostgresDSN)
	if err != nil {
		log.Error("error connecting to PostgreSQL: ", err)
		return
	}
	defer postgresDb.Close()

	// Проверка подключения к PostgreSQL
	if err = postgresDb.Ping(context.Background()); err != nil {
		log.Fatal("error while pinging PostgreSQL: ", err)
	}

	// User
	userRepository := UserRepository.CreateUserRepository(postgresDb)
	userUsecase := UserUsecase.CreateUserUsecase(userRepository)
	userDelivery := UserDelivery.CreateUserDelivery(userUsecase)

	panic("GRPC needed")
}
