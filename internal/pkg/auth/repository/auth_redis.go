package repository

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/pkg/utils/logging"
	"RPO_back/internal/pkg/utils/pgxiface"
	"context"
	"errors"
	"fmt"
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
func (r *AuthRepository) CreateSession(ctx context.Context, sessionID string, userID int64) error {
	funcName := "CreateSession"
	redisConn := r.redisDb.Conn(r.redisDb.Context())
	defer redisConn.Close()

	err := redisConn.Set(r.redisDb.Context(), fmt.Sprintf("%s%s", sessionPrefix, sessionID), userID, sessionLifeTime).Err()
	logging.Debugf(ctx, "%s query 1 has err: %v", funcName, err)
	if err != nil {
		return fmt.Errorf("%s (session): %w", funcName, err)
	}

	err = redisConn.SAdd(ctx, fmt.Sprintf("%s%d", userPrefix, userID), sessionID).Err()
	logging.Debugf(ctx, "%s query 2 has err: %v", funcName, err)
	if err != nil {
		return fmt.Errorf("%s (user): %w", funcName, err)
	}

	return nil
}

// RemoveSession удаляет сессию из Redis
func (r *AuthRepository) RemoveSession(ctx context.Context, sessionID string) error {
	funcName := "RemoveSession"
	redisConn := r.redisDb.Conn(r.redisDb.Context())
	defer redisConn.Close()

	res := redisConn.Get(r.redisDb.Context(), sessionPrefix+sessionID)
	err := res.Err()
	logging.Debugf(ctx, "%s query 1 has err: %v", funcName, err)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("%s (get): %w", funcName, errs.ErrNotFound)
		}
		return fmt.Errorf("%s (get): %w", funcName, err)
	}

	userID, err := strconv.ParseInt(res.Val(), 10, 64)
	logging.Debugf(ctx, "%s query 2 has err: %v", funcName, err)
	if err != nil {
		return fmt.Errorf("%s (atoi): %w", funcName, err)
	}

	err = redisConn.Del(r.redisDb.Context(), sessionPrefix+sessionID).Err()
	logging.Debugf(ctx, "%s query 3 has err: %v", funcName, err)
	if err != nil {
		return fmt.Errorf("%s (del): %w", funcName, err)
	}

	userKey := fmt.Sprintf("%s%d", userPrefix, userID)
	err = redisConn.SRem(ctx, userKey, sessionPrefix+sessionID).Err()
	logging.Debugf(ctx, "%s query 4 has err: %v", funcName, err)
	if err != nil {
		return fmt.Errorf("%s (srem): %w", funcName, err)
	}

	return nil
}

// DisplaceUserSessions удаляет все сессии пользователя из Redis, кроме одной сессии - sessionID
func (r *AuthRepository) DisplaceUserSessions(ctx context.Context, sessionID string, userID int64) error {
	funcName := "DisplaceUserSessions"
	setKey := fmt.Sprintf("%s%d", userPrefix, userID)

	redisConn := r.redisDb.Conn(r.redisDb.Context())
	defer redisConn.Close()

	sessions, err := redisConn.SMembers(ctx, setKey).Result()
	logging.Debugf(ctx, "%s query 1 has err: %v", funcName, err)
	if err != nil {
		return fmt.Errorf("%s: %w", funcName, err)
	}

	sessionsToDelete := make([]interface{}, 0)
	for _, session := range sessions {
		if session != sessionID {
			sessionsToDelete = append(sessionsToDelete, session)
			res := redisConn.Del(ctx, sessionPrefix+sessionID)
			err = res.Err()
			logging.Debugf(ctx, "%s query 2 has err: %v", funcName, err)
			if err != nil {
				return fmt.Errorf("%s: %w", funcName, err)
			}
		}
	}

	if len(sessionsToDelete) != 0 {
		res := redisConn.SRem(ctx, setKey, sessionsToDelete...)
		err = res.Err()
		logging.Debugf(ctx, "%s query 3 has err: %v", funcName, err)
		if err != nil {
			return fmt.Errorf("%s: %w", funcName, err)
		}
	}

	return nil
}

// CheckSession ходит в Redis и получает UserID (или не получает и даёт ошибку errs.ErrNotFound)
func (r *AuthRepository) CheckSession(ctx context.Context, sessionID string) (userID int64, err error) {
	funcName := "CheckSession"
	redisConn := r.redisDb.Conn(r.redisDb.Context())
	defer redisConn.Close()

	val, err := redisConn.Get(r.redisDb.Context(), sessionPrefix+sessionID).Result()
	logging.Debugf(ctx, "%s query has err: %v", funcName, err)
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
