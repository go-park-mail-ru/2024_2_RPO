package repository

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/auth"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type contextKey string

const (
	UserIDContextKey contextKey = "userId"
)

type AuthRepository struct {
	db      *pgxpool.Pool
	redisDb *redis.Client
}

func CreateAuthRepository(postgresDb *pgxpool.Pool, redisDb *redis.Client) *AuthRepository {
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

// Удаляет сессию из Redis
func (repo *AuthRepository) KillSessionRedis(sessionId string) error {
	redisConn := repo.redisDb.Conn(repo.redisDb.Context())
	defer redisConn.Close()
	if err := redisConn.Del(repo.redisDb.Context(), sessionId).Err(); err != nil {
		return err
	}
	return nil
}

func (repo *AuthRepository) RetrieveUserIdFromSessionId(sessionId string) (userId int, err error) {
	redisConn := repo.redisDb.Conn(repo.redisDb.Context())
	defer redisConn.Close()

	val, err := redisConn.Get(repo.redisDb.Context(), sessionId).Result()
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
	userID, ok := ctx.Value(UserIDContextKey).(string)
	return userID, ok
}

func (repo *AuthRepository) GetUserByEmail(email string) (user *models.User, err error) {
	user = &models.User{}
	selectError := repo.db.QueryRow(context.Background(), "SELECT u_id, nickname, email, description, joined_at, updated_at, password_hash FROM \"User\" WHERE email=$1", email).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Description,
		&user.JoinedAt,
		&user.UpdatedAt,
		&user.PasswordHash,
	)
	if selectError != nil {
		if errors.Is(selectError, pgx.ErrNoRows) {
			return nil, auth.ErrWrongCredentials

		}
		return nil, selectError
	}
	return user, nil
}

func (repo *AuthRepository) CreateUser(user *models.UserRegistration, hashedPassword string) (newUser *models.User, err error) {
	newUser = &models.User{}
	query := `INSERT INTO "user" (nickname, email, password_hash, description, joined_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING u_id, nickname, email, password_hash, description, joined_at, updated_at`

	err = repo.db.QueryRow(context.Background(), query, user.Name, user.Email, hashedPassword, "", time.Now(), time.Now()).Scan(
		&newUser.Id,
		&newUser.Name,
		&newUser.Email,
		&newUser.PasswordHash,
		&newUser.Description,
		&newUser.JoinedAt,
		&newUser.UpdatedAt,
	)
	return newUser, err
}

func (repo *AuthRepository) CheckUniqueCredentials(nickname string, email string) error {
	query1 := `SELECT COUNT(*) FROM "user" WHERE nickname = $1`
	query2 := `SELECT COUNT(*) FROM "user" WHERE email = $1`
	var count int
	err := repo.db.QueryRow(context.Background(), query1, nickname).Scan(&count)
	if err != nil {
		return fmt.Errorf("AuthRepository CheckUniqueCredentials: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("AuthRepository CheckUniqueCredentials: %w", auth.ErrBusyNickname)
	}
	err = repo.db.QueryRow(context.Background(), query2, email).Scan(&count)
	if err != nil {
		return fmt.Errorf("AuthRepository CheckUniqueCredentials: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("AuthRepository CheckUniqueCredentials: %w", auth.ErrBusyEmail)
	}
	return nil
}