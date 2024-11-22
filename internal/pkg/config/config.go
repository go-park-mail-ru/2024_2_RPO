package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// Эта структура нужна, чтобы все 3 сервиса могли
// подтягивать настройки из одной среды. Соответственно,
// то, что может быть использовано несколькими сервисами,
// кладём в общий конфиг, а индивидуальные параметры - в
// конфиг для отдельно взятого сервиса
type Config struct {
	RedisDSN      string
	PostgresDSN   string
	MaxUploadSize int64
	AuthURL       string
	UploadsDir    string
	ServerPort    string // Порт для всех TCP Listen-ов, в том числе GRPC

	Auth  *AuthConfig
	User  *UserConfig
	Board *BoardConfig
}

type AuthConfig struct {
	PostgresPoolSize int
	LogFile          string
}
type UserConfig struct {
	PostgresPoolSize int
	LogFile          string
}
type BoardConfig struct {
	PostgresPoolSize int
	LogFile          string
}

var (
	CurrentConfig *Config
)

// Проверить, есть ли данные переменные в env
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

// Удостовериться, что в env лежат все необходимые для работы приложения переменные
func ValidateEnv() error {

	err := checkEnv([]string{
		"POSTGRES_URL",
		"REDIS_URL",
		"MAX_UPLOAD_SIZE",
		"AUTH_GRPC_URL",
		"UPLOAD_DIR",
		"SERVER_PORT",

		"LOG_ROOT",

		"AUTH_POSTGRES_MAX_CONNS",
		"USER_POSTGRES_MAX_CONNS",
		"BOARD_POSTGRES_MAX_CONNS",

		"AUTH_LOG_FILE",
		"USER_LOG_FILE",
		"BOARD_LOG_FILE",
	})
	if err != nil {
		return err
	}

	return nil
}

func stringToInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 32)
	return int(i)
}

func LoadConfig() (err error) {
	err = ValidateEnv()
	if err != nil {
		return fmt.Errorf("LoadConfig: %w", err)
	}

	CurrentConfig = &Config{}
	CurrentConfig.Auth = &AuthConfig{}
	CurrentConfig.User = &UserConfig{}
	CurrentConfig.Board = &BoardConfig{}

	logRoot := os.Getenv("LOG_ROOT")

	CurrentConfig.PostgresDSN = os.Getenv("POSTGRES_URL")
	CurrentConfig.RedisDSN = os.Getenv("REDIS_URL")
	CurrentConfig.MaxUploadSize = int64(stringToInt(os.Getenv("MAX_UPLOAD_SIZE")))
	CurrentConfig.AuthURL = os.Getenv("AUTH_GRPC_URL")
	CurrentConfig.UploadsDir = os.Getenv("UPLOAD_DIR")
	CurrentConfig.ServerPort = os.Getenv("SERVER_PORT")
	CurrentConfig.Auth.PostgresPoolSize = stringToInt(os.Getenv("AUTH_POSTGRES_MAX_CONNS"))
	CurrentConfig.User.PostgresPoolSize = stringToInt(os.Getenv("USER_POSTGRES_MAX_CONNS"))
	CurrentConfig.Board.PostgresPoolSize = stringToInt(os.Getenv("BOARD_POSTGRES_MAX_CONNS"))
	CurrentConfig.Auth.LogFile = filepath.Join(logRoot, os.Getenv("AUTH_LOG_FILE"))
	CurrentConfig.User.LogFile = filepath.Join(logRoot, os.Getenv("USER_LOG_FILE"))
	CurrentConfig.Board.LogFile = filepath.Join(logRoot, os.Getenv("BOARD_LOG_FILE"))

	return nil
}
