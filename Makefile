DB_SSLMODE := require

# Подключаем .env файл
ifneq (,$(wildcard ./.env))
include .env
export
endif

BUILD_DIR := ./build

GOFLAGS := # Может, когда-нибудь пригодятся
LDFLAGS := -ldflags="-s -w" # Отключить дебаг-информацию

# Цели Makefile, которые не привязываются к файлам
.PHONY: build_auth build_user build_board run_tests test clean coverage run generate

build_auth:
	@echo "==> Building the application..."
	@go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/auth ./cmd/internal/auth

build_user:
	@echo "==> Building the application..."
	@go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/user ./cmd/internal/user

build_board:
	@echo "==> Building the application..."
	@go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/board ./cmd/internal/board

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

migrate-up:
	@which migrate
	@echo "==> Running migrations..."
	@echo "Provide DSN for migrations with superuser: >>> "
	@read MIGRATION_DSN; echo "MIGRATION_DSN: $$MIGRATION_DSN"; \
	migrate -path ./database/migrations -database $(MIGRATION_DSN) up

make-migrations:
	@echo "==> Let's generate migrations with Atlas!"
	@which atlas
	@echo "Provide migration name: >>> "
	@read MIGRATION_NAME; echo "MIGRATION_NAME: $$MIGRATION_NAME"; \
	atlas migrate diff $$MIGRATION_NAME.up --dir "file://database/migrations" --to "file://database/schema.sql" --dev-url "$(TEST_DATABASE_URL)"
