# 1 шаг - сборка
FROM golang:1.23-alpine AS build_stage
RUN apk add make
COPY . /go/src/pumpkin
WORKDIR /go/src/pumpkin
RUN make build

# 2 шаг
FROM alpine AS run_stage
WORKDIR /app_binary
COPY --from=build_stage /go/src/pumpkin/build/pumpkin_backend /app_binary/
RUN chmod +x ./pumpkin_backend
EXPOSE 8080/tcp
ENTRYPOINT ./pumpkin_backend
