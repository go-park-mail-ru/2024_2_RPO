package repository

import (
	"RPO_back/internal/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

// GetCardsForBoard возвращает все карточки, размещённые на доске
func (r *BoardRepository) GetCardsForBoard(boardID int) (cards []models.Card, err error) {
	query := `
	SELECT
		c.card_id,
		c.col_id,
		c.title,
		c.order_index,
		c.created_at,
		c.updated_at,
		c.cover_file_uuid
	FROM card c
	JOIN kanban_column kc ON c.col_id = kc.col_id
	WHERE kc.board_id = $1;
`
	rows, err := r.db.Query(context.Background(), query, boardID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	defer rows.Close()

	cards = make([]models.Card, 0)

	for rows.Next() {
		var card models.Card
		if err := rows.Scan(
			&card.Id,
			&card.Title,
			&card.Description,
			&card.ColumnId,
			&card.CreatedAt,
			&card.UpdatedAt,
		); err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	return cards, nil
}

// GetColumnsForBoard возвращает все колонки, которые есть на доске
func (r *BoardRepository) GetColumnsForBoard(boardID int) (columns []models.Column, err error) {
	query := `
	SELECT
		col_id,
		title
	FROM kanban_column
	WHERE board_id = $1;
	`
	rows, err := r.db.Query(context.Background(), query, boardID)
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
			&column.Id,
			&column.Title,
		); err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}

	return columns, nil
}

// В последующих трёх функциях boardId нужен для того, чтобы пользователь не смог,
// например, удалить карточку с другой доски по причине хакерской натуры своей.
// Что-то типа дополнительного уровня защиты

// CreateNewCard создаёт новую карточку
func (r *BoardRepository) CreateNewCard(boardID int, columnID int, title string) (newCard *models.Card, err error) {
	query := `
	WITH col_check AS (
		SELECT 1
		FROM kanban_column
		WHERE col_id = $2 AND board_id = $1
	)
	INSERT INTO card (col_id, title, created_at, updated_at, order_index)
	SELECT $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, COALESCE(MAX(order_index), 0) + 1
	FROM card
	WHERE col_id = $2
	RETURNING card_id, col_id, title, created_at, updated_at, order_index;
	`

	newCard = &models.Card{}
	if err = r.db.QueryRow(context.Background(), query, boardID, columnID, title).Scan(
		&newCard.Id,
		&newCard.Title,
		&newCard.Description,
		&newCard.ColumnId,
		&newCard.CreatedAt,
		&newCard.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return newCard, nil
}

// UpdateCard обновляет карточку
func (r *BoardRepository) UpdateCard(boardID int, cardID int, data models.CardPatchRequest) (updateCard *models.Card, err error) {
	query := `
	UPDATE card
	SET
		title = COALESCE(NULLIF($1, ''), title),
		description = COALESCE(NULLIF($2, ''), description),
		updated_at = CURRENT_TIMESTAMP
	FROM kanban_column
	WHERE card.col_id = kanban_column.col_id
		AND kanban_column.board_id = $3
		AND card.card_id = $4
	RETURNING
		card.card_id, card.title, card.description, card.col_id, card.created_at, card.updated_at
	`
	updateCard = &models.Card{}

	if err := r.db.QueryRow(context.Background(), query, data.NewTitle, data.NewDescription, data.ColumnId, boardID, cardID).Scan(
		&updateCard.Id,
		&updateCard.Title,
		&updateCard.Description,
		&updateCard.ColumnId,
		&updateCard.CreatedAt,
		&updateCard.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return updateCard, nil
}

// DeleteCard удаляет карточку
func (r *BoardRepository) DeleteCard(boardID int, cardID int) (err error) {
	query := `
		DELETE FROM card
		USING kanban_column
		WHERE card.col_id = kanban_column.col_id
			AND kanban_column.board_id = $1
			AND card.card_id = $2
	`
	_, err = r.db.Exec(context.Background(), query, boardID, cardID)
	if err != nil {
		return err
	}

	return nil
}
