package repository

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/utils/logging"
	"RPO_back/internal/pkg/utils/uploads"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// GetUserProfile получает из базы профиль пользователя
func (r *BoardRepository) GetUserProfile(ctx context.Context, userID int64) (user *models.UserProfile, err error) {
	query := `
	SELECT
	u_id,
	nickname,
	email,
	joined_at,
	updated_at,
	COALESCE(f.file_uuid::text, ''),
	COALESCE(f.file_extension, '')
	FROM "user" AS u
	LEFT JOIN user_uploaded_file AS f ON f.file_id=u.avatar_file_id
	WHERE u_id=$1;
	`
	rows := r.db.QueryRow(ctx, query, userID)
	user = &models.UserProfile{}
	var avatarUUID, avatarExt string
	err = rows.Scan(&user.ID, &user.Name, &user.Email,
		&user.JoinedAt,
		&user.UpdatedAt,
		&avatarUUID,
		&avatarExt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetUserProfile: %w", errs.ErrNotFound)
		}
		return nil, fmt.Errorf("GetUserProfile: %w", err)
	}
	user.AvatarImageURL = uploads.JoinFileURL(avatarUUID, avatarExt, uploads.DefaultAvatarURL)
	return user, nil
}

// GetMemberPermissions (предназначено для внутреннего использования)
// Возвращает информацию о правах участника на конкретной доске
// Если пользователя нет на доске, возвращает errs.ErrNotPermitted
//
// verbose - если равен true, запрос получит профили пригласившего
// пользователя и пользователя, внёсшего последнее изменение, а также
// загрузит аватарки. false -
// поля AddedBy и UpdatedBy будут установлены в nil. Но если
// verbose равен true, ещё не факт, что указанные поля
// будут не nil
func (r *BoardRepository) GetMemberPermissions(ctx context.Context, boardID int64, memberUserID int64, verbose bool) (member *models.MemberWithPermissions, err error) {
	query := `
	WITH board_check AS (
		SELECT 1
		FROM kanban_column
		WHERE board_id=$2
	)
	SELECT
	ub.role,
	ub.added_at,
	ub.updated_at,
	COALESCE(ub.added_by, -1),
	COALESCE(ub.updated_by, -1)
	FROM user_to_board AS ub
	WHERE ub.u_id=$1
	AND ub.board_id=$2;
	`
	member = &models.MemberWithPermissions{}
	// Получение профиля пользователя
	userProfile, err := r.GetUserProfile(ctx, memberUserID)
	if err != nil {
		return nil, fmt.Errorf("GetMemberPermissions (getting user profile): %w", err)
	}
	member.User = userProfile
	var addedByID, updatedByID int64
	rows := r.db.QueryRow(ctx, query, memberUserID, boardID)
	err = rows.Scan(
		&member.Role,
		&member.AddedAt,
		&member.UpdatedAt,
		&addedByID,
		&updatedByID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetMemberPermissions (getting user perms): %w", errs.ErrNotPermitted)
		}
		return nil, fmt.Errorf("GetMemberPermissions (getting user perms): %w", err)
	}
	if verbose {
		if addedByID != -1 {
			adder, err := r.GetUserProfile(ctx, addedByID)
			if err != nil {
				return nil, fmt.Errorf("GetMemberPermissions (getting adder profile): %w", err)
			}
			member.AddedBy = adder
		}
		if updatedByID != -1 {
			updater, err := r.GetUserProfile(ctx, updatedByID)
			if err != nil {
				return nil, fmt.Errorf("GetMemberPermissions (getting updater profile): %w", err)
			}
			member.UpdatedBy = updater
		}
	}
	return member, nil
}

