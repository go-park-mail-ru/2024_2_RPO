package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Тест для функции checkEnv
func TestCheckEnv(t *testing.T) {
	tests := []struct {
		name        string
		envSetup    map[string]string
		envVars     []string
		expectedErr bool
	}{
		{
			name:        "all env vars present",
			envSetup:    map[string]string{"VAR1": "value1", "VAR2": "value2"},
			envVars:     []string{"VAR1", "VAR2"},
			expectedErr: false,
		},
		{
			name:        "missing env vars",
			envSetup:    map[string]string{"VAR1": "value1"},
			envVars:     []string{"VAR1", "VAR2"},
			expectedErr: true,
		},
		{
			name:        "empty env var",
			envSetup:    map[string]string{"VAR1": "", "VAR2": "value2"},
			envVars:     []string{"VAR1", "VAR2"},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envSetup {
				os.Setenv(k, v)
			}
			err := checkEnv(tt.envVars)
			if (err != nil) != tt.expectedErr {
				t.Errorf("checkEnv() error = %v, expected error = %v", err, tt.expectedErr)
			}
			for k := range tt.envSetup {
				os.Unsetenv(k)
			}
		})
	}
}

// Тест для функции ValidateEnv
func TestValidateEnv(t *testing.T) {
	// Настройка переменных окружения
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("POSTGRES_PORT", "5432")
	os.Setenv("POSTGRES_USER", "user")
	os.Setenv("POSTGRES_PASSWORD", "password")
	os.Setenv("POSTGRES_DB", "database")
	os.Setenv("POSTGRES_SSLMODE", "disable")
	os.Setenv("SERVER_PORT", "8000")
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6379")
	os.Setenv("REDIS_PASSWORD", "redispassword")
	os.Setenv("CORS_ORIGIN", "*")
	os.Setenv("LOGS_FILE", "logs.txt")

	err := ValidateEnv()
	assert.Error(t, err)
}
