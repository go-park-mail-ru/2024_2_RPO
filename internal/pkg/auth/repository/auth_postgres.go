package repository

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/pkg/utils/logging"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

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

// GetUserPasswordHash получает хеш пароля пользователя
func (r *AuthRepository) GetUserPasswordHash(ctx context.Context, userID int) (passwordHash *string, err error) {
	query := `
	SELECT password_hash
	FROM "user"
	WHERE u_id = $1;
	`

	err = r.db.QueryRow(ctx, query, userID).Scan(&passwordHash)
	logging.Debug(ctx, "GetUserPasswordHash query has err: ", err)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetUserPasswordHash: %w", errs.ErrNotFound)
		}
		return nil, fmt.Errorf("GetUserPasswordHash: %w", err)
	}

	return passwordHash, nil
}
