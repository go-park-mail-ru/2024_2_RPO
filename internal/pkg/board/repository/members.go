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
		if errors.Is(err, pgx.ErrNoRows) {
			return "", 0, fmt.Errorf("%s (query): %w", funcName, errs.ErrNotFound)
		}
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
		JOIN card AS c ON c.card_id = cf.card_id
		WHERE c.card_id = $1;
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
func (r *BoardRepository) GetCardsForMove(ctx context.Context, destColumnID int64, cardID *int64) (columnFrom []models.Card, columnTo []models.Card, err error) {
	query := `
	SELECT c.card_id, c.col_id
	FROM card AS c
	WHERE c.col_id = $1 OR c.col_id = (SELECT col_id FROM card WHERE card_id=$2)
	ORDER BY c.order_index;
	`

	rows, err := r.db.Query(ctx, query, destColumnID, cardID)
	logging.Debug(ctx, "GetCardsForMove query has err: ", err)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, fmt.Errorf("GetCardsForMove (query): %w", errs.ErrNotFound)
		}
		return nil, nil, fmt.Errorf("GetCardsForMove (query): %w", err)
	}

	for rows.Next() {
		c := models.Card{}

		if err := rows.Scan(&c.ID, &c.ColumnID); err != nil {
			return nil, nil, fmt.Errorf("GetCardsForMove (scan): %w", err)
		}

		if c.ColumnID == destColumnID {
			columnTo = append(columnTo, c)
		} else {
			columnFrom = append(columnFrom, c)
		}
	}

	if len(columnFrom) == 0 {
		columnFrom = columnTo
	}

	return columnFrom, columnTo, nil
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
func (r *BoardRepository) RearrangeCards(ctx context.Context, column1 []models.Card, column2 []models.Card) (err error) {
	funcName := "RearrangeCards"
	query := `
	UPDATE card SET order_index=$1, col_id=$2 WHERE card_id=$3;
	`
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s (begin): %w", funcName, err)
	}

	batch := &pgx.Batch{}
	for idx, card := range column1 {
		batch.Queue(query, idx, card.ColumnID, card.ID)
	}
	for idx, card := range column2 {
		batch.Queue(query, idx, card.ColumnID, card.ID)
	}

	br := tx.SendBatch(ctx, batch)
	err = br.Close()
	logging.Debug(ctx, funcName, " batch query has err: ", err)
	if err != nil {
		return fmt.Errorf("%s (batch query): %w", funcName, err)
	}

	err = tx.Commit(ctx)
	logging.Debug(ctx, funcName, " commit has err: ", err)
	if err != nil {
		return fmt.Errorf("%s (commit): %w", funcName, err)
	}

	return nil
}

// RearrangeColumns обновляет позиции всех колонок, чтобы сделать порядок, как в слайсе
func (r *BoardRepository) RearrangeColumns(ctx context.Context, columns []models.Column) (err error) {
	funcName := "RearrangeColumns"
	query := `
	UPDATE kanban_column SET order_index=$1 WHERE col_id=$2;
	`
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s (begin): %w", funcName, err)
	}

	batch := &pgx.Batch{}
	for idx, col := range columns {
		batch.Queue(query, idx, col.ID)
	}

	br := tx.SendBatch(ctx, batch)
	err = br.Close()
	logging.Debug(ctx, funcName, " batch query has err: ", err)
	if err != nil {
		return fmt.Errorf("%s (batch query): %w", funcName, err)
	}

	err = tx.Commit(ctx)
	logging.Debug(ctx, funcName, " commit has err: ", err)
	if err != nil {
		return fmt.Errorf("%s (commit): %w", funcName, err)
	}

	return nil
}

// // RearrangeCheckList устанавливает порядок полей чеклиста как в слайсе
// func (r *BoardRepository) RearrangeCheckList(ctx context.Context, fields []models.CheckListField) (err error) {
// 	panic("not implemented")
// 	funcName := "RearrangeCheckList"
// 	query := ``
// 	batch := &pgx.Batch{}
// 	for _, col := range fields {
// 		batch.Queue(query, col.OrderIndex)
// 	}

