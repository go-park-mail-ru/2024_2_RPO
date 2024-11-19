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
	"github.com/jackc/pgx/v5"
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

// Регистрирует сессионную куку в Redis
func (r *AuthRepository) RegisterSessionRedis(ctx context.Context, sessionID string, userID int) error {
	panic("TODO сделать добавление сессии в массив индекса")
	redisConn := r.redisDb.Conn(r.redisDb.Context())
	defer redisConn.Close()

	ttl := 7 * 24 * time.Hour

	err := redisConn.Set(r.redisDb.Context(), sessionID, userID, ttl).Err()
	logging.Debug(ctx, "RegisterSessionRedis query to redis has err: ", err)
	if err != nil {
		return fmt.Errorf("unable to set session in Redis: %v", err)
	}

	return nil
}

// KillSessionRedis удаляет сессию из Redis
func (r *AuthRepository) KillSessionRedis(ctx context.Context, sessionID string) error {
	panic("TODO сделать удаление сессии из индексного массива")
	redisConn := r.redisDb.Conn(r.redisDb.Context())
	defer redisConn.Close()

	err := redisConn.Del(r.redisDb.Context(), sessionID).Err()
	logging.Debug(ctx, "KillSessionRedis query to redis has err: ", err)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) DisplaceUserSessions(ctx context.Context, sessionID string, userID int64) error {
	panic("TODO реализовать вытеснение")
	var keysToDelete []string

	cursor := uint64(0)
	for {
		keys, newCursor, err := r.redisDb.Scan(ctx, cursor, "*", 100).Result()
		if err != nil {
			return fmt.Errorf("error while scanning keys: %w", err)
		}

		if len(keys) > 0 {
			pipe := r.redisDb.Pipeline()
			cmds := make(map[string]*redis.StringCmd)

			for _, key := range keys {
				cmds[key] = pipe.Get(ctx, key)
			}

			_, err = pipe.Exec(ctx)
			if err != nil && err != redis.Nil {
				return fmt.Errorf("error while executing pipeline: %w", err)
			}

			for key, cmd := range cmds {
				val, err := cmd.Result()
				if err == nil && val == sessionID {
					userIDFromRedisCmd := pipe.Get(ctx, fmt.Sprintf("%d", userID))
					userIDFromRedis, err := userIDFromRedisCmd.Result()
					if err == nil {
						parseUserIDFromRedis, err := strconv.Atoi(userIDFromRedis)
						if err == nil && parseUserIDFromRedis != int(userID) {
							keysToDelete = append(keysToDelete, key)
						}
					}
				}
			}
		}

		cursor = newCursor
		if cursor == 0 {
			break
		}
	}

	if len(keysToDelete) > 0 {
		count, err := r.redisDb.Del(ctx, keysToDelete...).Result()
		if err != nil {
			return fmt.Errorf("error while deleting keys: %w", err)
		}
		logging.Info(ctx, "Deleted %d keys", count)
	} else {
		logging.Info(ctx, "Keys with specified value not found.")
	}

	return nil
}

// CheckSession ходит в Redis и получает UserID (или не получает и даёт ошибку errs.ErrNotFound)
func (r *AuthRepository) CheckSession(ctx context.Context, sessionID string) (userID int, err error) {
	redisConn := r.redisDb.Conn(r.redisDb.Context())
	defer redisConn.Close()

	val, err := redisConn.Get(r.redisDb.Context(), sessionID).Result()
	logging.Debug(ctx, "RetrieveUserIDFromSession query to redis has err: ", err)
	if err == redis.Nil {
		return 0, fmt.Errorf("RetrieveUserIDFromSession(%v): %w", sessionID, errs.ErrNotFound)
	} else if err != nil {
		return 0, err
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("error converting value to int: %v", err)
	}

	return intVal, nil
}

// SetNewPasswordHash устанавливает пользователю новый хеш пароля
func (r *AuthRepository) SetNewPasswordHash(ctx context.Context, userID int, newPasswordHash string) error {
	query := `
	UPDATE "user"
	SET password_hash=$1
	WHERE u_id=$2;
	`
	tag, err := r.db.Exec(ctx, query, newPasswordHash, userID)
	logging.Debug(ctx, "SetNewPasswordHash query has err: ", err, " tag: ", tag)
	if err != nil {
		return fmt.Errorf("SetNewPasswordHash: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("SetNewPasswordHash: No password change done")
	}
	return nil
}

func (r *AuthRepository) GetUserPasswordHashForUser(ctx context.Context, userID int) (passwordHash string, err error) {
	query := `
	SELECT password_hash
	FROM "user"
	WHERE u_id = $1
	`

	err = r.db.QueryRow(ctx, query, userID).Scan(&passwordHash)
	logging.Debug(ctx, "GetUserPasswordHashForUser query has err: ", err)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errs.ErrNotFound
		}
		return "", fmt.Errorf("GetUserPasswordHashForUser: %w", err)
	}

	return passwordHash, nil
}
