package repository

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/utils/logging"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// GetColumnsForBoard возвращает все колонки, которые есть на доске
func (r *BoardRepository) GetColumnsForBoard(ctx context.Context, boardID int64) (columns []models.Column, err error) {
	funcName := "GetColumnsForBoard"
	query := `
	SELECT
		col_id,
		title
	FROM kanban_column
	WHERE board_id = $1;
	`
	rows, err := r.db.Query(ctx, query, boardID)
	logging.Debug(ctx, funcName, " query has err: ", err)
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
func (r *BoardRepository) CreateColumn(ctx context.Context, boardID int64, title string) (newColumn *models.Column, err error) {
	funcName := "CreateColumn"
	query := `
		INSERT INTO kanban_column (board_id, title, created_at, updated_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING col_id, title;
	`

	newColumn = &models.Column{}
	err = r.db.QueryRow(ctx, query, boardID, title).Scan(
		&newColumn.ID,
		&newColumn.Title,
	)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return nil, err
	}

	return newColumn, nil
}

// UpdateColumn обновляет колонку на канбане
func (r *BoardRepository) UpdateColumn(ctx context.Context, columnID int64, data models.ColumnRequest) (updateColumn *models.Column, err error) {
	funcName := "UpdateColumn"
	query := `
		UPDATE kanban_column
		SET title = $1, updated_at = CURRENT_TIMESTAMP
		WHERE col_id = $2
		RETURNING col_id, title;
	`

	updateColumn = &models.Column{}
	err = r.db.QueryRow(ctx, query, data.NewTitle, columnID).Scan(
		&updateColumn.ID,
		&updateColumn.Title,
	)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return nil, fmt.Errorf("%s (query): %w", funcName, err)
	}

	return updateColumn, nil
}

// DeleteColumn убирает колонку с канбана
func (r *BoardRepository) DeleteColumn(ctx context.Context, columnID int64) (err error) {
	funcName := "DeleteColumn"
	query := `
		DELETE FROM kanban_column
		WHERE col_id = $1;
	`

	tag, err := r.db.Exec(ctx, query, columnID)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return fmt.Errorf("%s (query): %w", funcName, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s (query): %w", funcName, errs.ErrNotFound)
	}
	// Лишние карточки удалятся каскадно (за счёт ограничения FOREIGN KEY)

	return nil
}
