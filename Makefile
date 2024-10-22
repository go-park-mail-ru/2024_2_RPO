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

# Цели Makefile, которые не привязываются к файлам
.PHONY: all build test clean coverage run

all: build

build:
	@echo "==> Building the application..."
	@go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(SRC_DIR)

test:
	@echo "==> Running tests..."
	@go test $(GOFLAGS) -coverprofile coverage.out -v ./...

coverage: test
	@echo "==> Calculating coverage..."
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
	@migrate -path ./database/migrations -database $(DB_URL) up

migrate-down:
	@echo "==> Reverting migrations..."
	@migrate -path ./database/migrations -database $(DB_URL) down