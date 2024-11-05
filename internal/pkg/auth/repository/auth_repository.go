package repository

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v5"
)

type AuthRepository struct {
	db   pgxtype.Querier
	redisDb *redis.Client
}

func CreateAuthRepository(postgresDb pgxtype.Querier, redisDb *redis.Client) *AuthRepository {
	return &AuthRepository{
		db: postgresDb, redisDb: redisDb,
	}
}

// Регистрирует сессионную куку в Redis
func (repo *AuthRepository) RegisterSessionRedis(cookie string, userID int) error {
	redisConn := repo.redisDb.Conn(repo.redisDb.Context())
	defer redisConn.Close()

	ttl := 7 * 24 * time.Hour

	err := redisConn.Set(repo.redisDb.Context(), cookie, userID, ttl).Err()
	if err != nil {
		return fmt.Errorf("unable to set session in Redis: %v", err)
	}

	return nil
}

// KillSessionRedis удаляет сессию из Redis
func (repo *AuthRepository) KillSessionRedis(sessionId string) error {
	redisConn := repo.redisDb.Conn(repo.redisDb.Context())
	defer redisConn.Close()

	if err := redisConn.Del(repo.redisDb.Context(), sessionId).Err(); err != nil {
		return err
	}

	return nil
}

// RetrieveUserIdFromSessionId ходит в Redis и получает UserID (или не получает и даёт ошибку errs.ErrNotFound)
func (repo *AuthRepository) RetrieveUserIdFromSessionId(sessionId string) (userId int, err error) {
	redisConn := repo.redisDb.Conn(repo.redisDb.Context())
	defer redisConn.Close()

	val, err := redisConn.Get(repo.redisDb.Context(), sessionId).Result()
	if err == redis.Nil {
		return 0, fmt.Errorf("RetrieveUserIdFromSessionId(%v): %w", sessionId, errs.ErrNotFound)
	} else if err != nil {
		return 0, err
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("error converting value to int: %v", err)
	}

	return intVal, nil
}

// GetUserByEmail получает данные пользователя из базы по email
func (repo *AuthRepository) GetUserByEmail(email string) (user *models.UserProfile, err error) {
	query := `
	SELECT u_id, nickname, email, description,
	joined_at, updated_at, password_hash
	FROM "user"
	WHERE email=$1;`
	user = &models.UserProfile{}
	err = repo.db.QueryRow(context.Background(), query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Description,
		&user.JoinedAt,
		&user.UpdatedAt,
		&user.PasswordHash,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrWrongCredentials
		}
		return nil, err
	}
	return user, nil
}

// GetUserByID получает данные пользователя из базы по id
func (repo *AuthRepository) GetUserByID(userID int) (user *models.UserProfile, err error) {
	query := `
	SELECT u_id, nickname, email, description,
	joined_at, updated_at, password_hash
	FROM "user"
	WHERE u_id=$1;`
	user = &models.UserProfile{}
	err = repo.db.QueryRow(context.Background(), query, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Description,
		&user.JoinedAt,
		&user.UpdatedAt,
		&user.PasswordHash,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrWrongCredentials
		}
		return nil, fmt.Errorf("GetUserByID: %w", err)
	}
	return user, nil
}

// CreateUser создаёт пользователя (или не создаёт, если повторяются креды)
func (repo *AuthRepository) CreateUser(user *models.UserRegistration, hashedPassword string) (newUser *models.UserProfile, err error) {
	newUser = &models.UserProfile{}
	query := `INSERT INTO "user" (nickname, email, password_hash, description, joined_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING u_id, nickname, email, password_hash, description, joined_at, updated_at`

	err = repo.db.QueryRow(context.Background(), query, user.Name, user.Email, hashedPassword, "", time.Now(), time.Now()).Scan(
		&newUser.ID,
		&newUser.Name,
		&newUser.Email,
		&newUser.PasswordHash,
		&newUser.Description,
		&newUser.JoinedAt,
		&newUser.UpdatedAt,
	)
	return newUser, err
}

// CheckUniqueCredentials проверяет, существуют ли такие логин и email в базе
func (repo *AuthRepository) CheckUniqueCredentials(nickname string, email string) error {
	query1 := `SELECT COUNT(*) FROM "user" WHERE nickname = $1`
	query2 := `SELECT COUNT(*) FROM "user" WHERE email = $1`
	var count1, count2 int
	err := repo.db.QueryRow(context.Background(), query1, nickname).Scan(&count1)
	if err != nil {
		return fmt.Errorf("AuthRepository CheckUniqueCredentials (query1): %w", err)
	}
	err = repo.db.QueryRow(context.Background(), query2, email).Scan(&count2)
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
func (repo *AuthRepository) SetNewPasswordHash(userID int, newPasswordHash string) error {
	query := `
	UPDATE "user"
	SET password_hash=$1
	WHERE u_id=$2;
	`
	tag, err := repo.db.Exec(context.Background(), query, newPasswordHash, userID)
	if err != nil {
		return fmt.Errorf("SetNewPasswordHash: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("SetNewPasswordHash: No password change done")
	}
	return nil
}
