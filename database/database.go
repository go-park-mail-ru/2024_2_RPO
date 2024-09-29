package database

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
)

func HelloDatabase() {
	fmt.Println("Hello Database")
}

var (
	db   *sql.DB
	once sync.Once
)

// Устанавливает соединение с базой данных PostgreSQL.
func ConnectToDb(port int, user string, passwd string) error {
	var err error
	once.Do(func() {
		connStr := fmt.Sprintf("host=localhost port=%d user=%s password=%s dbname=pumpkin sslmode=disable", port, user, passwd)
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			err = errors.New(fmt.Sprintf("Database connection error: %s", err.Error()))
			return
		}
		err = db.Ping()
		if err != nil {
			err = errors.New(fmt.Sprintf("Database ping error: %s", err.Error()))
		}
	})
	return err
}

// GetDbConnection возвращает установленное соединение с базой данных.
// Если соединение не установлено, возвращается ошибка.
func GetDbConnection() (*sql.DB, error) {
	if db == nil {
		return nil, fmt.Errorf("No DB connection")
	}
	if db.Ping() == nil {
		return db, ConnectToDb(5432, "tarasovxx", "my_secure_password")
	}
	return db, nil
}