// GetMembersWithPermissions получает всех участников на конкретной
// доске с информацией об их правах и с разрешением профилей добавителя
// и пользователя, внёсшего последнее обновление в роль
func (r *BoardRepository) GetMembersWithPermissions(ctx context.Context, boardID int64, userID int64) (members []models.MemberWithPermissions, err error) {
	query := `
	SELECT

	m.u_id, m.nickname,
	m.email, m.joined_at, m.updated_at,

	ub.role, ub.added_at, ub.updated_at,

	COALESCE(adder.u_id,-1), adder.nickname, adder.email,
	adder.joined_at, adder.updated_at,

	COALESCE(updater.u_id,-1), updater.nickname, updater.email,
	updater.joined_at, updater.updated_at,

	COALESCE(f_m.file_uuid::text,''), COALESCE(f_m.file_extension,''),
	COALESCE(f_adder.file_uuid::text,''), COALESCE(f_adder.file_extension,''),
	COALESCE(f_updater.file_uuid::text,''), COALESCE(f_updater.file_extension,'')

	FROM "user" AS m
	JOIN user_to_board AS ub ON m.u_id=ub.u_id
	LEFT JOIN "user" AS adder ON adder.u_id=ub.added_by
	LEFT JOIN "user" AS updater ON updater.u_id=ub.updated_by

	LEFT JOIN user_uploaded_file AS f_m ON f_m.file_id=m.avatar_file_id
	LEFT JOIN user_uploaded_file AS f_adder ON f_adder.file_id=adder.avatar_file_id
	LEFT JOIN user_uploaded_file AS f_updater ON f_updater.file_id=updater.avatar_file_id
	WHERE ub.board_id=$1;
	`
	_, err = r.GetBoard(ctx, boardID, userID)
	if err != nil {
		return nil, fmt.Errorf("GetMembersWithPermissions (getting board): %w", errs.ErrNotFound)
	}
	rows, err := r.db.Query(ctx, query, boardID)
	logging.Debug(ctx, "GetMembersWithPermissions query has err: ", err)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("GetMembersWithPermissions (main query): %w", err)
	}
	for rows.Next() {
		var memberAvatarUUID, memberAvatarExt, adderAvatarUUID, adderAvatarExt, updaterAvatarUUID, updaterAvatarExt string
		field := models.MemberWithPermissions{}
		field.User = &models.UserProfile{}
		field.AddedBy = &models.UserProfile{}
		field.UpdatedBy = &models.UserProfile{}
		err := rows.Scan(
			&field.User.ID,
			&field.User.Name,
			&field.User.Email,
			&field.User.JoinedAt,
			&field.User.UpdatedAt,

			&field.Role, &field.AddedAt, &field.UpdatedAt,

			&field.AddedBy.ID,
			&field.AddedBy.Name,
			&field.AddedBy.Email,
			&field.AddedBy.JoinedAt,
			&field.AddedBy.UpdatedAt,

			&field.UpdatedBy.ID,
			&field.UpdatedBy.Name,
			&field.UpdatedBy.Email,
			&field.UpdatedBy.JoinedAt,
			&field.UpdatedBy.UpdatedAt,

			&memberAvatarUUID, &memberAvatarExt,
			&adderAvatarUUID, &adderAvatarExt,
			&updaterAvatarUUID, &updaterAvatarExt,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("GetMembersWithPermissions: %w", errs.ErrNotFound)
			}
			return nil, fmt.Errorf("GetMembersWithPermissions: %w", err)
		}
		if field.AddedBy.ID == -1 {
			field.AddedBy = nil
		}
		if field.UpdatedBy.ID == -1 {
			field.UpdatedBy = nil
		}

		field.User.AvatarImageURL = uploads.JoinFileURL(memberAvatarUUID, memberAvatarExt, uploads.DefaultAvatarURL)
		field.AddedBy.AvatarImageURL = uploads.JoinFileURL(adderAvatarUUID, adderAvatarExt, uploads.DefaultAvatarURL)
		field.UpdatedBy.AvatarImageURL = uploads.JoinFileURL(updaterAvatarUUID, updaterAvatarExt, uploads.DefaultAvatarURL)

		members = append(members, field)
	}
	return members, nil
}

// SetMemberRole устанавливает существующему участнику права (роль)
func (r *BoardRepository) SetMemberRole(ctx context.Context, userID int64, boardID int64, memberUserID int64, newRole string) (member *models.MemberWithPermissions, err error) {
	funcName := "SetMemberRole"
	query := `
	UPDATE user_to_board
	SET role=$1,
		updated_at=CURRENT_TIMESTAMP,
		updated_by=$2
	WHERE u_id=$3
	AND board_id=$4;
	`

	tag, err := r.db.Exec(ctx, query, memberUserID, boardID)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if tag.RowsAffected() == 0 {
		return nil, fmt.Errorf("%s (update): %w", funcName, errs.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s (update): %w", funcName, err)
	}
	member, err = r.GetMemberPermissions(ctx, boardID, memberUserID, true)
	if err != nil {
		return nil, fmt.Errorf("%s (get updated perms): %w", funcName, err)
	}
	return member, nil
}

// RemoveMember удаляет участника с доски
func (r *BoardRepository) RemoveMember(ctx context.Context, boardID int64, memberUserID int64) (err error) {
	query := `
	DELETE FROM user_to_board
	WHERE board_id=$1
	AND u_id=$2;
	`
	tag, err := r.db.Exec(ctx, query, boardID, memberUserID)
	logging.Debug(ctx, "RemoveMember query has err: ", err, " tag: ", tag)
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("RemoveMember: %w", errs.ErrNotFound)
	}
	if err != nil {
		return fmt.Errorf("RemoveMember: %w", err)
	}
	return nil
}

// AddMember добавляет участника на доску с правами "viewer"
func (r *BoardRepository) AddMember(ctx context.Context, boardID int64, adderID int64, memberUserID int64) (member *models.MemberWithPermissions, err error) {
	query := `
	INSERT INTO user_to_board (u_id, board_id, added_at, updated_at,
	last_visit_at, added_by, updated_by, "role") VALUES (
	$1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP,
	$3, $3, 'viewer'
	);
	`
	member, err = r.GetMemberPermissions(ctx, boardID, memberUserID, false)
	logging.Debug(ctx, "AddMember query has err: ", err)

	if (err != nil) && (!errors.Is(err, errs.ErrNotPermitted)) {
		return nil, fmt.Errorf("AddMember (get member): %w", err)
	}
	if err == nil {
		return nil, fmt.Errorf("AddMember (get member): %w", errs.ErrAlreadyExists)
	}
	_, err = r.db.Exec(ctx, query, memberUserID, boardID, adderID)
	if err != nil {
		return nil, fmt.Errorf("AddMember (insert): %w", err)
	}
	member, err = r.GetMemberPermissions(ctx, boardID, memberUserID, true)
	return member, err
}

// GetUserByNickname получает данные пользователя из базы по имени
func (r *BoardRepository) GetUserByNickname(ctx context.Context, nickname string) (user *models.UserProfile, err error) {
	query := `SELECT u_id, nickname, email, joined_at, updated_at FROM "user"
	WHERE nickname=$1;`
	user = &models.UserProfile{}
	err = r.db.QueryRow(ctx, query, nickname).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.JoinedAt,
		&user.UpdatedAt,
	)
	logging.Debug(ctx, "GetUserByNickname query has err: ", err)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetUserByNickname (query): %w", errs.ErrNotFound)
		}
		return nil, fmt.Errorf("GetUserByNickname (query): %w", err)
	}
	return user, nil
}

