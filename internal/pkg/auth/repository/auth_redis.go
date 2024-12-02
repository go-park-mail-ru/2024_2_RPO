package repository

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/pkg/utils/logging"
	"RPO_back/internal/pkg/utils/pgxiface"
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	sessionPrefix   = "s_" // Префикс для сессии
	userPrefix      = "u_" // Префикс для сета, в котором находятся все сессии данного пользователя
	sessionLifeTime = 7 * 24 * time.Hour
)

type AuthRepository struct {
	db      pgxiface.PgxIface
	redisDb *redis.Client
}

func CreateAuthRepository(postgresDb pgxiface.PgxIface, redisDb *redis.Client) *AuthRepository {
	return &AuthRepository{
		db: postgresDb, redisDb: redisDb,
	}
}

// RegisterSessionRedis регистрирует сессию в Redis
func (r *AuthRepository) RegisterSessionRedis(ctx context.Context, sessionID string, userID int64) error {
	redisConn := r.redisDb.Conn(r.redisDb.Context())
	defer redisConn.Close()

	fmt.Print("REGISTER SESSION user ", userID, "   session ", sessionID)

	err := redisConn.Set(r.redisDb.Context(), fmt.Sprintf("%s%s", sessionPrefix, sessionID), userID, sessionLifeTime).Err()
	logging.Debug(ctx, "RegisterSessionRedis query to redis has err: ", err)
	if err != nil {
		return fmt.Errorf("RegisterSessionRedis (session): %w", err)
	}

	err = redisConn.SAdd(ctx, fmt.Sprintf("%s%d", userPrefix, userID), sessionID).Err()
	if err != nil {
		return fmt.Errorf("RegisterSessionRedis (user): %w", err)
	}

	return nil
}

// KillSessionRedis удаляет сессию из Redis
func (r *AuthRepository) KillSessionRedis(ctx context.Context, sessionID string) error {
	redisConn := r.redisDb.Conn(r.redisDb.Context())
	defer redisConn.Close()

	res := redisConn.Get(r.redisDb.Context(), sessionPrefix+sessionID)
	if res.Err() != nil {
		if res.Err() == redis.Nil {
			return fmt.Errorf("KillSessionRedis (get): %w", errs.ErrNotFound)
		}
		return fmt.Errorf("KillSessionRedis (get): %w", res.Err())
	}

	userID, err := strconv.ParseInt(res.Val(), 10, 64)
	if err != nil {
		return fmt.Errorf("KillSessionRedis (atoi): %w", res.Err())
	}

	err = redisConn.Del(r.redisDb.Context(), sessionPrefix+sessionID).Err()
	if err != nil {
		return fmt.Errorf("KillSessionRedis (del): %w", err)
	}

	userKey := fmt.Sprintf("%s%d", userPrefix, userID)
	res2 := redisConn.SRem(ctx, userKey, sessionPrefix+sessionID)
	if res2.Err() != nil {
		return fmt.Errorf("KillSessionRedis (srem): %w", res2.Err())
	}

	return nil
}

// DisplaceUserSessions удаляет все сессии пользователя из Redis, кроме одной сессии - sessionID
func (r *AuthRepository) DisplaceUserSessions(ctx context.Context, sessionID string, userID int64) error {
	setKey := fmt.Sprintf("%s%d", userPrefix, userID)

	redisConn := r.redisDb.Conn(r.redisDb.Context())
	defer redisConn.Close()

	sessions, err := redisConn.SMembers(ctx, setKey).Result()
	if err != nil {
		log.Fatalf("DisplaceUserSessions (get user): %v", err)
	}

	sessionsToDelete := make([]interface{}, 0)
	for _, session := range sessions {
		if session != sessionID {
			sessionsToDelete = append(sessionsToDelete, session)
			res := redisConn.Del(ctx, sessionPrefix+sessionID)
			if res.Err() != nil {
				log.Fatalf("DisplaceUserSessions (del): %v", res.Err())
			}
		}
	}

	res := redisConn.SRem(ctx, setKey, sessionsToDelete...)
	if res.Err() != nil {
		log.Fatalf("DisplaceUserSessions (srem): %v", res.Err())
	}

	return nil
}

// CheckSession ходит в Redis и получает UserID (или не получает и даёт ошибку errs.ErrNotFound)
func (r *AuthRepository) CheckSession(ctx context.Context, sessionID string) (userID int64, err error) {
	redisConn := r.redisDb.Conn(r.redisDb.Context())
	defer redisConn.Close()

	val, err := redisConn.Get(r.redisDb.Context(), sessionPrefix+sessionID).Result()
	fmt.Printf("VALUE: %v\n", val)
	logging.Debug(ctx, "CheckSession query to redis has err: ", err)
	if err == redis.Nil {
		return 0, fmt.Errorf("CheckSession (get): %w", errs.ErrNotFound)
	} else if err != nil {
		return 0, err
	}

	intVal, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("CheckSession (atoi): %v", err)
	}

	return intVal, nil
}
