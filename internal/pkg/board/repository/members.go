package repository

import "RPO_back/internal/models"

// GetMemberPermissions (предназначено для внутреннего использования)
// Возвращает информацию о правах участника на конкретной доске
func (r *BoardRepository) GetMemberPermissions(boardID int, memberUserID int) (member *models.MemberWithPermissions, err error) {
	panic("Not implemented")
}

// GetMembersWithPermissions получает всех участников на конкретной доске с информацией об их правах
func (r *BoardRepository) GetMembersWithPermissions(boardID int) (members []models.MemberWithPermissions, err error) {
	panic("Not implemented")
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
