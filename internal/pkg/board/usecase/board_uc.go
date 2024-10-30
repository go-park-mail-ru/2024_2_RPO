package usecase

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/board/repository"
)

type BoardUsecase struct {
	boardRepository *repository.BoardRepository
}

func CreateBoardUsecase(boardRepository *repository.BoardRepository) *BoardUsecase {
	return &BoardUsecase{
		boardRepository: boardRepository,
	}
}

// CreateNewBoard создаёт новую доску и возвращает информацию о ней
func (uc *BoardUsecase) CreateNewBoard(userID int, data models.CreateBoardRequest) (newBoard *models.Board, err error) {
	panic("Not implemented")
}

// UpdateBoard обновляет информацию о доске и возвращает обновлённую информацию
func (uc *BoardUsecase) UpdateBoard(userID int, boardID int, data models.BoardPutRequest) (updatedBoard *models.Board, err error) {
	panic("Not implemented")
}

// DeleteBoard удаляет доску
func (uc *BoardUsecase) DeleteBoard(userID int, boardID int) error {
	panic("Not implemented")
}

// GetMyBoards получает все доски для пользователя
func (uc *BoardUsecase) GetMyBoards(userID int) (boards []models.Board, err error) {
	panic("Not implemented")
}

// GetMembersPermissions получает информацию о ролях всех участников доски
func (uc *BoardUsecase) GetMembersPermissions(userID int, boardID int) (data []models.MemberWithPermissions, err error) {
	panic("Not implemented")
}

// AddMember добавляет участника на доску с правами "viewer" и возвращает его права
func (uc *BoardUsecase) AddMember(userID int, boardID int, newMemberID int) (newMember *models.MemberWithPermissions, err error) {
	panic("Not implemented")
}

// UpdateMemberRole обновляет роль участника и возвращает обновлённые права
func (uc *BoardUsecase) UpdateMemberRole(userID int, boardID int, memberID int, newRole string) (updatedMember *models.MemberWithPermissions, err error) {
	panic("Not implemented")
}

// RemoveMember удаляет участника с доски
func (uc *BoardUsecase) RemoveMember(userID int, boardID int, memberID int) error {
	panic("Not implemented")
}

// GetBoardContent получает все карточки и колонки с доски, а также информацию о доске
func (uc *BoardUsecase) GetBoardContent(userID int, boardID int) (content *models.BoardContent, err error) {
	panic("Not implemented")
}

// CreateNewCard создаёт новую карточку и возвращает её
func (uc *BoardUsecase) CreateNewCard(userID int, boardID int, data *models.CardPatchRequest) (newCard *models.Card, err error) {
	panic("Not implemented")
}

// UpdateCard обновляет карточку и возвращает обновлённую версию
func (uc *BoardUsecase) UpdateCard(userID int, boardID int, cardID int, data *models.CardPatchRequest) (updatedCard *models.Card, err error) {
	panic("Not implemented")
}

// DeleteCard удаляет карточку
func (uc *BoardUsecase) DeleteCard(userID int, boardID int, cardID int) error {
	panic("Not implemented")
}

// CreateColumn создаёт колонку канбана на доске и возвращает её
func (uc *BoardUsecase) CreateColumn(userID int, boardID int, data *models.ColumnRequest) (newCol *models.Column, err error) {
	panic("Not implemented")
}

// UpdateColumn изменяет колонку и возвращает её обновлённую версию
func (uc *BoardUsecase) UpdateColumn(userID int, boardID int, columnID int, data *models.ColumnRequest) (updatedCol *models.Column, err error) {
	panic("Not implemented")
}

// DeleteColumn удаляет колонку
func (uc *BoardUsecase) DeleteColumn(userID int, boardID int, columnID int) error {
	panic("Not implemented")
}
