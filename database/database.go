package database

import (
	"RPO_back/auth"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"sync"

	_ "github.com/lib/pq"
)

var (
	port   int
	user   string
	passwd string
	db     *sql.DB
	mu     sync.Mutex
)

func GetUserId(w http.ResponseWriter, r *http.Request) (int, error) {
	sessionCookie, err := r.Cookie("session_id")
	if err != nil || sessionCookie.Value == "" {
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return 0, errors.New("No session cookie detected")
	}

	userId, err2 := auth.RetrieveUserIdFromSessionId(sessionCookie.Value)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusForbidden)
	}

	return userId, nil
}

func InitDBConnection(port_ int, user_ string, passwd_ string) error {
	port = port_
	user = user_
	passwd = passwd_
	err := ConnectToDb()
	if err == nil {
		fmt.Println("Successfully connected to Postgres!")
	}
	return err
}

// Устанавливает соединение с базой данных PostgreSQL.
func ConnectToDb() error {
	var err error
	connStr := fmt.Sprintf("host=localhost port=%d user=%s password=%s dbname=pumpkin sslmode=disable", port, user, passwd)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return errors.New(fmt.Sprintf("Database connection error: %s", err.Error()))
	}
	err = db.Ping()
	if err != nil {
		return errors.New(fmt.Sprintf("Database ping error: %s", err.Error()))
	}

	return err
}

// GetDbConnection возвращает установленное соединение с базой данных.
// Если соединение не установлено, происходит попытка переподключения. Если попытка неудачная, возвращается ошибка.
func GetDbConnection() (*sql.DB, error) {
	mu.Lock()
	defer mu.Unlock()

	// Если соединение еще не установлено, устанавливаем его
	if db == nil {
		var err error
		err = ConnectToDb()
		if err != nil {
			return nil, fmt.Errorf("не удалось установить соединение с БД: %w", err)
		}
		return db, nil
	}

	// Проверяем текущее состояние соединения
	if err := db.Ping(); err != nil {
		// Если соединение потеряно, пытаемся восстановить его
		db.Close() // Закрываем старое соединение
		var errReconnect error
		errReconnect = ConnectToDb()
		if errReconnect != nil {
			return nil, fmt.Errorf("соединение с БД закрыто и восстановить его не удалось: %w", errReconnect)
		}
	}

	return db, nil
}
