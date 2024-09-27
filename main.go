package main

import (
	"RPO_back/auth"
	"RPO_back/database"
	auth_handlers "RPO_back/handlers/auth"
	user_handlers "RPO_back/handlers/users"
	"RPO_back/utils"
	"fmt"
	"log"
	"net/http"
)

// pingHandler обрабатывает запросы к /ping
func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world")
}

// helloHandler обрабатывает запросы к /hello
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world")
}

func main() {
	// Создаём логгер
	logger := log.Default()

	// Обрабатываем файл .env
	serverConfig, err := utils.LoadDotEnv()
	if err != nil {
		log.Fatalf(err.Error())
	}
	logger.Printf("Server config: %#v", serverConfig)

	// Подключаемся к базе
	err2 := database.ConnectToDb(serverConfig.DbPort, serverConfig.DbUser, serverConfig.DbPasswd)
	if err2 != nil {
		log.Fatal(err2.Error())
	}

	// Подключаемся к Redis
	err3 := auth.ConnectToRedis(serverConfig.RedisPort, serverConfig.RedisUser, serverConfig.RedisPasswd)
	if err3 != nil {
		log.Fatal(err3.Error())
	}

	// Создаём новый маршрутизатор
	mux := http.NewServeMux()

	// Регистрируем обработчики
	mux.HandleFunc("/ping", user_handlers.GetUser)
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/auth/register", auth_handlers.RegisterUser)

	// Определяем адрес и порт для сервера
	addr := fmt.Sprintf(":%d", serverConfig.ServerPort)
	fmt.Printf("Сервер запущен на http://localhost%s\n", addr)

	// Запускаем сервер
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}