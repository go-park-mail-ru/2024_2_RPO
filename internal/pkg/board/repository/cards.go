package repository

import (
	"RPO_back/internal/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

// GetCardsForBoard возвращает все карточки, размещённые на доске
func (r *BoardRepository) GetCardsForBoard(boardID int) (cards []models.Card, err error) {
	query := `SELECT `
	tag, err := r.db.Query(context.Background(), query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
	}
	panic("Not implemented")
}

// GetColumnsForBoard возвращает все колонки, которые есть на доске
func (r *BoardRepository) GetColumnsForBoard(boardID int) (columns []models.Column, err error) {
	query := ``
	r.db.Query(context.Background(), query)
	panic("Not implemented")
}

// В последующих трёх функциях boardId нужен для того, чтобы пользователь не смог,
// например, удалить карточку с другой доски по причине хакерской натуры своей.
// Что-то типа дополнительного уровня защиты

// CreateNewCard создаёт новую карточку
func (r *BoardRepository) CreateNewCard(boardID int, columnID int, title string) (newCard *models.Card, err error) {
	query := ``
	r.db.Query(context.Background(), query)
	panic("Not implemented")
}

// UpdateCard обновляет карточку
func (r *BoardRepository) UpdateCard(boardID int, cardID int, data models.CardPatchRequest) (newCard *models.Card, err error) {
	query := ``
	r.db.Query(context.Background(), query)
	panic("Not implemented")
}

// DeleteCard удаляет карточку
func (r *BoardRepository) DeleteCard(boardID int, cardID int) (err error) {
	query := ``
	r.db.Query(context.Background(), query)
	panic("Not implemented")
}