// GetMemberFromCard получает права пользователя из ID карточки
func (r *BoardRepository) GetMemberFromCard(ctx context.Context, userID int64, cardID int64) (role string, boardID int64, err error) {
	funcName := "GetMemberFromCard"
	query := `
	SELECT
		ub.role,
		ub.board_id
	FROM card AS c
	LEFT JOIN kanban_column AS col ON col.col_id=c.col_id
	LEFT JOIN board AS b ON b.board_id=col.board_id
	LEFT JOIN user_to_board AS ub ON ub.board_id=b.board_id
	WHERE c.card_id=$1 AND ub.u_id=$2;
	`

	row := r.db.QueryRow(ctx, query, cardID, userID)
	err = row.Scan(
		&role,
		&boardID,
	)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return "", 0, fmt.Errorf("%s (query): %w", funcName, err)
	}
	return role, boardID, nil
}

// GetMemberFromCheckListField получает права пользователя из ID поля чеклиста
func (r *BoardRepository) GetMemberFromCheckListField(ctx context.Context, userID int64, fieldID int64) (role string, boardID int64, cardID int64, err error) {
	funcName := "GetMemberFromCheckListField"
	query := `
	SELECT
	utb.role, b.board_id, c.card_id
	FROM checklist_field AS cf
	JOIN card AS c ON cf.card_id = c.card_id
	JOIN kanban_column AS kc ON kc.col_id = c.col_id
	JOIN board AS b ON b.board_id = kc.board_id
	JOIN user_to_board AS utb ON utb.board_id = b.board_id
	WHERE utb.u_id = $1 AND cf.checklist_field_id = $2;
	`

	err = r.db.QueryRow(ctx, query, userID, fieldID).Scan(
		&role, &boardID, &cardID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", 0, 0, fmt.Errorf("%s (query): %w", funcName, errs.ErrNotFound)
		}
		return "", 0, 0, fmt.Errorf("%s (query): %w", funcName, err)
	}

	return role, boardID, cardID, err
}

// GetMemberFromAttachment получает права пользователя из ID вложения
func (r *BoardRepository) GetMemberFromAttachment(ctx context.Context, userID int64, attachmentID int64) (role string, boardID int64, cardID int64, err error) {
	funcName := "GetMemberFromAttachment"
	query := `
		SELECT utb.role, b.board_id, c.card_id
		FROM card_attachment AS ca
		JOIN card AS c ON ca.card_id = c.card_id
		JOIN kanban_column AS kc ON c.col_id = kc.col_id
		JOIN board AS b ON kc.board_id = b.board_id
		JOIN user_to_board AS utb ON utb.board_id = b.board_id
		WHERE utb.u_id = $1 AND ca.attachment_id = $2;
	`

	err = r.db.QueryRow(ctx, query, userID, attachmentID).Scan(
		&role, &boardID, &cardID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", 0, 0, fmt.Errorf("%s (query): %w", funcName, errs.ErrNotFound)
		}
		return "", 0, 0, fmt.Errorf("%s (query): %w", funcName, err)
	}

	return role, boardID, cardID, err
}

// GetMemberFromColumn получает права пользователя из ID колонки
func (r *BoardRepository) GetMemberFromColumn(ctx context.Context, userID int64, columnID int64) (role string, boardID int64, err error) {
	funcName := "GetMemberFromColumn"
	query := `
	SELECT utb.role, b.board_id
	FROM kanban_column AS kc
	JOIN board AS b ON kc.board_id = b.board_id
	JOIN user_to_board AS utb ON utb.board_id = b.board_id
	WHERE utb.u_id = $1 AND kc.col_id = $2;
	`

	err = r.db.QueryRow(ctx, query, userID, columnID).Scan(
		&role, &boardID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", 0, fmt.Errorf("%s (query): %w", funcName, errs.ErrNotFound)
		}
		return "", 0, fmt.Errorf("%s (query): %w", funcName, err)
	}

	return role, boardID, err
}

// GetMemberFromComment получает права пользователя из ID комментария
func (r *BoardRepository) GetMemberFromComment(ctx context.Context, userID int64, commentID int64) (role string, boardID int64, cardID int64, err error) {
	funcName := "GetMemberFromComment"
	query := `
		SELECT utb.role, b.board_id, c.card_id
		FROM card_comment AS cc
		JOIN card AS c ON cc.card_id = c.card_id
		JOIN kanban_column AS kc ON c.col_id = kc.col_id
		JOIN board AS b ON kc.board_id = b.board_id
		JOIN user_to_board AS utb ON utb.board_id = b.board_id
		WHERE utb.u_id = $1 AND cc.comment_id = $2;
	`

	err = r.db.QueryRow(ctx, query, userID, commentID).Scan(
		&role, &boardID, &cardID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", 0, 0, fmt.Errorf("%s (query): %w", funcName, errs.ErrNotFound)
		}
		return "", 0, 0, fmt.Errorf("%s (query): %w", funcName, err)
	}

	return role, boardID, cardID, err
}

// GetCardCheckList получает чеклисты для карточки
func (r *BoardRepository) GetCardCheckList(ctx context.Context, cardID int64) (checkList []models.CheckListField, err error) {
	funcName := "GetCardCheckList"
	query := `
		SELECT cf.checklist_field_id, cf.title, cf.created_at, cf.is_done
		FROM checklist_field AS cf
		JOIN card AS c ON cf.card_id = cf.card_id
		WHERE c.card_id = $1
		ORDER BY cf.order_index;
	`

	checkList = make([]models.CheckListField, 0)

	rows, err := r.db.Query(ctx, query, cardID)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s (query): %w", funcName, errs.ErrNotFound)
		}
		return nil, fmt.Errorf("%s (query): %w", funcName, err)
	}

	for rows.Next() {
		field := models.CheckListField{}
		if err := rows.Scan(&field.ID, &field.Title, &field.CreatedAt, &field.IsDone); err != nil {
			return nil, fmt.Errorf("%s (scan): %w", funcName, err)
		}

		checkList = append(checkList, field)
	}

	return checkList, nil
}

