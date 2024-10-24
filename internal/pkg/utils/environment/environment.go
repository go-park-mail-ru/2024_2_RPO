package environment

import (
	"errors"
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

	err := checkEnv([]string{"DATABASE_URL",
		"SERVER_PORT",
		"REDIS_URL",
		"CORS_ORIGIN"})
	if err != nil {
		return err
	}

	if _, err = strconv.ParseInt(os.Getenv("SERVER_PORT"), 10, 64); err != nil {
		return errors.New("SERVER_PORT env variable is invalid")
	}

	return nil
}
