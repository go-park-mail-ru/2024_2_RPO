# 1 шаг - сборка
FROM golang:1.23-alpine AS build_stage
RUN apk add make
WORKDIR /tmp
COPY ./go.mod ./go.sum /go/src/pumpkin/
WORKDIR /go/src/pumpkin
RUN go mod download
COPY ./Makefile /go/src/pumpkin/
COPY ./internal /go/src/pumpkin/internal/
COPY ./cmd/user/main.go /go/src/pumpkin/cmd/user/
RUN make build_user

# 2 шаг
FROM alpine AS run_stage
WORKDIR /app_binary
COPY --from=build_stage /go/src/pumpkin/build/user_app /app_binary/
RUN chmod +x ./user_app
EXPOSE 8080/tcp
ENTRYPOINT ./user_app
