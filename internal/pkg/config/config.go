package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	RedisDSN    string
	PostgresDSN string
	Auth        *AuthConfig
	User        *UserConfig
	Board       *BoardConfig
}

type AuthConfig struct {
	PostgresPoolSize int
	Port             int
	LogFile          string
}
type UserConfig struct {
	PostgresPoolSize int
	Port             int
	LogFile          string
}
type BoardConfig struct {
	PostgresPoolSize int
	Port             int
	LogFile          string
}

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

	err := checkEnv([]string{"POSTGRES_HOST",
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
