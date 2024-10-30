package repository

import "RPO_back/internal/models"

// GetCardsForBoard возвращает все карточки, размещённые на доске
func (r *BoardRepository) GetCardsForBoard(boardId int) (cards []models.Card, err error) {
	panic("Not implemented")
}

// GetColumnsForBoard возвращает все колонки, которые есть на доске
func (r *BoardRepository) GetColumnsForBoard(boardId int) (columns []models.Column, err error) {
	panic("Not implemented")
}

// В последующих трёх функциях boardId нужен для того, чтобы пользователь не смог,
// например, удалить карточку с другой доски по причине хакерской натуры своей.
// Что-то типа дополнительного уровня защиты

// CreateNewCard создаёт новую карточку
func (r *BoardRepository) CreateNewCard(boardId int, columnId int, title string) (newCard *models.Card, err error) {
	panic("Not implemented")
}

// UpdateCard обновляет карточку
func (r *BoardRepository) UpdateCard(boardId int, cardId int, data models.CardPatchRequest) (newCard *models.Card, err error) {
	panic("Not implemented")
}

// DeleteCard удаляет карточку
func (r *BoardRepository) DeleteCard(boardId int, cardId int) (err error) {
	panic("Not implemented")
}