// GetCardAssignedUsers получает пользователей, назначенных на карточку
func (r *BoardRepository) GetCardAssignedUsers(ctx context.Context, cardID int64) (assignedUsers []models.UserProfile, err error) {
	funcName := "GetCardAssignedUsers"
	query := `
		SELECT u.u_id,
		u.nickname,
		u.email,
		u.joined_at,
		u.updated_at,
		COALESCE(f.file_uuid::text, ''),
		COALESCE(f.file_extension::text, '')
		FROM card_user_assignment AS cua
		JOIN "user" AS u ON cua.u_id = u.u_id
		LEFT JOIN user_uploaded_file AS f ON f.file_id=u.avatar_file_id
		WHERE cua.card_id = $1;
	`
	rows, err := r.db.Query(ctx, query, cardID)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s (query): %w", funcName, errs.ErrNotFound)
		}
		return nil, fmt.Errorf("%s (query): %w", funcName, err)
	}

	assignedUsers = make([]models.UserProfile, 0)

	for rows.Next() {
		assignedProfile := models.UserProfile{}
		var avatarUUID, avatarExt string
		if err := rows.Scan(&assignedProfile.ID,
			&assignedProfile.Name,
			&assignedProfile.Email,
			&assignedProfile.JoinedAt,
			&assignedProfile.UpdatedAt,
			&avatarUUID,
			&avatarExt,
		); err != nil {
			return nil, fmt.Errorf("%s (scan): %w", funcName, err)
		}

		assignedProfile.AvatarImageURL = uploads.JoinFileURL(avatarUUID, avatarExt, uploads.DefaultAvatarURL)
		assignedUsers = append(assignedUsers, assignedProfile)
	}

	return assignedUsers, nil
}

// GetCardComments получает комментарии, оставленные на карточку
func (r *BoardRepository) GetCardComments(ctx context.Context, cardID int64) (comments []models.Comment, err error) {
	funcName := "GetCardComments"
	query := `
		SELECT cc.comment_id,
		cc.title,
		cc.created_at,
		cc.is_edited,

		u.u_id,
		u.nickname,
		u.email,
		u.joined_at,
		u.updated_at,
		COALESCE(f.file_uuid::text, ''),
		COALESCE(f.file_extension::text, '')

		FROM card_comment AS cc
		JOIN "user" AS u ON cc.created_by=u.u_id
		LEFT JOIN user_uploaded_file AS f ON f.file_id=u.avatar_file_id
		WHERE cc.card_id = $1;
	`

	comments = make([]models.Comment, 0)

	rows, err := r.db.Query(ctx, query, cardID)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s (query): %w", funcName, errs.ErrNotFound)
		}
		return nil, fmt.Errorf("%s (query): %w", funcName, err)
	}

	for rows.Next() {
		uP := models.UserProfile{}
		c := models.Comment{}
		var avatarUUID, avatarExt string

		if err := rows.Scan(&c.ID, &c.Text, &c.CreatedAt, &c.IsEdited, &uP.ID,
			&uP.Name, &uP.Email, &uP.JoinedAt,
			&uP.UpdatedAt, &avatarUUID, &avatarExt); err != nil {
			return nil, fmt.Errorf("%s (scan): %w", funcName, err)
		}
		uP.AvatarImageURL = uploads.JoinFileURL(avatarUUID, avatarExt, uploads.DefaultAvatarURL)
		c.CreatedBy = &uP

		comments = append(comments, c)
	}

	return comments, nil
}

// GetCardAttachments получает вложения к карточке
func (r *BoardRepository) GetCardAttachments(ctx context.Context, cardID int64) (attachments []models.Attachment, err error) {
	query := `
		SELECT ca.attachment_id,
		ca.original_name,
		ca.created_at,
		COALESCE(f.file_uuid::text, ''),
		COALESCE(f.file_extension::text, '')
		FROM card_attachment AS ca
		LEFT JOIN user_uploaded_file AS f ON f.file_id=ca.file_id
		WHERE ca.card_id = $1;
	`

	attachments = make([]models.Attachment, 0)

	rows, err := r.db.Query(ctx, query, cardID)
	logging.Debug(ctx, "GetCardAttachments query has err: ", err)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetCardAttachments (query): %w", errs.ErrNotFound)
		}
		return nil, fmt.Errorf("GetCardAttachments (query): %w", err)
	}

	for rows.Next() {
		a := models.Attachment{}
		var avatarUUID, avatarExt string

		if err := rows.Scan(&a.ID, &a.OriginalName, &a.CreatedAt, &avatarUUID, &avatarExt); err != nil {
			return nil, fmt.Errorf("GetCardAssignedUsers (scan): %w", err)
		}
		a.FileName = uploads.JoinFileURL(avatarUUID, avatarExt, uploads.DefaultAvatarURL)

		attachments = append(attachments, a)
	}

	return attachments, nil
}

