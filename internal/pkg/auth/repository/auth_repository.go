package repository

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	userIDContextKey    string = "userId"
	sessionIdCookieName string = "session_id"
)

type AuthRepository struct {
	postgresDb *pgxpool.Pool
	redisDb    *redis.Client
}

func CreateAuthRepository(postgresDb *pgxpool.Pool, redisDb *redis.Client) *AuthRepository {
	return &AuthRepository{
		postgresDb: postgresDb, redisDb: redisDb,
	}
}

// Регистрирует сессионную куку в Redis
func (this *AuthRepository) RegisterSessionRedis(cookie string, userID int) error {
	redisConn := this.redisDb.Conn(this.redisDb.Context())
	defer redisConn.Close()
	ttl := 7 * 24 * time.Hour
	err := redisConn.Set(this.redisDb.Context(), cookie, userID, ttl).Err()
	if err != nil {
		return fmt.Errorf("unable to set session in Redis: %v", err)
	}

	return nil
}

func (this *AuthRepository) retrieveUserIdFromSessionId(sessionId string) (userId int, err error) {
	redisConn := this.redisDb.Conn(this.redisDb.Context())
	defer redisConn.Close()

	val, err := redisConn.Get(this.redisDb.Context(), sessionId).Result()
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

// Function to retrieve user ID from request context
func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDContextKey).(string)
	return userID, ok
}

func (this *AuthRepository) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		userID, err := this.retrieveUserIdFromSessionId(cookie.Value)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:   "session",
				MaxAge: -1,
			})
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), userIDContextKey, userID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
