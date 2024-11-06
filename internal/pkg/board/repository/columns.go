package repository

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// GetColumnsForBoard возвращает все колонки, которые есть на доске
func (r *BoardRepository) GetColumnsForBoard(ctx context.Context, boardID int) (columns []models.Column, err error) {
	query := `
	SELECT
		col_id,
		title
	FROM kanban_column
	WHERE board_id = $1;
	`
	rows, err := r.db.Query(ctx, query, boardID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	defer rows.Close()

	columns = make([]models.Column, 0)

	for rows.Next() {
		var column models.Column
		if err := rows.Scan(
			&column.ID,
			&column.Title,
		); err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}

	return columns, nil
}

// CreateColumn создаёт колонку на канбане
func (r *BoardRepository) CreateColumn(ctx context.Context, boardID int, title string) (newColumn *models.Column, err error) {
	query := `
		INSERT INTO kanban_column (board_id, title, created_at, updated_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING col_id, title;
	`

	newColumn = &models.Column{}
	if err = r.db.QueryRow(ctx, query, boardID, title).Scan(
		&newColumn.ID,
		&newColumn.Title,
	); err != nil {
		return nil, err
	}

	return newColumn, nil
}

// UpdateColumn обновляет колонку на канбане
func (r *BoardRepository) UpdateColumn(ctx context.Context, boardID int, columnID int, data models.ColumnRequest) (updateColumn *models.Column, err error) {
	query := `
		UPDATE kanban_column
		SET title = $1, updated_at = CURRENT_TIMESTAMP
		WHERE col_id = $2 AND board_id = $3
		RETURNING col_id, title;
	`

	updateColumn = &models.Column{}
	if err = r.db.QueryRow(ctx, query, data.NewTitle, columnID, boardID).Scan(
		&updateColumn.ID,
		&updateColumn.Title,
	); err != nil {
		return nil, err
	}

	return updateColumn, nil
}

// DeleteColumn убирает колонку с канбана
func (r *BoardRepository) DeleteColumn(ctx context.Context, boardID int, columnID int) (err error) {
	query := `
		DELETE FROM kanban_column
		WHERE col_id = $1 AND board_id = $2;
	`

	tag, err := r.db.Exec(ctx, query, columnID, boardID)
	if err != nil {
		return fmt.Errorf("DeleteColumn (query): %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("DeleteColumn (query): %w", errs.ErrNotFound)
	}
	// Лишние карточки удалятся каскадно (за счёт ограничения FOREIGN KEY)

	return nil
}