// 	br := r.db.SendBatch(ctx, batch)
// 	err = br.Close()
// 	logging.Debug(ctx, funcName, " batch query has err: ", err)
// 	if err != nil {
// 		return fmt.Errorf("%s (batch query): %w", funcName, err)
// 	}
// 	return nil
// }

// AssignUserToCard назначает пользователя на карточку
func (r *BoardRepository) AssignUserToCard(ctx context.Context, cardID int64, assignedUserID int64) (assignedUser *models.UserProfile, err error) {
	funcName := "AssignUserToCard"
	query := `
		WITH update_card_user_assignment AS (
			INSERT INTO card_user_assignment (card_id, u_id) VALUES ($1, $2)
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
		FROM "user" AS u
		LEFT JOIN user_uploaded_file AS f ON f.file_id=u.avatar_file_id
		WHERE u.u_id=$2;
	`

	assignedUser = &models.UserProfile{}
	var fileUUID, fileExt string
	row := r.db.QueryRow(ctx, query, cardID, assignedUserID)
	err = row.Scan(&assignedUser.ID, &assignedUser.Name, &assignedUser.Email, &assignedUser.JoinedAt, &assignedUser.UpdatedAt, &fileUUID, &fileExt)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return nil, fmt.Errorf("%s (query): %w", funcName, err)
	}
	assignedUser.AvatarImageURL = uploads.JoinFileURL(fileUUID, fileExt, uploads.DefaultAvatarURL)
	return assignedUser, nil
}

// DeassignUserFromCard убирает назначение пользователя
func (r *BoardRepository) DeassignUserFromCard(ctx context.Context, cardID int64, assignedUserID int64) (err error) {
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
	SELECT;
	`
	_, err = r.db.Exec(ctx, query, cardID, assignedUserID)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return fmt.Errorf("%s (query): %w", funcName, err)
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
		return nil, fmt.Errorf("%s (query): %w", funcName, err)
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
		return nil, fmt.Errorf("%s (query): %w", funcName, err)
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

	_, err = r.db.Exec(ctx, query, commentID)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return fmt.Errorf("%s (query): %w", funcName, err)
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
	newField.Title = *field.Title

	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return nil, fmt.Errorf("%s (query): %w", funcName, err)
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
		return nil, fmt.Errorf("%s (query): %w", funcName, err)
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
func (r *BoardRepository) SetCardCover(ctx context.Context, userID int64, cardID int64, fileID int64) (updatedCard *models.Card, err error) {
	funcName := "SetCardCover"
	query := `
	WITH update_cover AS (
		UPDATE "card"
		SET updated_at=CURRENT_TIMESTAMP, cover_file_id=$1 WHERE card_id = $2
	),
	update_board AS (
		UPDATE board SET updated_at=CURRENT_TIMESTAMP WHERE board_id = (
			SELECT b.board_id
			FROM card AS c
			JOIN kanban_column AS kc ON c.col_id=kc.col_id
			JOIN board AS b ON b.board_id = kc.board_id
			WHERE c.card_id=$1
		)
	),
	update_user AS (
		UPDATE user_to_board SET last_visit_at=CURRENT_TIMESTAMP
		WHERE u_id=$3 AND board_id=(
			SELECT board_id
			FROM kanban_column AS cc
			JOIN kanban_card AS c ON c.col_id=cc.col_id
			WHERE card_id=$1
		)
	)
	SELECT;
	`

	updatedCard = &models.Card{}
	row := r.db.QueryRow(ctx, query, fileID, cardID, userID)
	err = row.Scan()
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return nil, fmt.Errorf("%s (query): %w", funcName, err)
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
func (r *BoardRepository) AddAttachment(ctx context.Context, userID int64, cardID int64, fileID int64, originalName string) (newAttachment *models.Attachment, err error) {
	funcName := "AddAttachment"
	query := `
	WITH insert_attachment AS (
		INSERT INTO card_attachment (card_id, file_id, original_name, attached_by) VALUES ($1, $2, $3, $4)
		RETURNING attachment_id, original_name, file_id, created_at
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
	SELECT a.attachment_id, a.original_name, a.created_at, f.file_uuid, f.file_extension
	FROM insert_attachment AS a
	JOIN user_uploaded_file AS f ON a.file_id = f.file_id;
	`

	newAttachment = &models.Attachment{}
	var fileUUID, fileExt string

	row := r.db.QueryRow(ctx, query, cardID, fileID, originalName, userID)
	err = row.Scan(&newAttachment.ID, &newAttachment.OriginalName,
		&newAttachment.CreatedAt, &fileUUID, &fileExt)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return nil, fmt.Errorf("%s (query): %w", funcName, err)
	}

	newAttachment.FileName = uploads.JoinFileURL(fileUUID, fileExt, "unknown-error")
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
	tag, err := r.db.Exec(ctx, query, attachmentID)
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
	UPDATE user_to_board
	SET invite_link_uuid = uuid_generate_v4()
	WHERE u_id = $1 AND board_id = $2
	RETURNING invite_link_uuid;
	`

	link = &models.InviteLink{}
	row := r.db.QueryRow(ctx, query, userID, boardID)
	err = row.Scan(&link.InviteLinkUUID)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return nil, fmt.Errorf("%s (query): %w", funcName, err)
	}
	return link, nil
}