// GetCardsForMove получает списки карточек на двух колонках. (карточки неполные)
// Нужно для Drag-n-Drop (колонки откуда перемещаем и куда)
func (r *BoardRepository) GetCardsForMove(ctx context.Context, col1ID int64, col2ID *int64) (column1 []models.Card, column2 []models.Card, err error) {
	query := `
	SELECT c.card_id, c.col_id, c.order_index
	FROM card AS c
	WHERE c.col_id = $1 OR c.col_id = $2
	ORDER BY c.order_index;
	`

	rows, err := r.db.Query(ctx, query, col1ID, col2ID)
	logging.Debug(ctx, "GetCardsForMove query has err: ", err)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, fmt.Errorf("GetCardsForMove (query): %w", errs.ErrNotFound)
		}
		return nil, nil, fmt.Errorf("GetCardsForMove (query): %w", err)
	}

	for rows.Next() {
		c := models.Card{}

		if err := rows.Scan(&c.ID, &c.ColumnID, c.OrderIndex); err != nil {
			return nil, nil, fmt.Errorf("GetCardsForMove (scan): %w", err)
		}

		if c.ColumnID == col1ID {
			column1 = append(column1, c)
		} else if col2ID != nil && c.ColumnID == *col2ID {
			column2 = append(column2, c)
		}
	}

	return column1, column2, nil
}

// GetColumnsForMove получает список всех колонок, чтобы сделать Drag-n-Drop
func (r *BoardRepository) GetColumnsForMove(ctx context.Context, boardID int64) (columns []models.Column, err error) {
	funcName := "GetColumnsForMove"
	query := `
	SELECT kc.col_id, kc.title, kc.order_index
	FROM kanban_column AS kc
	WHERE kc.board_id = $1
	ORDER BY kc.order_index;
	`

	rows, err := r.db.Query(ctx, query, boardID)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s (query): %w", funcName, errs.ErrNotFound)
		}
		return nil, fmt.Errorf("%s (query): %w", funcName, err)
	}

	for rows.Next() {
		c := models.Column{}

		if err := rows.Scan(&c.ID, &c.Title, &c.OrderIndex); err != nil {
			return nil, fmt.Errorf("%s (scan): %w", funcName, err)
		}

		columns = append(columns, c)
	}

	return columns, nil
}

// RearrangeCards обновляет позиции всех карточек колонки, чтобы сделать порядок, как в слайсе
func (r *BoardRepository) RearrangeCards(ctx context.Context, columnID int64, cards []models.Card) (err error) {
	funcName := "RearrangeCards"
	query := `
	WITH update_position_cards AS (
		UPDATE card SET order_index = $1 WHERE col_id = $2
	),
	update_board AS (
		UPDATE board SET updated_at = CURRENT_TIMESTAMP WHERE board_id = (
			SELECT board_id FROM kanban_column WHERE col_id = $2
		)
	)
	SELECT;
	`
	batch := &pgx.Batch{}
	for _, card := range cards {
		batch.Queue(query, card.OrderIndex)
	}

	br := r.db.SendBatch(ctx, batch)
	err = br.Close()
	logging.Debug(ctx, funcName, " batch query has err: ", err)
	if err != nil {
		return fmt.Errorf("%s (batch query): %w", funcName, err)
	}
	return nil
}

// RearrangeColumns обновляет позиции всех колонок, чтобы сделать порядок, как в слайсе
func (r *BoardRepository) RearrangeColumns(ctx context.Context, columns []models.Column) (err error) {
	panic("not implemented")
	funcName := "RearrangeColumns"
	query := ` WITH`
	batch := &pgx.Batch{}
	for _, col := range columns {
		batch.Queue(query, col.OrderIndex)
	}

	br := r.db.SendBatch(ctx, batch)
	err = br.Close()
	logging.Debug(ctx, funcName, " batch query has err: ", err)
	if err != nil {
		return fmt.Errorf("%s (batch query): %w", funcName, err)
	}
	return nil
}

// RearrangeCheckList устанавливает порядок полей чеклиста как в слайсе
func (r *BoardRepository) RearrangeCheckList(ctx context.Context, fields []models.CheckListField) (err error) {
	panic("not implemented")
	funcName := "RearrangeCheckList"
	query := ``
	batch := &pgx.Batch{}
	for _, col := range fields {
		batch.Queue(query, col.OrderIndex)
	}

	br := r.db.SendBatch(ctx, batch)
	err = br.Close()
	logging.Debug(ctx, funcName, " batch query has err: ", err)
	if err != nil {
		return fmt.Errorf("%s (batch query): %w", funcName, err)
	}
	return nil
}

// AssignUserToCard назначает пользователя на карточку
func (r *BoardRepository) AssignUserToCard(ctx context.Context, cardID int64, assignedUserID int64) (assignedUser *models.UserProfile, err error) {
	funcName := "AssignUserToCard"
	query := `
		WITH update_card_user_assignment AS (
			UPDATE card_user_assignment SET u_id = $2 WHERE card_id = $1
		),
		update_card AS (
			UPDATE "card" SET updated_at=CURRENT_TIMESTAMP WHERE card_id = $1
		),
		update_board AS (
			UPDATE board SET updated_at=CURRENT_TIMESTAMP WHERE board_id = (
				SELECT b.board_id
				FROM card AS c
				JOIN kanban_column AS kc ON c.col_id=kc.col_id
				JOIN board AS b ON b.board_id = kc.board_id
				WHERE c.card_id=$1
			)
		)
		SELECT u.u_id, u.nickname, u.email, u.joined_at, u.updated_at,
		COALESCE(f.file_uuid::text, ''), COALESCE(f.file_extension::text, '')
		FROM card_user_assignment AS cua
		JOIN "user" AS u ON cua.u_id = u.u_id
		LEFT JOIN user_uploaded_file AS f ON f.file_id=u.avatar_file_id
		WHERE cua.card_id = $1;
	`

	tag, err := r.db.Exec(ctx, query)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return nil, fmt.Errorf("%s (query): %w", funcName, err)
	}
	if tag.RowsAffected() == 0 {
		return nil, fmt.Errorf("%s (query): no rows affected", funcName)
	}
	return assignedUser, nil
}

