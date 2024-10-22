package utils

import (
	"RPO_back/internal/pkg/config"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

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

	dbPort_raw := strings.Split(strings.Split(strings.Split(os.Getenv("DB_URL"), "@")[1], ":")[1], "/")[0]
	serverPort_raw := os.Getenv("SERVER_PORT")
	redisPort_raw := strings.Split(strings.Split(os.Getenv("REDIS_URL"), "@")[1], ":")[1]
	dbPort, err1 := strconv.ParseInt(dbPort_raw, 10, 64)
	serverPort, err2 := strconv.ParseInt(serverPort_raw, 10, 64)
	redisPort, err3 := strconv.ParseInt(redisPort_raw, 10, 64)

	if err1 != nil || err2 != nil || err3 != nil {
		return nil, errors.New("DB's port or server port cant be parsed")
	}
	if dbPort < 0 || dbPort > 65535 {
		return nil, errors.New("DB port is out of range")
	}
	if redisPort < 0 || redisPort > 65535 {
		return nil, errors.New("redis port is out of range")
	}
	if serverPort < 0 || serverPort > 65535 {
		return nil, errors.New("server port is out of range")
	}
	if serverPort < 1000 {
		return nil, errors.New("server port should be less than 1000 because of no need to root")
	}

	var ret config.ServerConfig
    ret.DbUrl = os.Getenv("DB_URL")
    ret.ServerPort = int(serverPort)
    ret.RedisUrl = os.Getenv("REDIS_URL")

	return &ret, nil
}
