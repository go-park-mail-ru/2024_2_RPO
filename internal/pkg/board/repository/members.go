package repository

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/utils/logging"
	"RPO_back/internal/pkg/utils/uploads"
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/jackc/pgx/v5"
)

// GetUserProfile получает из базы профиль пользователя
func (r *BoardRepository) GetUserProfile(ctx context.Context, userID int) (user *models.UserProfile, err error) {
	query := `
	SELECT
	u_id,
	nickname,
	email,
	description,
	joined_at,
	updated_at,
	COALESCE(f.file_uuid::text, ''),
	COALESCE(f.file_extension, '')
	FROM "user" AS u
	LEFT JOIN user_uploaded_file AS f ON f.file_uuid=u.avatar_file_uuid
	WHERE u_id=$1;
	`
	rows := r.db.QueryRow(ctx, query, userID)
	user = &models.UserProfile{}
	var avatarUUID, avatarExt string
	err = rows.Scan(&user.ID, &user.Name, &user.Email,
		&user.Description,
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
func (r *BoardRepository) GetMemberPermissions(ctx context.Context, boardID int, memberUserID int, verbose bool) (member *models.MemberWithPermissions, err error) {
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
	var addedByID, updatedByID int
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
	if verbose == true {
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
func (r *BoardRepository) GetMembersWithPermissions(ctx context.Context, boardID int, userID int) (members []models.MemberWithPermissions, err error) {
	query := `
	SELECT

	m.u_id, m.nickname,
	m.email, m.description,
	m.joined_at, m.updated_at,

	ub.role, ub.added_at, ub.updated_at,

	COALESCE(adder.u_id,-1), adder.nickname, adder.email,
	adder.description, adder.joined_at, adder.updated_at,

	COALESCE(updater.u_id,-1), updater.nickname, updater.email,
	updater.description, updater.joined_at, updater.updated_at,

	COALESCE(f_m.file_uuid::text,''), COALESCE(f_m.file_extension,''),
	COALESCE(f_adder.file_uuid::text,''), COALESCE(f_adder.file_extension,''),
	COALESCE(f_updater.file_uuid::text,''), COALESCE(f_updater.file_extension,'')

	FROM "user" AS m
	JOIN user_to_board AS ub ON m.u_id=ub.u_id
	LEFT JOIN "user" AS adder ON adder.u_id=ub.added_by
	LEFT JOIN "user" AS updater ON updater.u_id=ub.updated_by

	LEFT JOIN user_uploaded_file AS f_m ON f_m.file_uuid=m.avatar_file_uuid
	LEFT JOIN user_uploaded_file AS f_adder ON f_adder.file_uuid=adder.avatar_file_uuid
	LEFT JOIN user_uploaded_file AS f_updater ON f_updater.file_uuid=updater.avatar_file_uuid
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
			&field.User.Description,
			&field.User.JoinedAt,
			&field.User.UpdatedAt,

			&field.Role, &field.AddedAt, &field.UpdatedAt,

			&field.AddedBy.ID,
			&field.AddedBy.Name,
			&field.AddedBy.Email,
			&field.AddedBy.Description,
			&field.AddedBy.JoinedAt,
			&field.AddedBy.UpdatedAt,

			&field.UpdatedBy.ID,
			&field.UpdatedBy.Name,
			&field.UpdatedBy.Email,
			&field.UpdatedBy.Description,
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
func (r *BoardRepository) SetMemberRole(ctx context.Context, boardID int, memberUserID int, newRole string) (member *models.MemberWithPermissions, err error) {
	query := `
	UPDATE user_to_board
	SET role='%s',
	updated_at=CURRENT_TIMESTAMP
	WHERE u_id=$1
	AND board_id=$2;
	`

	// Дополнительная проверка для защиты от SQL-инъекций
	if !slices.Contains([]string{"viewer", "editor", "editor_chief", "admin"}, newRole) {
		return nil, fmt.Errorf("Unknown role: %s", newRole)
	}
	query = fmt.Sprintf(query, newRole)

	tag, err := r.db.Exec(ctx, query, memberUserID, boardID)
	logging.Debug(ctx, "SetMemberRole query has err: ", err)
	if tag.RowsAffected() == 0 {
		return nil, fmt.Errorf("SetMemberRole (update): %w", errs.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("SetMemberRole (update): %w", err)
	}
	member, err = r.GetMemberPermissions(ctx, boardID, memberUserID, true)
	if err != nil {
		return nil, fmt.Errorf("SetMemberRole (get updated perms): %w", err)
	}
	return member, nil
}

// RemoveMember удаляет участника с доски
func (r *BoardRepository) RemoveMember(ctx context.Context, boardID int, memberUserID int) (err error) {
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
func (r *BoardRepository) AddMember(ctx context.Context, boardID int, adderID int, memberUserID int) (member *models.MemberWithPermissions, err error) {
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
	query := `SELECT u_id, nickname, email, description, joined_at, updated_at FROM "user"
	WHERE nickname=$1;`
	user = &models.UserProfile{}
	err = r.db.QueryRow(ctx, query, nickname).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Description,
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
func (r *BoardRepository) GetMemberFromCard(ctx context.Context, userID int, cardID int64) (role string, boardID int64, err error) {
	query := `
	SELECT
	FROM card AS c
	LEFT JOIN kanban_column AS col ON col.col_id=c.col_id
	LEFT JOIN board AS b ON b.board_id=col.board_id
	LEFT JOIN user_to_board AS ub ON ub.board_id=b.board_id
	WHERE c.card_id=$1 AND ub.u_id=$2;
	`
}

// GetMemberFromCheckListField получает права пользователя из ID поля чеклиста
func (r *BoardRepository) GetMemberFromCheckListField(ctx context.Context, userID int64, fieldID int64) (role string, boardID int64, cardID int64, err error) {
	panic("not implemented")
}

// GetMemberFromAttachment получает права пользователя из ID вложения
func (r *BoardRepository) GetMemberFromAttachment(ctx context.Context, userID int64, attachmentID int64) (role string, boardID int64, cardID int64, err error) {
	panic("not implemented")
}

// GetMemberFromColumn получает права пользователя из ID колонки
func (r *BoardRepository) GetMemberFromColumn(ctx context.Context, userID int64, columnID int64) (role string, boardID int64, err error) {
	panic("not implemented")
}

// GetMemberFromComment получает права пользователя из ID комментария
func (r *BoardRepository) GetMemberFromComment(ctx context.Context, userID int64, commentID int64) (role string, boardID int64, cardID int64, err error) {
	panic("not implemented")
}

// GetCardCheckList получает чеклист для карточки
func (r *BoardRepository) GetCardCheckList(ctx context.Context, cardID int64) (err error) {
	panic("not implemented")
}

// GetCardAssignedUsers получает пользователей, назначенных на карточку
func (r *BoardRepository) GetCardAssignedUsers(ctx context.Context, cardID int64) (assignedUsers []models.UserProfile, err error) {
	panic("not implemented")
}

// GetCardComments получает комментарии, оставленные на карточку
func (r *BoardRepository) GetCardComments(ctx context.Context, cardID int64) (comments []models.Comment, err error) {
	panic("not implemented")
}

// GetCardAttachments получает вложения к карточке
func (r *BoardRepository) GetCardAttachments(ctx context.Context, cardID int64) (attachments []models.Attachment, err error) {
	panic("not implemented")
}

// GetCardsForMove получает списки карточек на двух колонках.
// Нужно для Drag-n-Drop (колонки откуда перемещаем и куда)
func (r *BoardRepository) GetCardsForMove(ctx context.Context, colID int64) (column []models.Card, err error) {
	panic("not implemented")
}

// GetColumnsForMove получает список всех колонок, чтобы сделать Drag-n-Drop
func (r *BoardRepository) GetColumnsForMove(ctx context.Context, boardID int64) (columns []models.Column, err error) {
	panic("not implemented")
}

// RearrangeCards обновляет позиции всех карточек колонки, чтобы сделать порядок, как в слайсе
func (r *BoardRepository) RearrangeCards(ctx context.Context, columnID int64, cards []models.Card) (err error) {
	panic("not implemented")
}

// RearrangeColumns обновляет позиции всех колонок, чтобы сделать порядок, как в слайсе
func (r *BoardRepository) RearrangeColumns(ctx context.Context, columns []models.Column) (err error) {
	panic("not implemented")
}

// AssignUserToCard назначает пользователя на карточку
func (r *BoardRepository) AssignUserToCard(ctx context.Context, cardID int64, assignedUserID int64) (err error) {
	panic("not implemented")
}

// DeassignUserFromCard убирает назначение пользователя
func (r *BoardRepository) DeassignUserFromCard(ctx context.Context, cardID int64, assignedUserID int64) (err error) {
	panic("not implemented")
}

// CreateComment добавляет на карточку комментарий
func (r *BoardRepository) CreateComment(ctx context.Context, userID int64, cardID int64, comment *models.CommentRequest) (newComment *models.Comment, err error) {
	panic("not implemented")
}

// UpdateComment редактирует комментарий
func (r *BoardRepository) UpdateComment(ctx context.Context, commentID int64, update *models.CommentRequest) (updatedComment *models.Comment, err error) {
	panic("not implemented")
}

// DeleteComment удаляет комментарий
func (r *BoardRepository) DeleteComment(ctx context.Context, commentID int64) (err error) {
	panic("not implemented")
}

// CreateCheckListField создаёт поле чеклиста и добавляет его в конец
func (r *BoardRepository) CreateCheckListField(ctx context.Context, cardID int64, field *models.CheckListFieldPostRequest) (err error) {
	panic("not implemented")
}

// UpdateCheckListField обновляет одно поле чеклиста
func (r *BoardRepository) UpdateCheckListField(ctx context.Context, fieldID int64, update *models.CheckListFieldPatchRequest) (updatedField *models.CheckListField, err error) {
	panic("not implemented")
}

// UpdateCheckList устанавливает порядок полей чеклиста как в слайсе
func (r *BoardRepository) ReorderCheckList(ctx context.Context, fields []models.CheckListField) (err error) {
	panic("not implemented")
}

// SetCardCover устанавливает файл обложки карточки
func (r *BoardRepository) SetCardCover(ctx context.Context, userID int64, cardID int64, originalName string, fileID int64) (updatedCard *models.Card, err error) {
	panic("not implemented")
}

// RemoveCardCover удаляет обложку карточки
func (r *BoardRepository) RemoveCardCover(ctx context.Context, cardID int64) (err error) {
	panic("not implemented")
}

// AddAttachment добавляет файл вложения в карточку
func (r *BoardRepository) AddAttachment(ctx context.Context, userID int64, cardID int64, originalName string, fileID int64) (newAttachment *models.Attachment, err error) {
	panic("not implemented")
}

// RemoveAttachment удаляет вложение
func (r *BoardRepository) RemoveAttachment(ctx context.Context, attachmentID int64) (err error) {
	panic("not implemented")
}

// PullInviteLink заменяет для доски индивидуальную ссылку-приглашение и возвращает новую ссылку
func (r *BoardRepository) PullInviteLink(ctx context.Context, userID int64, boardID int64) (link *models.InviteLink, err error) {
	panic("not implemented")
}

// DeleteInviteLink удаляет ссылку-приглашение
func (r *BoardRepository) DeleteInviteLink(ctx context.Context, userID int64, boardID int64) (err error) {
	panic("not implemented")
}

// FetchInvite возвращает информацию о доске, куда чела пригласили
func (r *BoardRepository) FetchInvite(ctx context.Context, inviteUUID string) (board *models.Board, err error) {
	panic("not implemented")
}

// AcceptInvite добавляет приглашённого пользователя на доску с правами зрителя
func (r *BoardRepository) AcceptInvite(ctx context.Context, userID int64, boardID int64, invitedUserID int64, inviteUUID string) (board *models.Board, err error) {
	panic("not implemented")
}
