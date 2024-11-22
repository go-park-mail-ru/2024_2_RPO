package config

import (
	"errors"
	"fmt"
	"os"
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
	ServerPort    string

	Auth  *AuthConfig
	User  *UserConfig
	Board *BoardConfig
}

type AuthConfig struct {
	PostgresPoolSize int
	LogFile          string
	GrpcServerPort   string
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
		"POSTGRES_HOST",
		"POSTGRES_PORT",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"POSTGRES_DB",
		"POSTGRES_SSLMODE",

		"SERVER_PORT",

		"REDIS_HOST",
		"REDIS_PORT",
		"REDIS_PASSWORD",
		"CORS_ORIGIN",
		"LOGS_FILE",
	})
	if err != nil {
		return err
	}

	if _, err = strconv.ParseInt(os.Getenv("SERVER_PORT"), 10, 64); err != nil {
		return errors.New("SERVER_PORT env variable is invalid")
	}

	return nil
}

func LoadConfig() (err error) {
	return fmt.Errorf("config not implemented")
}
