package utils

import (
	"RPO_back/internal/pkg/config"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Проверить, все ли необходимые для работы переменные окружения заданы
func checkEnv(envVars []string) error {

	var missingVars []string

	for _, envVar := range envVars {
		if value, exists := os.LookupEnv(envVar); !exists || value == "" {
			missingVars = append(missingVars, envVar)
		}
	}

	if len(missingVars) > 0 {
		return fmt.Errorf("error: this env vars are missing: %v", missingVars)
	} else {
		return nil
	}
}

// Загрузить файл .env и вытянуть оттуда данные
func LoadDotEnv() (config_ *config.ServerConfig, err_ error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("can't load .env file: %s", err.Error())
	} 

	err = checkEnv([]string{"DB_URL", "SERVER_PORT", "REDIS_URL"})
	if err != nil {
		return nil, err
	}

	port, err1 := strconv.ParseInt(os.Getenv("SERVER_PORT"), 10, 64)
	if err1 != nil {
		return nil, fmt.Errorf("error converting SERVER_PORT to int: %v", err1)
	}
	
	var ret config.ServerConfig
    ret.DbUrl = os.Getenv("DB_URL")
    ret.ServerPort = int(port)
    ret.RedisUrl = os.Getenv("REDIS_URL")

	return &ret, nil
}
