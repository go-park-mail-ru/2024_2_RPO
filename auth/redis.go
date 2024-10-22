package auth

import (
	"context"
	"fmt"
	"net/url"
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
func ConnectToRedis(url_ string) error {
	var err error
	once.Do(func() {
		url, err_ := url.Parse(url_)
		if err_ != nil {
			err = fmt.Errorf("invalid Redis URL: %s", err.Error())
			return
		}

		passwd, _ := url.User.Password()  
		host := url.Hostname()
		port := url.Port()
		rdb = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", host, port),
			Password: passwd,
			DB:       0, // Выбор стандартной БД
		})

		_, err = rdb.Ping(ctx).Result()
	})

	if err != nil {
		err = fmt.Errorf("redis connection error: %s", err.Error())
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
