package repository

import "RPO_back/internal/models"

// CreateColumn создаёт колонку на канбане
func (r *BoardRepository) CreateColumn(boardId int, title string) (newColumn *models.Column, err error) {
	panic("Not implemented")
}

// UpdateColumn обновляет колонку на канбане
func (r *BoardRepository) UpdateColumn(boardId int, columnId int, data models.ColumnRequest) (newColumn *models.Column, err error) {
	panic("Not implemented")
}

// DeleteColumn убирает колонку с канбана
func (r *BoardRepository) DeleteColumn(boardId int, columnId int) error {
	panic("Not implemented")
}
