DB_SSLMODE := require

# Подключаем .env файл
ifneq (,$(wildcard ./.env))
include .env
export
endif

APP_NAME := pumpkin_backend
BUILD_DIR := ./build
SRC_DIR := ./cmd/pumpkin_backend

GOFLAGS := # Может, когда-нибудь пригодятся
LDFLAGS := -ldflags="-s -w" # Отключить дебаг-информацию

DATABASE_URL="postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSLMODE}"

# Цели Makefile, которые не привязываются к файлам
.PHONY: all build test clean coverage run

all: build

build:
	@echo "==> Building the application..."
	@go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(SRC_DIR)

run_tests:
	@echo "==> Running tests..."
	@go test $(GOFLAGS) -coverprofile coverage_raw.out -v ./...

test: run_tests
	@echo "==> Calculating coverage..."
	@grep -vi "mock" coverage_raw.out | cat >coverage.out
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out -o=coverage.html
	@echo "==> Done! Check coverage.html file!"


clean:
	@echo "==> Cleaning up..."
	@rm -rf $(BUILD_DIR)

run:
	@go run ${SRC_DIR}

docker-build:
	@echo "==> Building Docker image..."
	@docker build -t $(APP_NAME):latest .

# Миграции для базы данных
migrate-up:
	@echo "==> Running migrations..."
	@migrate -path ./database/migrations -database $(DATABASE_URL) up

migrate-down:
	@echo "==> Reverting migrations..."
	@migrate -path ./database/migrations -database $(DATABASE_URL) down