// DeassignUserFromCard убирает назначение пользователя
func (r *BoardRepository) DeassignUserFromCard(ctx context.Context, cardID int64, assignedUserID int64) (err error) {
	panic("not implemented")
	funcName := "DeassignUserFromCard"
	query := `
		WITH delete_card_user_assignment AS (
			DELETE FROM card_user_assignment WHERE card_id = $1 AND u_id = $2
		),
		update_card AS (
			UPDATE "card" SET updated_at=CURRENT_TIMESTAMP WHERE card_id = $1
		),
		update_board AS (
			UPDATE board SET updated_at=CURRENT_TIMESTAMP WHERE board_id = (
				SELECT b.board_id
				FROM card AS c
				JOIN kanban_column AS kc ON c.col_id=kc.col_id
				JOIN board AS b ON b.board_id = kc.board_id
				WHERE c.card_id=$1
			)
		)
		)
		)
	`
	tag, err := r.db.Exec(ctx, query)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return fmt.Errorf("%s (query): %w", funcName, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s (query): no rows affected", funcName)
	}
	return nil
}

// CreateComment добавляет на карточку комментарий
func (r *BoardRepository) CreateComment(ctx context.Context, userID int64, cardID int64, comment *models.CommentRequest) (newComment *models.Comment, err error) {
	funcName := "CreateComment"
	query := `
		WITH insert_comment AS (
			INSERT INTO card_comment (card_id, created_by, title) VALUES ($1, $2, $3)
			RETURNING comment_id, title, is_edited, created_by, created_at
		),
		update_card AS (
			UPDATE "card" SET updated_at=CURRENT_TIMESTAMP WHERE card_id = $1
		),
		update_board AS (
			UPDATE board SET updated_at=CURRENT_TIMESTAMP WHERE board_id = (
				SELECT b.board_id
				FROM card AS c
				JOIN kanban_column AS kc ON c.col_id=kc.col_id
				JOIN board AS b ON b.board_id = kc.board_id
				WHERE c.card_id=$1
			)
		)
		SELECT i.comment_id, i.title,
			i.is_edited, i.created_at,
			u.u_id, u.nickname, u.email,u.joined_at, u.updated_at,
			COALESCE(f.file_uuid::text, ''),
			COALESCE(f.file_extension, '')
		FROM insert_comment AS i
		JOIN "user" AS u ON u.u_id=i.created_by
		LEFT JOIN user_uploaded_file AS f ON f.file_id=u.avatar_file_id;

	`

	newComment = &models.Comment{}
	newComment.CreatedBy = &models.UserProfile{}
	row := r.db.QueryRow(ctx, query, cardID, userID, comment.Text)
	var fileUUID, fileExtension string
	err = row.Scan(
		&newComment.ID,
		&newComment.Text,
		&newComment.IsEdited,
		&newComment.CreatedAt,
		&newComment.CreatedBy.ID,
		&newComment.CreatedBy.Name,
		&newComment.CreatedBy.Email,
		&newComment.CreatedBy.JoinedAt,
		&newComment.CreatedBy.UpdatedAt,
		&fileUUID,
		&fileExtension,
	)
	newComment.CreatedBy.AvatarImageURL = uploads.JoinFileURL(fileUUID, fileExtension, uploads.DefaultAvatarURL)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return nil, fmt.Errorf("%s (query): %w", err)
	}
	return newComment, err
}

// UpdateComment редактирует комментарий
func (r *BoardRepository) UpdateComment(ctx context.Context, commentID int64, update *models.CommentRequest) (updatedComment *models.Comment, err error) {
	funcName := "UpdateComment"
	query := `
	WITH update_comment AS (
		UPDATE card_comment
		SET updated_at=CURRENT_TIMESTAMP, title=$2, is_edited=TRUE
		WHERE comment_id = $1
		RETURNING comment_id, title, is_edited, created_at, created_by
	),
	update_card AS (
		UPDATE "card" SET updated_at=CURRENT_TIMESTAMP WHERE card_id = (
			SELECT card_id FROM card_comment WHERE comment_id=$1
		)
	),
	update_board AS (
		UPDATE board SET updated_at=CURRENT_TIMESTAMP WHERE board_id = (
			SELECT b.board_id
			FROM card_comment AS cc
			JOIN card AS c ON cc.card_id=c.card_id
			JOIN kanban_column AS kc ON c.col_id=kc.col_id
			JOIN board AS b ON b.board_id = kc.board_id
			WHERE cc.comment_id=$1
		)
	)
	SELECT c.comment_id, c.title,
		c.is_edited, c.created_at,
		u.u_id, u.nickname, u.email,u.joined_at, u.updated_at,
		COALESCE(f.file_uuid::text, ''),
		COALESCE(f.file_extension, '')
	FROM update_comment AS c
	JOIN "user" AS u ON u.u_id=i.created_by
	LEFT JOIN user_uploaded_file AS f ON f.file_id=u.avatar_file_id;
	`

	updatedComment = &models.Comment{}
	row := r.db.QueryRow(ctx, query, commentID, update.Text)
	err = row.Scan()
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return nil, fmt.Errorf("%s (query): %w", err)
	}
	return updatedComment, nil
}

// DeleteComment удаляет комментарий
func (r *BoardRepository) DeleteComment(ctx context.Context, commentID int64) (err error) {
	funcName := "DeleteComment"
	query := `
	WITH delete_comment AS (
		DELETE FROM card_comment WHERE comment_id=$1
	),
	update_card AS (
		UPDATE "card" SET updated_at=CURRENT_TIMESTAMP WHERE card_id=(
			SELECT card_id FROM card_comment WHERE comment_id=$1
		)
	),
	update_board AS (
		UPDATE board SET updated_at=CURRENT_TIMESTAMP WHERE board_id=(
			SELECT b.board_id FROM card_comment AS c
			JOIN card AS cc ON c.card_id=cc.card_id
			JOIN kanban_column AS k ON k.col_id=cc.col_id
			JOIN board AS b ON b.board_id=k.board_id
			WHERE c.comment_id=$1
		)
	)
	SELECT;
	`

	tag, err := r.db.Exec(ctx, query)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return fmt.Errorf("%s (query): %w", funcName, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s (query): no rows affected", funcName)
	}
	return nil
}

