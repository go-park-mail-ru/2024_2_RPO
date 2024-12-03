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
.PHONY: build_auth build_user build_board build_poll build_all run_tests test clean coverage run generate

build_auth:
	@echo "==> Building auth application..."
	@GOOS=linux go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/auth_app ./cmd/auth

build_user:
	@echo "==> Building user application..."
	@GOOS=linux go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/user_app ./cmd/user

build_board:
	@echo "==> Building board application..."
	@GOOS=linux go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/board_app ./cmd/board

build_poll:
	@echo "==> Building poll application..."
	@GOOS=linux go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/poll_app ./cmd/poll

build_all: build_auth build_user build_poll build_board
	@echo "==> All applications are built!"

run_tests:
	@echo "==> Running tests..."
	@go test $(GOFLAGS) -coverprofile ../coverage_raw.out -v ./...

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
	@migrate -path ./migrations -database $(SUPERUSER_DSN) up

make-migrations:
	@echo "==> Let's generate migrations with Atlas!"
	@which atlas
	@echo "Provide migration name: >>> "
	@read MIGRATION_NAME; echo "MIGRATION_NAME: $$MIGRATION_NAME"; \
	atlas migrate diff $$MIGRATION_NAME.up --dir "file://database/migrations" --to "file://database/schema.sql" --dev-url "$(TEST_DATABASE_URL)"
