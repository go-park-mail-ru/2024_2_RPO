package database

import (
	"RPO_back/models"
	"context"
	"database/sql"
	"fmt"
)

func GetUserByID(userID int) (*models.User, error) {
	query := `
        SELECT u_id, nickname, email, description, joined_at, updated_at
        FROM "User"
        WHERE u_id = $1
    `
	conn, err := GetDbConnection()
	if err != nil {
		return nil, err
	}
	row := conn.QueryRow(context.Background(), query, userID)

	var user models.User
	err = row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Description,
		&user.JoinedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("User with ID %d not found", userID)
		}
		return nil, fmt.Errorf("Error while retrieving user: %w", err)
	}

	return &user, nil
}