// DeleteInviteLink удаляет ссылку-приглашение
func (r *BoardRepository) DeleteInviteLink(ctx context.Context, userID int64, boardID int64) (err error) {
	funcName := "DeleteInviteLink"
	query := `
	WITH delete_invite_link AS (
		UPDATE user_to_board SET invite_link_uuid = NULL WHERE u_id = $1 AND board_id = $2
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
		SELECT b.board_id, b.name, b.created_at, b.updated_at,
		COALESCE(f.file_uuid::text, ''), COALESCE(f.file_extension, '')
		FROM user_to_board AS ub
		JOIN board AS b ON ub.board_id = b.board_id
		LEFT JOIN user_uploaded_file AS f ON f.file_id=b.background_image_id
		WHERE ub.invite_link_uuid = $1::UUID;
	`

	board = &models.Board{}
	var fileUUID, fileExt string
	row := r.db.QueryRow(ctx, query, inviteUUID)
	err = row.Scan(
		&board.ID,
		&board.Name,
		&board.CreatedAt,
		&board.UpdatedAt,
		&fileUUID,
		&fileExt,
	)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return nil, fmt.Errorf("%s (query): %w", funcName, err)
	}
	board.BackgroundImageURL = uploads.JoinFileURL(fileUUID, fileExt, uploads.DefaultBackgroundURL)
	return board, nil
}

// AcceptInvite добавляет приглашённого пользователя на доску с правами зрителя
func (r *BoardRepository) AcceptInvite(ctx context.Context, userID int64, boardID int64, inviteUUID string) (err error) {
	funcName := "AcceptInvite"
	query := `
	WITH update_board AS (
		UPDATE board SET updated_at=CURRENT_TIMESTAMP WHERE board_id = $1
	)
	INSERT INTO user_to_board (u_id, board_id, role, added_by, updated_by)
	VALUES ($2, $1, 'viewer',
		(SELECT u_id FROM user_to_board WHERE invite_link_uuid=$3),
		(SELECT u_id FROM user_to_board WHERE invite_link_uuid=$3)
	);
	`

	tag, err := r.db.Exec(ctx, query, boardID, userID, inviteUUID)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return fmt.Errorf("%s (query): %w", funcName, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s (query): no rows affected", funcName)
	}
	return nil
}