// CreateCheckListField создаёт поле чеклиста и добавляет его в конец
func (r *BoardRepository) CreateCheckListField(ctx context.Context, cardID int64, field *models.CheckListFieldPostRequest) (newField *models.CheckListField, err error) {
	funcName := "CreateCheckListField"
	query := `
	WITH insert_field AS (
		INSERT INTO checklist_field (card_id, title, order_index) VALUES ($1, $2, 12345)
		RETURNING checklist_field_id AS id
	)
	SELECT id FROM insert_field;
	`

	newField = &models.CheckListField{}
	row := r.db.QueryRow(ctx, query, cardID, field.Title)
	err = row.Scan(&newField.ID) //&newField.Title, &newField.CreatedAt, &newField.IsDone,
	newField.Title = field.Title

	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return nil, fmt.Errorf("%s (query): %w", err)
	}
	return newField, nil
}

// UpdateCheckListField обновляет одно поле чеклиста
func (r *BoardRepository) UpdateCheckListField(ctx context.Context, fieldID int64, update *models.CheckListFieldPatchRequest) (updatedField *models.CheckListField, err error) {
	funcName := "UpdateCheckListField"
	query := `
	WITH update_field AS (
		UPDATE checklist_field SET title=COALESCE($2,title), is_done=COALESCE($3,is_done) WHERE checklist_field_id = $1
	),
	update_card AS (
		UPDATE "card" SET updated_at=CURRENT_TIMESTAMP WHERE card_id = (
			SELECT card_id FROM checklist_field WHERE checklist_field_id=$1
		)
	),
	update_board AS (
		UPDATE board SET updated_at=CURRENT_TIMESTAMP WHERE board_id = (
			SELECT b.board_id
			FROM checklist_field AS cf
			JOIN card AS c ON cf.card_id=c.card_id
			JOIN kanban_column AS kc ON c.col_id=kc.col_id
			JOIN board AS b ON b.board_id = kc.board_id
			WHERE cf.checklist_field_id=$1
		)
	)
	SELECT f.checklist_field_id, f.title, f.created_at, f.is_done
	FROM checklist_field AS f
	WHERE f.checklist_field_id=$1;
	`

	updatedField = &models.CheckListField{}
	row := r.db.QueryRow(ctx, query, fieldID, update.Title, update.IsDone)
	err = row.Scan(&updatedField.ID, &updatedField.Title, &updatedField.CreatedAt, &updatedField.IsDone)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return nil, fmt.Errorf("%s (query): %w", err)
	}
	if update.IsDone != nil {
		updatedField.IsDone = *update.IsDone
	}
	if update.Title != nil {
		updatedField.Title = *update.Title
	}
	return updatedField, nil
}

func (r *BoardRepository) DeleteCheckListField(ctx context.Context, fieldID int64) error {
	funcName := "UpdateCheckListField"
	query := `
	DELETE FROM checklist_field
	WHERE checklist_field_id=$1;`

	tag, err := r.db.Exec(ctx, query, fieldID)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return fmt.Errorf("%s (query): %w", funcName, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s (query): no rows affected", funcName)
	}
	return nil
}

// SetCardCover устанавливает файл обложки карточки
func (r *BoardRepository) SetCardCover(ctx context.Context, userID int64, cardID int64, file *models.UploadedFile) (updatedCard *models.Card, err error) {
	funcName := "SetCardCover"
	query := `
	WITH update_cover AS (
		UPDATE "card" SET updated_at=CURRENT_TIMESTAMP, cover_file_id=$1 WHERE card_id = $2
	),
	update_board AS (
		UPDATE board SET updated_at=CURRENT_TIMESTAMP WHERE board_id = (
			SELECT b.board_id
			FROM card AS c
			JOIN kanban_column AS kc ON c.col_id=kc.col_id
			JOIN board AS b ON b.board_id = kc.board_id
			WHERE c.card_id=$1
		)
	)
	SELECT;
	`

	updatedCard = &models.Card{}
	row := r.db.QueryRow(ctx, query)
	err = row.Scan()
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return nil, fmt.Errorf("%s (query): %w", err)
	}
	return updatedCard, nil
}

// RemoveCardCover удаляет обложку карточки
func (r *BoardRepository) RemoveCardCover(ctx context.Context, cardID int64) (err error) {
	funcName := "RemoveCardCover"
	query := `
	WITH delete_cover AS (
		UPDATE "card" SET updated_at=CURRENT_TIMESTAMP, cover_file_id=NULL WHERE card_id = $1
	),
	update_board AS (
		UPDATE board SET updated_at=CURRENT_TIMESTAMP WHERE board_id = (
			SELECT b.board_id
			FROM card AS c
			JOIN kanban_column AS kc ON c.col_id=kc.col_id
			JOIN board AS b ON b.board_id = kc.board_id
			WHERE c.card_id=$1
		)
	)
	SELECT;
	`

	tag, err := r.db.Exec(ctx, query)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return fmt.Errorf("%s (query): %w", funcName, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s (query): no rows affected", funcName)
	}
	return nil
}

