package misc

import (
	"RPO_back/internal/pkg/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectToPgx подключается к PostgreSQL и делает ping
func ConnectToPgx(maxConns int) (db *pgxpool.Pool, err error) {
	// Конфиг pgx
	config, err := pgxpool.ParseConfig(config.CurrentConfig.PostgresDSN)
	if err != nil {
		return nil, fmt.Errorf("Error creating pgx config: %w", err)
	}
	config.MaxConns = int32(maxConns)

	// Подключение к PostgreSQL
	db, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("error creating PostgreSQL connection: %w", err)
	}

	// Проверка подключения к PostgreSQL
	if err = db.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("error while pinging PostgreSQL: %w", err)

	}
	return db, nil
}
