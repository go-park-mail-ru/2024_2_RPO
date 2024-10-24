package repository

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	url string
	db  *pgxpool.Pool
	mu  sync.Mutex
)

// Регистрирует сессионную куку в Redis
func RegisterSessionRedis(cookie string, userID int) error {
	ttl := 7 * 24 * time.Hour
	err := GetRedisConnection().Set(ctx, cookie, userID, ttl).Err()
	if err != nil {
		return fmt.Errorf("unable to set session in Redis: %v", err)
	}

	return nil
}

func RetrieveUserIdFromSessionId(sessionId string) (userId int, err error) {
	val, err := rdb.Get(ctx, sessionId).Result()
	if err == redis.Nil {
		return 0, fmt.Errorf("session cookie is invalid or expired: %s", sessionId)
	} else if err != nil {
		return 0, err
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("error converting value to int: %v", err)
	}

	return intVal, nil
}
