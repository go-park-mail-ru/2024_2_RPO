package repository

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// GetUserProfile получает из базы профиль пользователя
func (r *BoardRepository) GetUserProfile(userID int) (user *models.UserProfile, err error) {
	query := `
	SELECT u_id, nickname, email, description, joined_at, updated_at
	FROM "user"
	WHERE u_id=$1;
	`
	rows := r.db.QueryRow(context.Background(), query, userID)
	user = &models.UserProfile{}
	err = rows.Scan(&user.Id, &user.Name, &user.Email,
		&user.Description,
		&user.JoinedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetUserProfile: %w", errs.ErrNotFound)
		}
		return nil, fmt.Errorf("GetUserProfile: %w", err)
	}
	return user, nil
}

// GetMemberPermissions (предназначено для внутреннего использования)
// Возвращает информацию о правах участника на конкретной доске
//
// getAdderInfo - если равен true, запрос получит профили пригласившего
// пользователя и пользователя, внёсшего последнее изменение. false -
// поля AddedBy и UpdatedBy будут установлены в nil. Но если
// getAdderInfo равен true, ещё не факт, что указанные поля
// будут не nil
func (r *BoardRepository) GetMemberPermissions(boardID int, memberUserID int, getAdderInfo bool) (member *models.MemberWithPermissions, err error) {
	query := `
	SELECT
	ub.role
	ub.added_at,
	ub.updated_at,
	COALESCE(ub.added_by, -1)
	COALESCE(ub.updated_by, -1)
	FROM user_to_board AS ub
	WHERE ub.u_id=$1
	AND ub.board_id=$2;
	`
	member = &models.MemberWithPermissions{}
	// Получение профиля пользователя
	userProfile, err := r.GetUserProfile(memberUserID)
	if err != nil {
		return nil, fmt.Errorf("GetMemberPermissions (getting user profile): %w", err)
	}
	// Проверка на то, что доска существует
	_, err = r.GetBoard(int64(boardID))
	if err != nil {
		return nil, fmt.Errorf("GetMemberPermissions (getting board): %w", err)
	}
	member.User = userProfile
	var addedByID, updatedByID int
	rows := r.db.QueryRow(context.Background(), query, memberUserID, boardID)
	err = rows.Scan(&member.Role,
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
	if getAdderInfo == true {
		if addedByID != -1 {
			adder, err := r.GetUserProfile(addedByID)
			if err != nil {
				return nil, fmt.Errorf("GetMemberPermissions (getting adder profile): %w", err)
			}
			member.AddedBy = adder
		}
		if updatedByID != -1 {
			updater, err := r.GetUserProfile(updatedByID)
			if err != nil {
				return nil, fmt.Errorf("GetMemberPermissions (getting last updater profile): %w", err)
			}
			member.UpdatedBy = updater
		}
	}
	return member, nil
}

// GetMembersWithPermissions получает всех участников на конкретной
// доске с информацией об их правах и с разрешением профилей добавителя
// и пользователя, внёсшего последнее обновление в роль
func (r *BoardRepository) GetMembersWithPermissions(boardID int) (members []models.MemberWithPermissions, err error) {
	query := `
	SELECT

	m.u_id, m.nickname,
	m.email, m.description,
	m.joined_at, m.updated_at,

	ub.role, ub.added_at, ub.updated_at,

	COALESCE(adder.u_id,-1), adder.nickname, adder.email,
	adder.description, adder.joined_at, adder.updated_at,

	COALESCE(updater.u_id,-1), updater.nickname, updater.email,
	updater.description, updater.joined_at, updater.updated_at

	FROM "user" AS m
	JOIN user_to_board AS ub ON m.u_id=ub.u_id
	LEFT JOIN "user" AS adder ON adder.u_id=ub.added_by
	LEFT JOIN "user" AS updater ON updater.u_id=ub.updated_by
	WHERE ub.b_id=$1;
	`
	_, err = r.GetBoard(int64(boardID))
	if err != nil {
		return nil, fmt.Errorf("GetMembersWithPermissions (getting board): %w", errs.ErrNotFound)
	}
	rows, err := r.db.Query(context.Background(), query, boardID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("GetMembersWithPermissions (main query): %w", err)
	}
	for rows.Next() {
		err := rows.Scan()
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("GetMembersWithPermissions: %w", errs.ErrNotFound)
			}
			return nil, fmt.Errorf("GetMembersWithPermissions: %w", err)
		}
	}

}

// SetMemberRole устанавливает существующему участнику права (роль)
func (r *BoardRepository) SetMemberRole(boardID int, memberUserID int, newRole string) (member *models.MemberWithPermissions, err error) {
	panic("Not implemented")
}

// RemoveMember удаляет участника с доски
func (r *BoardRepository) RemoveMember(boardID int, memberUserID int) (err error) {
	panic("Not implemented")
}

// AddMember добавляет участника на доску с правами "viewer"
func (r *BoardRepository) AddMember(boardID int, memberUserID int) (member *models.MemberWithPermissions, err error) {
	panic("Not implemented")
}
