APP_NAME := pumpkin_backend
BUILD_DIR := ./build
SRC_DIR := ./cmd/pumpkin_backend

# Определите стандартные флаги сборки
GOFLAGS :=
LDFLAGS := -ldflags="-s -w" # Отключить дебаг-информацию

# Цели Makefile, которые не привязываются к файлам
.PHONY: all build test clean coverage

all: build

build:
	@echo "==> Building the application..."
	@GOOS=linux GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(SRC_DIR)

test:
	@echo "==> Running tests..."
	@go test $(GOFLAGS) -v ./...

coverage:
	@echo "==> Calculating coverage..."

clean:
	@echo "==> Cleaning up..."
	@rm -rf $(BUILD_DIR)

docker-build:
	@echo "==> Building Docker image..."
	@docker build -t $(APP_NAME):latest .

# Миграции для базы данных
migrate-up:
	@echo "==> Running migrations..."
	@migrate -path ./migrations -database YOUR_DB_CONNECTION_URL up

migrate-down:
	@echo "==> Reverting migrations..."
	@migrate -path ./migrations -database YOUR_DB_CONNECTION_URL down
