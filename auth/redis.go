package auth

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	rdb  *redis.Client
	once sync.Once
	ctx  = context.Background()
)

// Устанавливает соединение с Redis и сохраняет клиент в глобальную переменную.
func ConnectToRedis(port int, _ string, passwd string) error {
	var err error
	once.Do(func() {
		rdb = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("127.0.0.1:%d", port),
			Password: passwd,
			DB:       0, // Выбор стандартной БД
		})

		_, err = rdb.Ping(ctx).Result()
	})
	if err != nil {
		err = errors.New(fmt.Sprintf("Redis connection error: %s", err.Error()))
	}
	return err
}

// Возвращает подключение к Redis
func GetRedisConnection() *redis.Client {
	return rdb
}

// Регистрирует сессионную куку в Redis
func RegisterSessionRedis(cookie string, userID int) error {
	ttl := 7 * 24 * time.Hour
	err := GetRedisConnection().Set(ctx, cookie, userID, ttl).Err()
	if err != nil {
		return fmt.Errorf("Unable to set session in Redis: %v", err)
	}

	return nil
}

func RetrieveUserIdFromSessionId(sessionId string) (userId int, err error) {
	val, err := rdb.Get(ctx, sessionId).Result()
	if err == redis.Nil {
		return 0, fmt.Errorf("Session cookie is invalid or expired: %s", sessionId)
	} else if err != nil {
		return 0, err
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("Error converting value to int: %v", err)
	}

	return intVal, nil
}
