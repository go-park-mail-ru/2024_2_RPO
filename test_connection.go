package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TODO сделать интеграционное тестирование
func TestRegisterUser(t *testing.T) {
	router, err := initializeApp()
	if err != nil {
		t.Fatalf("Ошибка инициализации приложения: %v", err)
	}

	// Создаём тестовый запрос
	req, err := http.NewRequest("POST", "/auth/register", nil)
	if err != nil {
		t.Fatalf("Ошибка создания запроса: %v", err)
	}

	// Создаём ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()

	// Отправляем запрос
	router.ServeHTTP(rr, req)

	// Проверяем статус код
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Ожидался статус %v, получен %v", http.StatusOK, status)
	}

	// Дополнительные проверки ответа
}
