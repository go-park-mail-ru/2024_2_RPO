package repository

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func CreateUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// GetUserProfile возвращает профиль пользователя
func (r *UserRepository) GetUserProfile(userID int) (profile *models.UserProfile, err error) {
	query := `
        SELECT u.u_id, u.nickname, u.email, u.description, u.joined_at, u.updated_at
        FROM "user" AS u
		LEFT JOIN user_uploaded_file AS f ON u.avatar_file_uuid=f.file_uuid
        WHERE u_id = $1
    `
	row := r.db.QueryRow(context.Background(), query, userID)

	var user models.UserProfile
	err = row.Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Description,
		&user.JoinedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetUserById: %w", errs.ErrNotFound)
		}
		return nil, fmt.Errorf("GetUserById: %w", err)
	}

	return &user, nil
}

// UpdateUserProfile обновляет профиль пользователя
func (r *UserRepository) UpdateUserProfile(userID int, data models.UserProfileUpdate) (newProfile *models.UserProfile, err error) {
	panic("Not implemented")
}

// TODO загрузка аватарки
