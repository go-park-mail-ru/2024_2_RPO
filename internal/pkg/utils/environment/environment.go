package environment

import (
	"RPO_back/internal/pkg/config"
	"fmt"
	"os"
	"strconv"
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

	err := checkEnv([]string{"DB_URL", "SERVER_PORT", "REDIS_URL"})
	if err != nil {
		return err
	}

	port, err1 := strconv.ParseInt(os.Getenv("SERVER_PORT"), 10, 64)
	if err1 != nil {
		return fmt.Errorf("error converting SERVER_PORT to int: %v", err1)
	}

	var ret config.ServerConfig
	ret.DbUrl = os.Getenv("DB_URL")
	ret.ServerPort = int(port)
	ret.RedisUrl = os.Getenv("REDIS_URL")

	return nil
}
