package repository

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
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
func (repo *AuthRepository) RegisterSessionRedis(ctx context.Context, cookie string, userID int) error {
	redisConn := repo.redisDb.Conn(repo.redisDb.Context())
	defer redisConn.Close()

	ttl := 7 * 24 * time.Hour

	err := redisConn.Set(repo.redisDb.Context(), cookie, userID, ttl).Err()
	logging.Debug(ctx, "RegisterSessionRedis query to redis has err: ", err)
	if err != nil {
		return fmt.Errorf("unable to set session in Redis: %v", err)
	}

	return nil
}

// KillSessionRedis удаляет сессию из Redis
func (repo *AuthRepository) KillSessionRedis(ctx context.Context, sessionID string) error {
	redisConn := repo.redisDb.Conn(repo.redisDb.Context())
	defer redisConn.Close()

	err := redisConn.Del(repo.redisDb.Context(), sessionID).Err()
	logging.Debug(ctx, "KillSessionRedis query to redis has err: ", err)
	if err != nil {
		return err
	}

	return nil
}

// RetrieveUserIdFromSessionId ходит в Redis и получает UserID (или не получает и даёт ошибку errs.ErrNotFound)
func (repo *AuthRepository) RetrieveUserIDFromSession(ctx context.Context, sessionID string) (userID int, err error) {
	redisConn := repo.redisDb.Conn(repo.redisDb.Context())
	defer redisConn.Close()

	val, err := redisConn.Get(repo.redisDb.Context(), sessionID).Result()
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

// GetUserByID получает данные пользователя из базы по id
func (repo *AuthRepository) GetUserByID(ctx context.Context, userID int) (user *models.UserProfile, err error) {
	query := `
	SELECT u_id, nickname, email, description,
	joined_at, updated_at, password_hash
	FROM "user"
	WHERE u_id=$1;`
	user = &models.UserProfile{}
	err = repo.db.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Description,
		&user.JoinedAt,
		&user.UpdatedAt,
		&user.PasswordHash,
	)
	logging.Debug(ctx, "GetUserByID query has err: ", err)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrWrongCredentials
		}
		return nil, fmt.Errorf("GetUserByID: %w", err)
	}
	return user, nil
}

// CreateUser создаёт пользователя (или не создаёт, если повторяются креды)
func (repo *AuthRepository) CreateUser(ctx context.Context, user *models.UserRegistration, hashedPassword string) (newUser *models.UserProfile, err error) {
	newUser = &models.UserProfile{}
	query := `INSERT INTO "user" (nickname, email, password_hash, description, joined_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING u_id, nickname, email, password_hash, description, joined_at, updated_at`

	err = repo.db.QueryRow(ctx, query, user.Name, user.Email, "", time.Now(), time.Now()).Scan(
		&newUser.ID,
		&newUser.Name,
		&newUser.Email,
		&newUser.Description,
		&newUser.JoinedAt,
		&newUser.UpdatedAt,
	)
	logging.Debug(ctx, "CreateUser query has err: ", err)
	return newUser, err
}

// CheckUniqueCredentials проверяет, существуют ли такие логин и email в базе
func (repo *AuthRepository) CheckUniqueCredentials(ctx context.Context, nickname string, email string) error {
	query1 := `SELECT COUNT(*) FROM "user" WHERE nickname = $1;`
	query2 := `SELECT COUNT(*) FROM "user" WHERE email = $1;`
	var count1, count2 int
	err := repo.db.QueryRow(ctx, query1, nickname).Scan(&count1)
	logging.Debug(ctx, "CheckUniqueCredentials query 1 has err: ", err)
	if err != nil {
		return fmt.Errorf("AuthRepository CheckUniqueCredentials (query1): %w", err)
	}
	err = repo.db.QueryRow(ctx, query2, email).Scan(&count2)
	logging.Debug(ctx, "CheckUniqueCredentials query 2 has err: ", err)
	if err != nil {
		return fmt.Errorf("AuthRepository CheckUniqueCredentials (query2): %w", err)
	}
	if count1 > 0 && count2 > 0 {
		return fmt.Errorf("AuthRepository CheckUniqueCredentials: %w %w", errs.ErrBusyNickname, errs.ErrBusyEmail)
	} else if count1 > 0 {
		return fmt.Errorf("AuthRepository CheckUniqueCredentials: %w", errs.ErrBusyNickname)
	} else if count1 > 0 {
		return fmt.Errorf("AuthRepository CheckUniqueCredentials: %w", errs.ErrBusyEmail)
	}
	return nil
}

// SetNewPasswordHash устанавливает пользователю новый хеш пароля
func (repo *AuthRepository) SetNewPasswordHash(ctx context.Context, userID int, newPasswordHash string) error {
	query := `
	UPDATE "user"
	SET password_hash=$1
	WHERE u_id=$2;
	`
	tag, err := repo.db.Exec(ctx, query, newPasswordHash, userID)
	logging.Debug(ctx, "SetNewPasswordHash query has err: ", err, " tag: ", tag)
	if err != nil {
		return fmt.Errorf("SetNewPasswordHash: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("SetNewPasswordHash: No password change done")
	}
	return nil
}
