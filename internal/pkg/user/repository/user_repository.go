package repository

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/utils/logging"
	"RPO_back/internal/pkg/utils/uploads"
	"context"
	"errors"
	"fmt"
	"time"

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
func (r *UserRepository) GetUserProfile(ctx context.Context, userID int) (profile *models.UserProfile, err error) {
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
	row := r.db.QueryRow(ctx, query, userID)

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
	logging.Debug(ctx, "GetUserProfile query has err: ", err)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetUserProfile: %w", errs.ErrNotFound)
		}
		return nil, fmt.Errorf("GetUserProfile: %w", err)
	}
	user.AvatarImageURL = uploads.JoinFileURL(fileUUID, fileExt, uploads.DefaultAvatarURL)

	return &user, nil
}

// UpdateUserProfile обновляет профиль пользователя
func (r *UserRepository) UpdateUserProfile(ctx context.Context, userID int, data models.UserProfileUpdateRequest) (newProfile *models.UserProfile, err error) {
	query1 := `SELECT COUNT(*) FROM "user" WHERE email=$1 AND u_id!=$2;`
	query2 := `SELECT COUNT(*) FROM "user" WHERE nickname=$1 AND u_id!=$2;`
	query3 := `
	UPDATE "user"
	SET email=$1, nickname=$2
	WHERE u_id=$3;`
	var nicknameCount, emailCount int
	row := r.db.QueryRow(ctx, query1, data.Email, userID)
	err = row.Scan(&emailCount)
	logging.Debug(ctx, "UpdateUserProfile query 1 has err: ", err)
	if err != nil {
		return nil, fmt.Errorf("UpdateUserProfile (check unique email): %w", err)
	}
	row = r.db.QueryRow(ctx, query2, data.NewName, userID)
	err = row.Scan(&nicknameCount)
	logging.Debug(ctx, "UpdateUserProfile query 2 has err: ", err)
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
	tag, err := r.db.Exec(ctx, query3, data.Email, data.NewName, userID)
	logging.Debug(ctx, "UpdateUserProfile query 3 has err: ", err)
	if err != nil {
		return nil, fmt.Errorf("UpdateUserProfile (action): %w", err)
	}
	if tag.RowsAffected() == 0 {
		return nil, fmt.Errorf("UpdateUserProfile (action): UPDATE made no changes")
	}
	newProfile, err = r.GetUserProfile(ctx, userID)
	return
}

func (r *UserRepository) SetUserAvatar(ctx context.Context, userID int, fileExtension string, fileSize int) (fileName string, err error) {
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
	row := r.db.QueryRow(ctx, query1, fileExtension, userID, fileSize)
	err = row.Scan(&fileUUID)
	logging.Debug(ctx, "SetUserAvatar query 1 has err: ", err)
	if err != nil {
		return "", fmt.Errorf("SetUserAvatar (register file): %w", err)
	}
	tag, err := r.db.Exec(ctx, query2, fileUUID, userID)
	logging.Debug(ctx, "SetUserAvatar query 2 has err: ", err)
	if err != nil {
		return "", fmt.Errorf("SetUserAvatar (update user): %w", err)
	}
	if tag.RowsAffected() == 0 {
		return "", fmt.Errorf("SetUserAvatar (update user): no rows affected")
	}
	return uploads.JoinFilePath(fileUUID, fileExtension), nil
}

// GetUserByEmail получает данные пользователя из базы по email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (user *models.UserProfile, err error) {
	query := `
	SELECT u_id, nickname, email, description,
	joined_at, updated_at
	FROM "user"
	WHERE email=$1;`
	user = &models.UserProfile{}
	err = r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Description,
		&user.JoinedAt,
		&user.UpdatedAt,
	)
	logging.Debug(ctx, "GetUserByEmail query has err: ", err)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrWrongCredentials
		}
		return nil, err
	}
	return user, nil
}

// CreateUser создаёт пользователя (или не создаёт, если повторяются креды)
func (r *UserRepository) CreateUser(ctx context.Context, user *models.UserRegisterRequest) (newUser *models.UserProfile, err error) {
	newUser = &models.UserProfile{}
	query := `INSERT INTO "user" (nickname, email, password_hash, description, joined_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING u_id, nickname, email, password_hash, description, joined_at, updated_at`

	err = r.db.QueryRow(ctx, query, user.Name, user.Email, "", time.Now(), time.Now()).Scan(
		&newUser.ID,
		&newUser.Name,
		&newUser.Email,
		&newUser.Description,
		&newUser.JoinedAt,
		&newUser.UpdatedAt,
	)
	logging.Debug(ctx, "CreateUser query has err: ", err)
	return newUser, err
}

// CheckUniqueCredentials проверяет, существуют ли такие логин и email в базе
func (r *UserRepository) CheckUniqueCredentials(ctx context.Context, nickname string, email string) error {
	funcName := `UserRepository.CheckUniqueCredentials`
	query := `SELECT nickname, email FROM "user" WHERE nickname = $1 OR email=$2;`
	var emailCount, nicknameCount int
	rows, err := r.db.Query(ctx, query, nickname)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("%s: %w", funcName, err)
	}
	for rows.Next() {
		var knownNickname, knownEmail string
		err := rows.Scan(&knownNickname, &knownEmail)
		if err != nil {
			return fmt.Errorf("%s: %w", funcName, err)
		}
		if knownEmail == email {
			emailCount++
		}
		if knownNickname == nickname {
			nicknameCount++
		}
	}
	if emailCount > 0 && nicknameCount > 0 {
		return fmt.Errorf("%s: %w %w", funcName, errs.ErrBusyNickname, errs.ErrBusyEmail)
	} else if nicknameCount > 0 {
		return fmt.Errorf("%s: %w", funcName, errs.ErrBusyNickname)
	} else if emailCount > 0 {
		return fmt.Errorf("%s: %w", funcName, errs.ErrBusyEmail)
	}
	return nil
}
