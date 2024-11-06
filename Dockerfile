# 1 шаг - сборка
FROM golang:1.23-alpine AS build_stage
RUN apk add make wget
WORKDIR /tmp
RUN wget https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-arm64.tar.gz -O migrate.tar.gz
RUN tar -xzf migrate.tar.gz
RUN cp ./migrate /usr/bin/migrate
RUN rm -rf migrate migrate.tar.gz
COPY ./go.mod /go/src/pumpkin/
COPY ./go.sum /go/src/pumpkin/
WORKDIR /go/src/pumpkin
RUN go mod download
COPY . /go/src/pumpkin
RUN make build

# 2 шаг
FROM alpine AS run_stage
RUN apk add make
WORKDIR /app_binary
COPY ./database /app/database
COPY ./Makefile /app/Makefile
COPY --from=build_stage /go/src/pumpkin/build/pumpkin_backend /app_binary/
COPY --from=build_stage /usr/bin/migrate /usr/bin/migrate
RUN chmod +x ./pumpkin_backend
RUN chmod +x /usr/bin/migrate
EXPOSE 8080/tcp
ENTRYPOINT ./pumpkin_backend
