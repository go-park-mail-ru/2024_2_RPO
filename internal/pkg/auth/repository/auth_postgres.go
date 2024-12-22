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
func (r *AuthRepository) SetNewPasswordHash(ctx context.Context, userID int64, newPasswordHash string) error {
	funcName := "SetNewPasswordHash"
	query := `
	UPDATE "user"
	SET password_hash=$1
	WHERE u_id=$2;
	`

	tag, err := r.db.Exec(ctx, query, newPasswordHash, userID)
	logging.Debugf(ctx, "%s query has err: %v", funcName, err)
	if err != nil {
		return fmt.Errorf("%s: %w", funcName, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: No password change done", funcName)
	}
	return nil
}

// GetUserPasswordHash получает хеш пароля пользователя
func (r *AuthRepository) GetUserPasswordHash(ctx context.Context, userID int64) (passwordHash *string, err error) {
	funcName := "GetUserPasswordHash"
	query := `
	SELECT password_hash
	FROM "user"
	WHERE u_id = $1;
	`

	err = r.db.QueryRow(ctx, query, userID).Scan(&passwordHash)
	logging.Debugf(ctx, "%s query has err: %v", funcName, err)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", funcName, errs.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", funcName, err)
	}

	return passwordHash, nil
}
