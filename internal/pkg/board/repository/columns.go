package repository

import (
	"RPO_back/internal/models"
	"context"
)

// CreateColumn создаёт колонку на канбане
func (r *BoardRepository) CreateColumn(boardId int, title string) (newColumn *models.Column, err error) {
	query := `
		INSERT INTO kanban_column (board_id, title, created_at, updated_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING col_id, title;
	`

	newColumn = &models.Column{}
	if err = r.db.QueryRow(context.Background(), query, boardId, title).Scan(
		&newColumn.Id,
		&newColumn.Title,
	); err != nil {
		return nil, err
	}

	return newColumn, nil
}

// UpdateColumn обновляет колонку на канбане
func (r *BoardRepository) UpdateColumn(boardId int, columnId int, data models.ColumnRequest) (updateColumn *models.Column, err error) {
	query := `
		UPDATE kanban_column
		SET title = $1, updated_at = CURRENT_TIMESTAMP
		WHERE col_id = $2 AND board_id = $3
		RETURNING col_id, title;
	`

	updateColumn = &models.Column{}
	if err = r.db.QueryRow(context.Background(), query, data.NewTitle, boardId, columnId).Scan(
		&updateColumn.Id,
		&updateColumn.Title,
	); err != nil {
		return nil, err
	}

	return updateColumn, nil
}

// DeleteColumn убирает колонку с канбана
func (r *BoardRepository) DeleteColumn(boardId int, columnId int) (err error) {
	query := `
		DELETE FROM kanban_column
		WHERE col_id = $1 AND board_id = $2;
	`

	_, err = r.db.Exec(context.Background(), query, boardId, columnId)
	if err != nil {
		return err
	}

	return nil
}
