package usecase

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/board"
	"errors"
	"fmt"
)

var roleLevels = make(map[string]int)

func init() {
	roleLevels["viewer"] = 0
	roleLevels["editor"] = 1
	roleLevels["editor_chief"] = 2
	roleLevels["admin"] = 3
}

type BoardUsecase struct {
	boardRepository board.BoardRepo
}

func CreateBoardUsecase(boardRepository board.BoardRepo) *BoardUsecase {
	return &BoardUsecase{
		boardRepository: boardRepository,
	}
}

// CreateNewBoard создаёт новую доску и возвращает информацию о ней
func (uc *BoardUsecase) CreateNewBoard(userID int, data models.CreateBoardRequest) (newBoard *models.Board, err error) {
	newBoard, err = uc.boardRepository.CreateBoard(data.Name, userID)
	if err != nil {
		return nil, err
	}
	_, err = uc.boardRepository.AddMember(newBoard.ID, userID, userID)
	if err != nil {
		return nil, err
	}
	_, err = uc.boardRepository.SetMemberRole(newBoard.ID, userID, "admin")
	if err != nil {
		return nil, err
	}
	return newBoard, nil
}

// UpdateBoard обновляет информацию о доске и возвращает обновлённую информацию
func (uc *BoardUsecase) UpdateBoard(userID int, boardID int, data models.BoardPutRequest) (updatedBoard *models.Board, err error) {
	deleterMember, err := uc.boardRepository.GetMemberPermissions(boardID, userID, false)
	if err != nil {
		return nil, fmt.Errorf("GetMembersPermissions (getting editor perm-s): %w", err)
	}
	if deleterMember.Role != "admin" && deleterMember.Role != "editor_chief" {
		return nil, fmt.Errorf("GetMembersPermissions (checking): %w", errs.ErrNotPermitted)
	}
	updatedBoard, err = uc.boardRepository.UpdateBoard(boardID, &data)
	return
}

// DeleteBoard удаляет доску
func (uc *BoardUsecase) DeleteBoard(userID int, boardID int) error {
	deleterMember, err := uc.boardRepository.GetMemberPermissions(boardID, userID, false)
	if err != nil {
		return fmt.Errorf("GetMembersPermissions (getting deleter perm-s): %w", err)
	}
	if deleterMember.Role != "admin" {
		return fmt.Errorf("GetMembersPermissions (checking): %w", errs.ErrNotPermitted)
	}
	err = uc.boardRepository.DeleteBoard(boardID)
	if err != nil {
		return fmt.Errorf("GetMembersPermissions (action): %w", err)
	}
	return nil
}

