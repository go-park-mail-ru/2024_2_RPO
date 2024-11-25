# 1 шаг - сборка
FROM golang:1.23-alpine AS build_stage
RUN apk add make wget
WORKDIR /tmp
RUN wget https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-arm64.tar.gz -O migrate.tar.gz
RUN tar -xzf migrate.tar.gz
RUN cp ./migrate /usr/bin/migrate
RUN rm -rf migrate migrate.tar.gz
COPY ./go.mod ./go.sum /go/src/pumpkin/
WORKDIR /go/src/pumpkin
RUN go mod download
COPY ./Makefile /go/src/pumpkin/
COPY ./internal /go/src/pumpkin/internal/
COPY ./cmd/auth/main.go /go/src/pumpkin/cmd/auth/
RUN make build_auth

# 2 шаг
FROM ubuntu AS run_stage
RUN apt update && apt install make -y
WORKDIR /app_binary
COPY ./database /app/database
COPY ./Makefile /app/Makefile
COPY --from=build_stage /go/src/pumpkin/build/auth_app /app_binary/
COPY --from=build_stage /usr/bin/migrate /usr/bin/migrate
RUN chmod +x ./auth_app
RUN chmod +x /usr/bin/migrate
EXPOSE 8080/tcp
ENTRYPOINT ./auth_app
