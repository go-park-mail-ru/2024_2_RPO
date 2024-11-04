package repository

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/utils/uploads"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultUserAvatar string = "/static/img/KarlMarks.jpg"
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
        SELECT
		u.u_id,
		u.nickname,
		u.email,
		u.description,
		u.joined_at,
		u.updated_at,
		COALESCE(f.file_uuid::text, ''),
		COALESCE(f.file_extension, '')
        FROM "user" AS u
		LEFT JOIN user_uploaded_file AS f ON u.avatar_file_uuid=f.file_uuid
        WHERE u_id = $1;
    `
	row := r.db.QueryRow(context.Background(), query, userID)

	var user models.UserProfile
	var fileUUID, fileExt string
	err = row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Description,
		&user.JoinedAt,
		&user.UpdatedAt,
		&fileUUID,
		&fileExt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetUserById: %w", errs.ErrNotFound)
		}
		return nil, fmt.Errorf("GetUserById: %w", err)
	}
	user.AvatarImageURL = uploads.JoinFileURL(fileUUID, fileExt, defaultUserAvatar)

	return &user, nil
}

// UpdateUserProfile обновляет профиль пользователя
func (r *UserRepository) UpdateUserProfile(userID int, data models.UserProfileUpdate) (newProfile *models.UserProfile, err error) {
	query1 := `SELECT COUNT(*) FROM "user" WHERE email=$1 AND u_id!=$2;`
	query2 := `SELECT COUNT(*) FROM "user" WHERE nickname=$1 AND u_id!=$2;`
	query3 := `
	UPDATE "user"
	SET email=$1, nickname=$2
	WHERE u_id=$3;`
	var nicknameCount, emailCount int
	row := r.db.QueryRow(context.Background(), query1, data.Email, userID)
	err = row.Scan(&emailCount)
	if err != nil {
		return nil, fmt.Errorf("UpdateUserProfile (check unique email): %w", err)
	}
	row = r.db.QueryRow(context.Background(), query2, data.NewName, userID)
	err = row.Scan(&nicknameCount)
	if err != nil {
		return nil, fmt.Errorf("UpdateUserProfile (check unique nick): %w", err)
	}
	if nicknameCount != 0 && emailCount != 0 {
		return nil, fmt.Errorf("UpdateUserProfile (check unique): %w %w", errs.ErrBusyEmail, errs.ErrBusyNickname)
	}
	if nicknameCount != 0 {
		return nil, fmt.Errorf("UpdateUserProfile (check unique): %w", errs.ErrBusyNickname)
	}
	if emailCount != 0 {
		return nil, fmt.Errorf("UpdateUserProfile (check unique): %w", errs.ErrBusyEmail)
	}
	tag, err := r.db.Exec(context.Background(), query3, data.Email, data.NewName, userID)
	if err != nil {
		return nil, fmt.Errorf("UpdateUserProfile (action): %w", err)
	}
	if tag.RowsAffected() == 0 {
		return nil, fmt.Errorf("UpdateUserProfile (action): UPDATE made no changes")
	}
	newProfile, err = r.GetUserProfile(userID)
	return
}

func (r *UserRepository) SetUserAvatar(userID int, fileExtension string, fileSize int) (fileName string, err error) {
	query1 := `
	INSERT INTO user_uploaded_file
	(file_extension, created_at, created_by, "size")
	VALUES ($1, CURRENT_TIMESTAMP, $2, $3)
	RETURNING file_uuid::text;
	`
	query2 := `
	UPDATE "user"
	SET avatar_file_uuid=to_uuid($1)
	WHERE u_id=$2;`
	var fileUUID string
	row := r.db.QueryRow(context.Background(), query1, fileExtension, userID, fileSize)
	err = row.Scan(&fileUUID)
	if err != nil {
		return "", fmt.Errorf("SetUserAvatar (register file): %w", err)
	}
	tag, err := r.db.Exec(context.Background(), query2, fileUUID, userID)
	if err != nil {
		return "", fmt.Errorf("SetUserAvatar (update user): %w", err)
	}
	if tag.RowsAffected() == 0 {
		return "", fmt.Errorf("SetUserAvatar (update user): no rows affected")
	}
	return uploads.JoinFilePath(fileUUID, fileExtension), nil
}
