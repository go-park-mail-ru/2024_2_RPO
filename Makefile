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
.PHONY: all build run_tests test clean coverage run generate

all: build

build:
	@echo "==> Building the application..."
	@go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(SRC_DIR)

run_tests: generate
	@echo "==> Running tests..."
	@go test $(GOFLAGS) -coverprofile coverage_raw.out -v ./...

test: run_tests
	@echo "==> Calculating coverage..."
	@grep -vi "mock" coverage_raw.out | cat >coverage.out
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out -o=coverage.html
	@echo "==> Done! Check coverage.html file!"

generate:
	@echo "==> Generating mocks and protobuf..."
	@go generate ./...

clean:
	@echo "==> Cleaning up..."
	@rm -rf $(BUILD_DIR)

run:
	@go run ${SRC_DIR}

migrate-up:
	@echo "==> Running migrations..."
	@migrate -path ./database/migrations -database $(DATABASE_URL) up

make-migrations:
	@echo "==> Let's generate migrations with Atlas!"
	@which atlas
	@echo "Provide migration name: >>> "
	@read MIGRATION_NAME; echo "MIGRATION_NAME: $$MIGRATION_NAME"; \
	atlas migrate diff $$MIGRATION_NAME.up --dir "file://database/migrations" --to "file://database/schema.sql" --dev-url "$(TEST_DATABASE_URL)"
