package repository

import (
	"RPO_back/auth"
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	url 	string
	db     *pgxpool.Pool
	mu     sync.Mutex
)

func GetUserId(w http.ResponseWriter, r *http.Request) (int, error) {
	sessionCookie, err := r.Cookie("session_id")
	if err != nil || sessionCookie.Value == "" {
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return 0, errors.New("no session cookie detected")
	}

	userId, err2 := auth.RetrieveUserIdFromSessionId(sessionCookie.Value)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusForbidden)
	}

	return userId, nil
}

func InitDBConnection(url_ string) error {
	url = url_
	err := ConnectToDb()
	if err == nil {
		fmt.Println("Successfully connected to Postgres!")
	}
	return err
}

// Устанавливает соединение с базой данных PostgreSQL.
func ConnectToDb() error {
	var err error
	db, err = pgxpool.New(context.Background(), url)
	if err != nil {
		return fmt.Errorf("database connection error: %s", err)
	}

	conn, err := db.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("database ping error: %s", err)
	}
	defer conn.Release()

	return err
}

// GetDbConnection возвращает установленное соединение с базой данных.
// Если соединение не установлено, происходит попытка переподключения. Если попытка неудачная, возвращается ошибка.
func GetDbConnection() (*pgxpool.Pool, error) {
	mu.Lock()
	defer mu.Unlock()

	// Если соединение еще не установлено, устанавливаем его
	if db == nil {
		err := ConnectToDb()
		if err != nil {
			return nil, fmt.Errorf("не удалось установить соединение с БД: %w", err)
		}
		return db, nil
	}

	// Проверяем текущее состояние соединения
	if _, err := db.Acquire(context.Background()); err != nil {
		// Если соединение потеряно, пытаемся восстановить его
		db.Close() // Закрываем старое соединение
		errReconnect := ConnectToDb()
		if errReconnect != nil {
			return nil, fmt.Errorf("соединение с БД закрыто и восстановить его не удалось: %w", errReconnect)
		}
	}

	return db, nil
}
