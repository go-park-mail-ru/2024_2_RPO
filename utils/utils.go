package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	DbPort      int
	DbUser      string
	DbPasswd    string
	ServerPort  int
	RedisPort   int
	RedisUser   string
	RedisPasswd string
}

// Проверить, все ли необходимые для работы переменные окружения заданы
func checkEnv(envVars []string) error {

	var missingVars []string

	for _, envVar := range envVars {
		if value, exists := os.LookupEnv(envVar); !exists || value == "" {
			missingVars = append(missingVars, envVar)
		}
	}

	if len(missingVars) > 0 {
		return errors.New(fmt.Sprintf("Error: this env vars are missing: %v\n", missingVars))
	} else {
		return nil
	}
}

// Загрузить файл .env и вытянуть оттуда данные
func LoadDotEnv() (config *ServerConfig, err error) {
	err2 := godotenv.Load()
	if err2 != nil {
		return nil, errors.New("Cant load .env file")
	}

	err3 := checkEnv([]string{"DB_PASSWORD", "DB_USER", "DB_PORT", "SERVER_PORT", "REDIS_PORT", "REDIS_USER", "REDIS_PASSWORD"})
	if err3 != nil {
		return nil, err3
	}

	dbPort_raw := os.Getenv("DB_PORT")
	serverPort_raw := os.Getenv("SERVER_PORT")
	redisPort_raw := os.Getenv("REDIS_PORT")
	dbPort, err4 := strconv.ParseInt(dbPort_raw, 10, 64)
	serverPort, err5 := strconv.ParseInt(serverPort_raw, 10, 64)
	redisPort, err6 := strconv.ParseInt(redisPort_raw, 10, 64)

	if err4 != nil || err5 != nil || err6 != nil {
		return nil, errors.New("DB port or server port cant be parsed")
	}
	if dbPort < 0 || dbPort > 65535 {
		return nil, errors.New("DB port is out of range")
	}
	if redisPort < 0 || redisPort > 65535 {
		return nil, errors.New("Redis port is out of range")
	}
	if serverPort < 0 || serverPort > 65535 {
		return nil, errors.New("Server port is out of range")
	}
	if serverPort < 1000 {
		return nil, errors.New("Server port should be less than 1000 because of no need to root")
	}

	var ret ServerConfig
	ret.DbPasswd = os.Getenv("DB_PASSWORD")
	ret.DbUser = os.Getenv("DB_USER")
	ret.RedisUser = os.Getenv("REDIS_USER")
	ret.RedisPasswd = os.Getenv("REDIS_PASSWORD")
	ret.ServerPort = int(serverPort)
	ret.RedisPort = int(redisPort)
	ret.DbPort = int(dbPort)

	return &ret, nil
}