// AddAttachment добавляет файл вложения в карточку
func (r *BoardRepository) AddAttachment(ctx context.Context, userID int64, cardID int64, file *models.UploadedFile) (newAttachment *models.Attachment, err error) {
	funcName := "AddAttachment"
	query := `
	WITH insert_attachment AS (
		INSERT INTO card_attachment (card_id, file_id, original_name, attached_by) VALUES ($1, $2, $3, $4)
		RETURNING attachment_id
	),
	update_card AS (
		UPDATE "card" SET updated_at=CURRENT_TIMESTAMP WHERE card_id = $1
	),
	update_board AS (
		UPDATE board SET updated_at=CURRENT_TIMESTAMP WHERE board_id = (
			SELECT b.board_id
			FROM card_attachment AS ca
			JOIN card AS c ON ca.card_id=c.card_id
			JOIN kanban_column AS kc ON c.col_id=kc.col_id
			JOIN board AS b ON b.board_id = kc.board_id
			WHERE ca.attachment_id=$1
		)
	)
	SELECT uuf.file_uuid, uuf.file_extension
	FROM card_attachment AS ca
	JOIN user_uploaded_file AS uuf ON ca.file_id = uuf.file_id
	WHERE ca.file_id = $2;
	`

	newAttachment = &models.Attachment{}
	row := r.db.QueryRow(ctx, query)
	err = row.Scan()
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return nil, fmt.Errorf("%s (query): %w", err)
	}
	return newAttachment, nil
}

// RemoveAttachment удаляет вложение
func (r *BoardRepository) RemoveAttachment(ctx context.Context, attachmentID int64) (err error) {
	funcName := "RemoveAttachment"
	query := `
	WITH delete_attachment AS (
		DELETE FROM card_attachment WHERE attachment_id = $1
	),
	update_card AS (
		UPDATE card SET updated_at=CURRENT_TIMESTAMP WHERE card_id = (
			SELECT card_id FROM card_attachment WHERE attachment_id=$1
		)
	),
	update_board AS (
		UPDATE board SET updated_at=CURRENT_TIMESTAMP WHERE board_id = (
			SELECT b.board_id
			FROM card_attachment AS ca
			JOIN card AS c ON ca.card_id=c.card_id
			JOIN kanban_column AS kc ON c.col_id=kc.col_id
			JOIN board AS b ON b.board_id = kc.board_id
			WHERE ca.attachment_id=$1
		)
	)
	SELECT;
	`
	tag, err := r.db.Exec(ctx, query)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return fmt.Errorf("%s (query): %w", funcName, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s (query): no rows affected", funcName)
	}
	return nil
}

// PullInviteLink заменяет для доски индивидуальную ссылку-приглашение и возвращает новую ссылку
func (r *BoardRepository) PullInviteLink(ctx context.Context, userID int64, boardID int64) (link *models.InviteLink, err error) {
	funcName := "PullInviteLink"
	query := `
	WITH update_invite_link AS (
		UPDATE user_to_board SET invite_uuid = uuid_generate_v4() WHERE u_id = $1 AND board_id = $2
	)
	SELECT;
	`

	link = &models.InviteLink{}
	row := r.db.QueryRow(ctx, query)
	err = row.Scan()
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return nil, fmt.Errorf("%s (query): %w", err)
	}
	return link, nil
}

// DeleteInviteLink удаляет ссылку-приглашение
func (r *BoardRepository) DeleteInviteLink(ctx context.Context, userID int64, boardID int64) (err error) {
	funcName := "DeleteInviteLink"
	query := `
	WITH delete_invite_link AS (
		UPDATE user_to_board SET invite_uuid = NULL WHERE u_id = $1 AND board_id = $2
	)
	SELECT;
	`

	tag, err := r.db.Exec(ctx, query)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return fmt.Errorf("%s (query): %w", funcName, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s (query): no rows affected", funcName)
	}
	return nil
}

// FetchInvite возвращает информацию о доске, куда чела пригласили
func (r *BoardRepository) FetchInvite(ctx context.Context, inviteUUID string) (board *models.Board, err error) {
	funcName := "FetchInvite"
	query := `
		SELECT b.board_id, b.name, b.created_at, b.updated_at, ub.last_visit_at,
		COALESCE(f.file_uuid::text, ''), COALESCE(f.file_extension, '')
		FROM user_to_board AS utb
		JOIN board AS b ON utb.board_id = b.board_id
		LEFT JOIN user_uploaded_file AS f ON f.file_id=b.background_image_id
		WHERE utb.invite_uuid = $1;
	`

	board = &models.Board{}
	row := r.db.QueryRow(ctx, query)
	err = row.Scan()
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return nil, fmt.Errorf("%s (query): %w", err)
	}
	return board, nil
}

// AcceptInvite добавляет приглашённого пользователя на доску с правами зрителя
func (r *BoardRepository) AcceptInvite(ctx context.Context, userID int64, boardID int64, invitedUserID int64, inviteUUID string) (board *models.Board, err error) {
	funcName := "AcceptInvite"
	query := `
	WITH update_board AS (
		UPDATE board SET updated_at=CURRENT_TIMESTAMP WHERE board_id = $1
	),
	update_user_to_board AS (
		INSERT INTO user_to_board (u_id, board_id, role) VALUES ($2, $1, 'viewer')
	)
	SELECT
		b.board_id,
        b.name,
        b.created_at,
        b.updated_at,
        ub.last_visit_at,
        COALESCE(file.file_uuid::text,''),
        COALESCE(file.file_extension,'')
    FROM board AS b
    LEFT JOIN user_to_board AS ub ON ub.board_id = b.board_id AND ub.u_id = $1
    LEFT JOIN user_uploaded_file AS file ON file.file_id=b.background_image_id
    WHERE b.board_id = $1;
	`

	board = &models.Board{}
	row := r.db.QueryRow(ctx, query)
	err = row.Scan()
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return nil, fmt.Errorf("%s (query): %w", err)
	}
	return board, nil
}
