package repository

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/utils/errs"
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (this *UserRepository) GetUserById(userID int) (*models.User, error) {
	query := `
        SELECT u.u_id, u.nickname, u.email, u.description, u.joined_at, u.updated_at
        FROM "user" AS u
		LEFT JOIN user_uploaded_file AS f ON u.avatar_file_uuid=f.file_uuid
        WHERE u_id = $1
    `
	row := this.db.QueryRow(context.Background(), query, userID)

	var user models.User
	err := row.Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Description,
		&user.JoinedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("GetUserById: %w", errs.ErrNotFound)
		}
		return nil, fmt.Errorf("GetUserById: %w", err)
	}

	return &user, nil
}
