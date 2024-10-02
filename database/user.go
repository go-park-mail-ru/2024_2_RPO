package database

import (
	"RPO_back/models"
	"database/sql"
	"fmt"
)

func GetUserByID(userID int) (*models.User, error) {
	query := `
        SELECT u_id, nickname, email, description, joined_at, updated_at
        FROM "User"
        WHERE u_id = $1
    `
	conn, err0 := GetDbConnection()
	if err0 != nil {
		return nil, err0
	}
	row := conn.QueryRow(query, userID)

	var user models.User
	err := row.Scan(
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
