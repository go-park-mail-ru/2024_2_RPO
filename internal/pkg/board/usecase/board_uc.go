package usecase

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/board/repository"
	"fmt"
)

var permissiveTable = make(map[string]int)

func init() {
	permissiveTable["viewer"] = 0
	permissiveTable["editor"] = 1
	permissiveTable["editor_chief"] = 2
	permissiveTable["admin"] = 3
}

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
	_, err = uc.boardRepository.GetMemberPermissions(boardID, userID, false)
	if err != nil {
		return nil, fmt.Errorf("GetMembersPermissions (permissions): %w", err)
	}
	permissions, err := uc.boardRepository.GetMembersWithPermissions(boardID)
	if err != nil {
		return nil, fmt.Errorf("GetMembersPermissions (query): %w", err)
	}
	return permissions, nil
}

// AddMember добавляет участника на доску с правами "viewer" и возвращает его права
func (uc *BoardUsecase) AddMember(userID int, boardID int, newMemberID int) (newMember *models.MemberWithPermissions, err error) {
	adderMember, err := uc.boardRepository.GetMemberPermissions(boardID, userID, false)
	if err != nil {
		return nil, fmt.Errorf("GetMembersPermissions (permissions): %w", err)
	}
	if (adderMember.Role != "admin") && (adderMember.Role != "editor_chief") {
		return nil, fmt.Errorf("GetMembersPermissions (permissions): %w", errs.ErrNotPermitted)
	}
	newMember, err = uc.boardRepository.AddMember(boardID, userID, newMemberID)
	if err != nil {
		return nil, fmt.Errorf("GetMembersPermissions (action): %w", err)
	}
	return newMember, nil
}

// UpdateMemberRole обновляет роль участника и возвращает обновлённые права
func (uc *BoardUsecase) UpdateMemberRole(userID int, boardID int, memberID int, newRole string) (updatedMember *models.MemberWithPermissions, err error) {
	updaterMember, err := uc.boardRepository.GetMemberPermissions(boardID, userID, false)
	if err != nil {
		return nil, fmt.Errorf("UpdateMemberRole (updater permissions): %w", err)
	}
	memberToUpdate, err := uc.boardRepository.GetMemberPermissions(boardID, userID, false)
	if err != nil {
		return nil, fmt.Errorf("UpdateMemberRole (member permissions): %w", err)
	}
	if updaterMember.Role != "admin" {
		if (updaterMember.Role != "admin") && (updaterMember.Role != "editor_chief") {
			return nil, fmt.Errorf("UpdateMemberRole (check): %w", errs.ErrNotPermitted)
		}
		if permissiveTable[updaterMember.Role] <= permissiveTable[newRole] {
			return nil, fmt.Errorf("UpdateMemberRole (check): %w", errs.ErrNotPermitted)
		}
		if permissiveTable[updaterMember.Role] <= permissiveTable[memberToUpdate.Role] {
			return nil, fmt.Errorf("UpdateMemberRole (check): %w", errs.ErrNotPermitted)
		}
	}
	updatedMember, err = uc.boardRepository.SetMemberRole(boardID, memberID, newRole)
	if err != nil {
		return nil, fmt.Errorf("UpdateMemberRole (action): %w", err)
	}
	return updatedMember, nil
}

// RemoveMember удаляет участника с доски
func (uc *BoardUsecase) RemoveMember(userID int, boardID int, memberID int) error {
	removerMember, err := uc.boardRepository.GetMemberPermissions(boardID, userID, false)
	if err != nil {
		fmt.Errorf("RemoveMember (remover permissions): %w", err)
	}
	memberToUpdate, err := uc.boardRepository.GetMemberPermissions(boardID, userID, false)
	if err != nil {
		fmt.Errorf("RemoveMember (member permissions): %w", err)
	}
	if removerMember.Role != "admin" {
		if (removerMember.Role != "admin") && (removerMember.Role != "editor_chief") {
			return fmt.Errorf("RemoveMember (check): %w", errs.ErrNotPermitted)
		}

		if permissiveTable[removerMember.Role] <= permissiveTable[memberToUpdate.Role] {
			fmt.Errorf("RemoveMember (check): %w", errs.ErrNotPermitted)
		}
	}
	err = uc.boardRepository.RemoveMember(boardID, memberID)
	if err != nil {
		fmt.Errorf("UpdateMemberRole (action): %w", err)
	}
	return nil
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