// GetMyBoards получает все доски для пользователя
func (uc *BoardUsecase) GetMyBoards(userID int) (boards []models.Board, err error) {
	boards, err = uc.boardRepository.GetBoardsForUser(userID)
	return
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
func (uc *BoardUsecase) AddMember(userID int, boardID int, addRequest *models.AddMemberRequest) (newMember *models.MemberWithPermissions, err error) {
	adderMember, err := uc.boardRepository.GetMemberPermissions(boardID, userID, false)
	if err != nil {
		return nil, fmt.Errorf("GetMembersPermissions (permissions): %w", err)
	}
	if (adderMember.Role != "admin") && (adderMember.Role != "editor_chief") {
		return nil, fmt.Errorf("GetMembersPermissions (permissions): %w", errs.ErrNotPermitted)
	}
	newMemberProfile, err := uc.boardRepository.GetUserByNickname(addRequest.MemberNickname)
	if err != nil {
		return nil, fmt.Errorf("GetMembersPermissions (get new user ID): %w", err)
	}
	newMember, err = uc.boardRepository.AddMember(boardID, userID, newMemberProfile.ID)
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
	memberToUpdate, err := uc.boardRepository.GetMemberPermissions(boardID, memberID, false)
	if err != nil {
		return nil, fmt.Errorf("UpdateMemberRole (member permissions): %w", err)
	}
	if updaterMember.Role != "admin" {
		if (updaterMember.Role != "admin") && (updaterMember.Role != "editor_chief") {
			return nil, fmt.Errorf("UpdateMemberRole (check1): %w", errs.ErrNotPermitted)
		}
		if roleLevels[updaterMember.Role] <= roleLevels[newRole] {
			return nil, fmt.Errorf("UpdateMemberRole (check2): %w", errs.ErrNotPermitted)
		}
		if roleLevels[updaterMember.Role] <= roleLevels[memberToUpdate.Role] {
			return nil, fmt.Errorf("UpdateMemberRole (check3): %w", errs.ErrNotPermitted)
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
		return fmt.Errorf("RemoveMember (remover permissions): %w", err)
	}
	memberToRemove, err := uc.boardRepository.GetMemberPermissions(boardID, memberID, false)
	if err != nil {
		return fmt.Errorf("RemoveMember (member permissions): %w", err)
	}
	fmt.Printf("%s removes %s\n", removerMember.Role, memberToRemove.Role)
	if removerMember.Role != "admin" && userID != memberID {
		if (removerMember.Role != "admin") && (removerMember.Role != "editor_chief") {
			return fmt.Errorf("RemoveMember (check1): %w", errs.ErrNotPermitted)
		}

		if roleLevels[removerMember.Role] <= roleLevels[memberToRemove.Role] {
			return fmt.Errorf("RemoveMember (check2): %w", errs.ErrNotPermitted)
		}
	}
	err = uc.boardRepository.RemoveMember(boardID, memberID)
	if err != nil {
		return fmt.Errorf("UpdateMemberRole (action): %w", err)
	}
	return nil
}

// GetBoardContent получает все карточки и колонки с доски, а также информацию о доске
func (uc *BoardUsecase) GetBoardContent(userID int, boardID int) (content *models.BoardContent, err error) {
	userPermissions, err := uc.boardRepository.GetMemberPermissions(boardID, userID, false)
	if err != nil {
		if errors.Is(err, errs.ErrNotPermitted) {
			return nil, fmt.Errorf("GetBoardContent: %w", errs.ErrNotPermitted)
		}
		if errors.Is(err, errs.ErrNotFound) {
			return nil, fmt.Errorf("GetBoardContent: %w", errs.ErrNotFound)
		}
		return nil, fmt.Errorf("GetBoardContent (add GetMemberPermissions): %w", err)
	}

	cards, err := uc.boardRepository.GetCardsForBoard(boardID)
	if err != nil {
		return nil, fmt.Errorf("GetBoardContent (add GetCardsForBoard): %w", err)
	}

	cols, err := uc.boardRepository.GetColumnsForBoard(boardID)
	if err != nil {
		return nil, fmt.Errorf("GetBoardContent (add GetColumnsForBoard): %w", err)
	}

	info, err := uc.boardRepository.GetBoard(boardID)
	if err != nil {
		return nil, fmt.Errorf("GetBoardContent (add GetBoard): %w", err)
	}

	return &models.BoardContent{
		Cards:     cards,
		Columns:   cols,
		BoardInfo: info,
		MyRole:    userPermissions.Role,
	}, nil
}

// CreateNewCard создаёт новую карточку и возвращает её
func (uc *BoardUsecase) CreateNewCard(userID int, boardID int, data *models.CardPutRequest) (newCard *models.Card, err error) {
	perms, err := uc.boardRepository.GetMemberPermissions(boardID, userID, false)
	if err != nil {
		if errors.Is(err, errs.ErrNotPermitted) {
			return nil, fmt.Errorf("CreateNewCard (get permissions): %w", err)
		}
		if errors.Is(err, errs.ErrNotFound) {
			return nil, fmt.Errorf("CreateNewCard (get permissions): %w", err)
		}
		return nil, fmt.Errorf("CreateNewCard (add GetMemberPermissions): %w", err)
	}
	if perms.Role == "viewer" {
		return nil, fmt.Errorf("CreateNewCard (check): %w", errs.ErrNotPermitted)
	}

	card, err := uc.boardRepository.CreateNewCard(boardID, data.NewColumnId, data.NewTitle)
	if err != nil {
		return nil, fmt.Errorf("CreateNewCard (add CreateNewCard): %w", err)
	}

	return &models.Card{
		ID:        card.ID,
		Title:     card.Title,
		ColumnID:  card.ColumnID,
		CreatedAt: card.CreatedAt,
		UpdatedAt: card.UpdatedAt,
	}, nil
}

// UpdateCard обновляет карточку и возвращает обновлённую версию
func (uc *BoardUsecase) UpdateCard(userID int, boardID int, cardID int, data *models.CardPutRequest) (updatedCard *models.Card, err error) {
	perms, err := uc.boardRepository.GetMemberPermissions(boardID, userID, false)
	if err != nil {
		if errors.Is(err, errs.ErrNotPermitted) {
			return nil, fmt.Errorf("UpdateCard: %w", err)
		}
		if errors.Is(err, errs.ErrNotFound) {
			return nil, fmt.Errorf("UpdateCard: %w", err)
		}
		return nil, fmt.Errorf("UpdateCard (add GetMemberPermissions): %w", err)
	}
	if perms.Role == "viewer" {
		return nil, fmt.Errorf("UpdateCard (check): %w", errs.ErrNotPermitted)
	}

	updatedCard, err = uc.boardRepository.UpdateCard(boardID, cardID, *data)
	if err != nil {
		return nil, fmt.Errorf("UpdateCard (repo): %w", err)
	}

	return &models.Card{
		ID:        updatedCard.ID,
		Title:     updatedCard.Title,
		ColumnID:  updatedCard.ColumnID,
		CreatedAt: updatedCard.CreatedAt,
		UpdatedAt: updatedCard.UpdatedAt,
	}, nil
}

// DeleteCard удаляет карточку
func (uc *BoardUsecase) DeleteCard(userID int, boardID int, cardID int) (err error) {
	perms, err := uc.boardRepository.GetMemberPermissions(boardID, userID, false)
	if err != nil {
		if errors.Is(err, errs.ErrNotPermitted) {
			return err
		}
		if errors.Is(err, errs.ErrNotFound) {
			return err
		}
		return err
	}
	if perms.Role == "viewer" {
		return fmt.Errorf("DeleteCard (check): %w", errs.ErrNotPermitted)
	}

	err = uc.boardRepository.DeleteCard(boardID, cardID)
	if err != nil {
		return err
	}

	return nil
}

// CreateColumn создаёт колонку канбана на доске и возвращает её
func (uc *BoardUsecase) CreateColumn(userID int, boardID int, data *models.ColumnRequest) (newCol *models.Column, err error) {
	perms, err := uc.boardRepository.GetMemberPermissions(boardID, userID, false)
	if err != nil {
		if errors.Is(err, errs.ErrNotPermitted) {
			return nil, fmt.Errorf("CreateColumn: %w", err)
		}
		if errors.Is(err, errs.ErrNotFound) {
			return nil, fmt.Errorf("CreateColumn: %w", err)
		}
		return nil, fmt.Errorf("CreateColumn (add GetMemberPermissions): %w", err)
	}
	if perms.Role == "viewer" {
		return nil, fmt.Errorf("CreateColumn (check): %w", errs.ErrNotPermitted)
	}

	column, err := uc.boardRepository.CreateColumn(boardID, data.NewTitle)
	if err != nil {
		return nil, fmt.Errorf("CreateColumn (add CreateColumn): %w", err)
	}

	return &models.Column{
		Id:    column.Id,
		Title: column.Title,
	}, nil
}

// UpdateColumn изменяет колонку и возвращает её обновлённую версию
func (uc *BoardUsecase) UpdateColumn(userID int, boardID int, columnID int, data *models.ColumnRequest) (updatedCol *models.Column, err error) {
	perms, err := uc.boardRepository.GetMemberPermissions(boardID, userID, false)
	if err != nil {
		if errors.Is(err, errs.ErrNotPermitted) {
			return nil, fmt.Errorf("UpdateColumn: %w", err)
		}
		if errors.Is(err, errs.ErrNotFound) {
			return nil, fmt.Errorf("UpdateColumn: %w", err)
		}
		return nil, fmt.Errorf("UpdateColumn (add GetMemberPermissions): %w", err)
	}
	if perms.Role == "viewer" {
		return nil, fmt.Errorf("UpdateColumn (check): %w", errs.ErrNotPermitted)
	}

	updatedCol, err = uc.boardRepository.UpdateColumn(boardID, columnID, *data)
	if err != nil {
		return nil, fmt.Errorf("UpdateColumn (add UpdateColumn): %w", err)
	}

	return &models.Column{
		Id:    updatedCol.Id,
		Title: updatedCol.Title,
	}, nil
}

// DeleteColumn удаляет колонку
func (uc *BoardUsecase) DeleteColumn(userID int, boardID int, columnID int) (err error) {
	perms, err := uc.boardRepository.GetMemberPermissions(boardID, userID, false)
	if err != nil {
		if errors.Is(err, errs.ErrNotPermitted) {
			return err
		}
		if errors.Is(err, errs.ErrNotFound) {
			return err
		}
		return err
	}
	if perms.Role == "viewer" {
		return fmt.Errorf("DeleteColumn (check): %w", errs.ErrNotPermitted)
	}

	err = uc.boardRepository.DeleteColumn(boardID, columnID)
	if err != nil {
		return err
	}

	return nil
}

// ReplaceBackground заменяет задний фон доски и возвращает обновлённую доску
func (uc *BoardUsecase) ReplaceBackground(userID int, originalFileName string) (updatedBoard models.Board, fileName string, err error) {
	panic("Not implemented")
}
