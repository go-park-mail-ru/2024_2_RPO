package usecase

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/board/repository"
	"errors"
	"fmt"
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

	info, err := uc.boardRepository.GetBoard(int64(boardID))
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
func (uc *BoardUsecase) CreateNewCard(userID int, boardID int, data *models.CardPatchRequest) (newCard *models.Card, err error) {
	_, err = uc.boardRepository.GetMemberPermissions(boardID, userID, false)
	if err != nil {
		if errors.Is(err, errs.ErrNotPermitted) {
			return nil, fmt.Errorf("CreateNewCard: %w", err)
		}
		if errors.Is(err, errs.ErrNotFound) {
			return nil, fmt.Errorf("CreateNewCard: %w", err)
		}
		return nil, fmt.Errorf("CreateNewCard (add GetMemberPermissions): %w", err)
	}

	card, err := uc.boardRepository.CreateNewCard(boardID, data.ColumnId, data.NewTitle)
	if err != nil {
		return nil, fmt.Errorf("CreateNewCard (add CreateNewCard): %w", err)
	}

	return &models.Card{
		Id:          card.Id,
		Title:       card.Title,
		Description: card.Description,
		ColumnId:    card.ColumnId,
		CreatedAt:   card.CreatedAt,
		UpdatedAt:   card.UpdatedAt,
	}, nil
}

// UpdateCard обновляет карточку и возвращает обновлённую версию
func (uc *BoardUsecase) UpdateCard(userID int, boardID int, cardID int, data *models.CardPatchRequest) (updatedCard *models.Card, err error) {
	_, err = uc.boardRepository.GetMemberPermissions(boardID, userID, false)
	if err != nil {
		if errors.Is(err, errs.ErrNotPermitted) {
			return nil, fmt.Errorf("UpdateCard: %w", err)
		}
		if errors.Is(err, errs.ErrNotFound) {
			return nil, fmt.Errorf("UpdateCard: %w", err)
		}
		return nil, fmt.Errorf("UpdateCard (add GetMemberPermissions): %w", err)
	}

	updatedCard, err = uc.boardRepository.UpdateCard(boardID, cardID, *data)
	if err != nil {
		return nil, fmt.Errorf("UpdateCard (add UpdateCard): %w", err)
	}

	return &models.Card{
		Id:          updatedCard.Id,
		Title:       updatedCard.Title,
		Description: updatedCard.Description,
		ColumnId:    updatedCard.ColumnId,
		CreatedAt:   updatedCard.CreatedAt,
		UpdatedAt:   updatedCard.UpdatedAt,
	}, nil
}

// DeleteCard удаляет карточку
func (uc *BoardUsecase) DeleteCard(userID int, boardID int, cardID int) (err error) {
	_, err = uc.boardRepository.GetMemberPermissions(boardID, userID, false)
	if err != nil {
		if errors.Is(err, errs.ErrNotPermitted) {
			return err
		}
		if errors.Is(err, errs.ErrNotFound) {
			return err
		}
		return err
	}

	err = uc.boardRepository.DeleteCard(boardID, cardID)
	if err != nil {
		return err
	}

	return nil
}

// CreateColumn создаёт колонку канбана на доске и возвращает её
func (uc *BoardUsecase) CreateColumn(userID int, boardID int, data *models.ColumnRequest) (newCol *models.Column, err error) {
	_, err = uc.boardRepository.GetMemberPermissions(boardID, userID, false)
	if err != nil {
		if errors.Is(err, errs.ErrNotPermitted) {
			return nil, fmt.Errorf("CreateColumn: %w", err)
		}
		if errors.Is(err, errs.ErrNotFound) {
			return nil, fmt.Errorf("CreateColumn: %w", err)
		}
		return nil, fmt.Errorf("CreateColumn (add GetMemberPermissions): %w", err)
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
	_, err = uc.boardRepository.GetMemberPermissions(boardID, userID, false)
	if err != nil {
		if errors.Is(err, errs.ErrNotPermitted) {
			return nil, fmt.Errorf("UpdateColumn: %w", err)
		}
		if errors.Is(err, errs.ErrNotFound) {
			return nil, fmt.Errorf("UpdateColumn: %w", err)
		}
		return nil, fmt.Errorf("UpdateColumn (add GetMemberPermissions): %w", err)
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
	_, err = uc.boardRepository.GetMemberPermissions(boardID, userID, false)
	if err != nil {
		if errors.Is(err, errs.ErrNotPermitted) {
			return err
		}
		if errors.Is(err, errs.ErrNotFound) {
			return err
		}
		return err
	}

	err = uc.boardRepository.DeleteColumn(boardID, columnID)
	if err != nil {
		return err
	}

	return nil
}
