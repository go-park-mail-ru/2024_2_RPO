package main

import (
	user_handlers "RPO_back/handlers/users"
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
	// Создаём новый маршрутизатор
	mux := http.NewServeMux()

	// Регистрируем обработчики
	mux.HandleFunc("/ping", user_handlers.GetUser)
	mux.HandleFunc("/hello", helloHandler)

	// Определяем адрес и порт для сервера
	addr := ":8080"
	fmt.Printf("Сервер запущен на %s\n", addr)

	// Запускаем сервер
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
