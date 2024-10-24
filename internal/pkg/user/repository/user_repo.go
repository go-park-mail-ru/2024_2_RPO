package repository

import (
	"RPO_back/internal/models"
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	postgresDb *pgxpool.Pool
}

func (this *UserRepository) GetUserByID(userID int) (*models.User, error) {
	query := `
        SELECT u_id, nickname, email, description, joined_at, updated_at
        FROM "User"
        WHERE u_id = $1
    `
	row := this.postgresDb.QueryRow(context.Background(), query, userID)

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
			return nil, fmt.Errorf("user with ID %d not found", userID)
		}
		return nil, fmt.Errorf("error while retrieving user: %w", err)
	}

	return &user, nil
}
