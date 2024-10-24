package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Структуры для запросов и ответов
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

type Board struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Background  string    `json:"background"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func generateUniqueEmail() string {
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	return fmt.Sprintf("user_%s@example.com", timestamp)
}

func generateUniqueName() string {
	return fmt.Sprintf("User_%d", time.Now().UnixNano())
}

func TestIntegrationFlow(t *testing.T) {
	// TODO переработать этот тест
	return
	mux, err3 := initializeApp()
	if err3 != nil {
		t.Error("Failed init server: " + err3.Error())
		t.FailNow()
	}
	// Инициализация тестового сервера
	server := httptest.NewServer(mux) // Здесь замените nil на ваш хендлер
	defer server.Close()

	// Создание HTTP клиента с поддержкой cookie
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{
		Jar: jar,
	}

	// 1) Регистрация
	registerURL := fmt.Sprintf("%s/auth/register", server.URL)
	fmt.Println(registerURL)
	registerReq := RegisterRequest{
		Name:     generateUniqueName(),
		Email:    generateUniqueEmail(),
		Password: "securepassword123",
	}
	reqBody, _ := json.Marshal(registerReq)
	resp, err := client.Post(registerURL, "application/json", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Пользователь должен быть успешно зарегистрирован")

	// Проверка наличия куки
	u, _ := url.Parse(server.URL)
	cookies := client.Jar.Cookies(u)
	assert.NotEmpty(t, cookies, "Куки должны быть установлены после регистрации")

	// 2) Разлогинивание
	logoutURL := fmt.Sprintf("%s/auth/logout", server.URL)
	req, err := http.NewRequest("POST", logoutURL, nil)
	assert.NoError(t, err)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Пользователь должен быть успешно разлогинен")

	// Проверка удаления куки
	cookies = client.Jar.Cookies(u)
	assert.Empty(t, cookies, "Куки должны быть удалены после разлогинивания")

	// 3) Попытка логина с неверным паролем
	loginURL := fmt.Sprintf("%s/auth/login", server.URL)
	loginReq := LoginRequest{
		Email:    registerReq.Email,
		Password: "wrongpassword",
	}
	reqBody, _ = json.Marshal(loginReq)
	resp, err = client.Post(loginURL, "application/json", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Логин должен не получиться с неверным паролем")

	// 4) Логин с правильным паролем
	loginReq.Password = registerReq.Password
	reqBody, _ = json.Marshal(loginReq)
	resp, err = client.Post(loginURL, "application/json", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Логин должен быть успешным с правильным паролем")

	// Проверка наличия куки после логина
	cookies = client.Jar.Cookies(u)
	assert.NotEmpty(t, cookies, "Куки должны быть установлены после успешного логина")

	// 5) Получение информации о пользователе
	userInfoURL := fmt.Sprintf("%s/users/me", server.URL)
	resp, err = client.Get(userInfoURL)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Должна вернуться информация о пользователе")

	var userInfo User
	json.NewDecoder(resp.Body).Decode(&userInfo)
	assert.Equal(t, registerReq.Name, userInfo.Name, "Имя пользователя должно совпадать")
	assert.Equal(t, registerReq.Email, userInfo.Email, "Email пользователя должен совпадать")

	// 6) Добавление двух досок
	boardURL := fmt.Sprintf("%s/boards", server.URL)
	boards := []Board{}

	for i := 1; i <= 2; i++ {
		boardReq := map[string]string{
			"name":        fmt.Sprintf("Новая доска %d", i),
			"description": "Описание новой доски",
			"background":  "#FF5733",
		}
		reqBody, _ := json.Marshal(boardReq)
		resp, err := client.Post(boardURL, "application/json", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Доска должна быть успешно создана")

		var board Board
		json.NewDecoder(resp.Body).Decode(&board)
		boards = append(boards, board)
	}

	// 7) Удаление одной из досок
	deleteBoardID := boards[0].ID
	deleteURL := fmt.Sprintf("%s/boards/board_%d", server.URL, deleteBoardID)
	req, err = http.NewRequest("DELETE", deleteURL, nil)
	assert.NoError(t, err)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Доска должна быть успешно удалена")

	// 8) Получение списка досок
	getBoardsURL := fmt.Sprintf("%s/boards/my", server.URL)
	resp, err = client.Get(getBoardsURL)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Должен вернуться список досок")

	var myBoards []Board
	json.NewDecoder(resp.Body).Decode(&myBoards)
	assert.Len(t, myBoards, 1, "Должна остаться одна доска после удаления")
	assert.Equal(t, boards[1].ID, myBoards[0].ID, "Оставшаяся доска должна соответствовать созданной")
}
