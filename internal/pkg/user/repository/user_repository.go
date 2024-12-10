package repository

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/utils/logging"
	"RPO_back/internal/pkg/utils/pgxiface"
	"RPO_back/internal/pkg/utils/uploads"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db pgxiface.PgxIface
}

func CreateUserRepository(db pgxiface.PgxIface) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// GetUserProfile возвращает профиль пользователя
func (r *UserRepository) GetUserProfile(ctx context.Context, userID int64) (profile *models.UserProfile, err error) {
	query := `
        SELECT
		u.u_id,
		u.nickname,
		u.email,
		u.joined_at,
		u.updated_at,
		COALESCE(f.file_uuid::text, ''),
		COALESCE(f.file_extension, '')
        FROM "user" AS u
		LEFT JOIN user_uploaded_file AS f ON u.avatar_file_id=f.file_id
        WHERE u_id = $1;
    `
	row := r.db.QueryRow(ctx, query, userID)

	var user models.UserProfile
	var fileUUID, fileExt string
	err = row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
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
func (r *UserRepository) UpdateUserProfile(ctx context.Context, userID int64, data models.UserProfileUpdateRequest) (newProfile *models.UserProfile, err error) {
	funcName := "UpdateUserProfile"
	query1 := `SELECT COUNT(*) FROM "user" WHERE (email=$1 OR nickname=$2) AND u_id!=$3;`
	query2 := `
	UPDATE "user" u
	SET email = $1, nickname = $2
	WHERE u.u_id = $3
	RETURNING
	    u.u_id,
	    u.nickname,
	    u.email,
	    u.joined_at,
	    u.updated_at,
	    COALESCE(
	        (SELECT f.file_uuid::text FROM user_uploaded_file f WHERE f.file_id = u.avatar_file_id),
	        ''
	    ) AS file_uuid,
	    COALESCE(
	        (SELECT f.file_extension FROM user_uploaded_file f WHERE f.file_id = u.avatar_file_id),
	        ''
	    ) AS file_extension;
	`

	var duplicateCount int

	row := r.db.QueryRow(ctx, query1, data.Email, userID)
	err = row.Scan(&duplicateCount)
	logging.Debugf(ctx, "%s query 1 has err: %v", funcName, err)
	if err != nil {
		return nil, fmt.Errorf("%s (check unique): %w", funcName, err)
	}
	if duplicateCount != 0 {
		return nil, fmt.Errorf("%s (check unique): %w", funcName, errs.ErrAlreadyExists)
	}

	newProfile = &models.UserProfile{}
	var fileUUID, fileExt string

	row = r.db.QueryRow(ctx, query2, data.Email, data.NewName, userID)
	err = row.Scan(&newProfile.ID,
		&newProfile.Name,
		&newProfile.Email,
		&newProfile.JoinedAt,
		&newProfile.UpdatedAt,
		&fileUUID,
		&fileExt)
	logging.Debugf(ctx, "%s query 2 has err: %v", funcName, err)
	if err != nil {
		return nil, fmt.Errorf("%s (action): %w", funcName, err)
	}

	newProfile.AvatarImageURL = uploads.JoinFileURL(fileUUID, fileExt, uploads.DefaultAvatarURL)

	return
}

func (r *UserRepository) SetUserAvatar(ctx context.Context, userID int64, avatarFileID int64) error {
	funcName := "SetUserAvatar"
	query := `
	UPDATE "user"
	SET avatar_file_id=$1
	WHERE u_id=$2;
	`
	tag, err := r.db.Exec(ctx, query, avatarFileID, userID)
	logging.Debugf(ctx, "%s query has err: ", funcName, err)
	if err != nil {
		return fmt.Errorf("%s (update): %w", funcName, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s (update): no rows affected", funcName)
	}
	return nil
}

// GetUserByEmail получает данные пользователя из базы по email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (user *models.UserProfile, err error) {
	funcName := "GetUserByEmail"
	query := `
	SELECT u.u_id, u.nickname, u.email,
	u.joined_at, u.updated_at,
	COALESCE(f.file_uuid::text, ''), COALESCE(f.file_extension, '')
	FROM "user" AS u
	LEFT JOIN user_uploaded_file AS f ON f.file_id=u.avatar_file_id
	WHERE u.email=$1;`

	user = &models.UserProfile{}
	var fileUUID, fileExt string
	err = r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.JoinedAt,
		&user.UpdatedAt,
		&fileUUID, &fileExt,
	)
	logging.Debugf(ctx, "%s query has err: %v", funcName, err)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrWrongCredentials
		}
		return nil, err
	}
	user.AvatarImageURL = uploads.JoinFileURL(fileUUID, fileExt, uploads.DefaultAvatarURL)
	return user, nil
}

// CreateUser создаёт пользователя (или не создаёт, если повторяются креды)
func (r *UserRepository) CreateUser(ctx context.Context, user *models.UserRegisterRequest) (newUser *models.UserProfile, err error) {
	funcName := "CreateUser"
	newUser = &models.UserProfile{}
	query := `INSERT INTO "user" (nickname, email, password_hash, csat_poll_dt)
              VALUES ($1, $2, NULL, (CURRENT_TIMESTAMP+$3)) RETURNING u_id, nickname, email, joined_at, updated_at`

	err = r.db.QueryRow(ctx, query, user.Name, user.Email, 24*7*time.Hour).Scan(
		&newUser.ID,
		&newUser.Name,
		&newUser.Email,
		&newUser.JoinedAt,
		&newUser.UpdatedAt,
	)
	logging.Debugf(ctx, "%s query has err: %v", funcName, err)
	return newUser, err
}

// CheckUniqueCredentials проверяет, существуют ли такие логин и email в базе
func (r *UserRepository) CheckUniqueCredentials(ctx context.Context, nickname string, email string) error {
	funcName := "CheckUniqueCredentials"
	query := `SELECT nickname, email FROM "user" WHERE nickname = $1 OR email=$2;`
	var emailCount, nicknameCount int
	rows, err := r.db.Query(ctx, query, nickname, email)
	logging.Debugf(ctx, "%s query has err: %v", funcName, err)
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
	}
	if nicknameCount > 0 {
		return fmt.Errorf("%s: %w", funcName, errs.ErrBusyNickname)
	}
	if emailCount > 0 {
		return fmt.Errorf("%s: %w", funcName, errs.ErrBusyEmail)
	}
	return nil
}

func (r *UserRepository) DeduplicateFile(ctx context.Context, file *models.UploadedFile) (fileNames []string, fileIDs []int64, err error) {
	return uploads.DeduplicateFile(ctx, r.db, file)
}
func (r *UserRepository) RegisterFile(ctx context.Context, file *models.UploadedFile) (fileID int64, fileUUID string, err error) {
	return uploads.RegisterFile(ctx, r.db, file)
}
