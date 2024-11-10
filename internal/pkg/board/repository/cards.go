package repository

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/utils/logging"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

// GetCardsForBoard возвращает все карточки, размещённые на доске
func (r *BoardRepository) GetCardsForBoard(ctx context.Context, boardID int) (cards []models.Card, err error) {
	query := `
	SELECT
		c.card_id,
		c.col_id,
		c.title,
		c.created_at,
		c.updated_at
	FROM card c
	JOIN kanban_column kc ON c.col_id = kc.col_id
	WHERE kc.board_id = $1;
`
	rows, err := r.db.Query(ctx, query, boardID)
	logging.Debug(ctx, "GetCardsForBoard query has err: ", err)

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
		err := rows.Scan(
			&card.ID,
			&card.ColumnID,
			&card.Title,
			&card.CreatedAt,
			&card.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	return cards, nil
}

// В последующих трёх функциях boardId нужен для того, чтобы пользователь не смог,
// например, удалить карточку с другой доски по причине хакерской натуры своей.
// Что-то типа дополнительного уровня защиты

// CreateNewCard создаёт новую карточку
func (r *BoardRepository) CreateNewCard(ctx context.Context, boardID int, columnID int, title string) (newCard *models.Card, err error) {
	query := `
	WITH col_check AS (
		SELECT 1
		FROM kanban_column
		WHERE col_id = $2 AND board_id = $1
	)
	INSERT INTO card (col_id, title, created_at, updated_at, order_index)
	VALUES ($2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 0) -- Временное решение, TODO для drag-n-drop сделать подзапрос
	RETURNING card_id, col_id, title, created_at, updated_at;
	`

	newCard = &models.Card{}
	err = r.db.QueryRow(ctx, query, boardID, columnID, title).Scan(
		&newCard.ID,
		&newCard.ColumnID,
		&newCard.Title,
		&newCard.CreatedAt,
		&newCard.UpdatedAt,
	)
	logging.Debug(ctx, "CreateNewCard query has err: ", err)
	if err != nil {
		return nil, err
	}

	return newCard, nil
}

// UpdateCard обновляет карточку
func (r *BoardRepository) UpdateCard(ctx context.Context, boardID int, cardID int, data models.CardPutRequest) (updateCard *models.Card, err error) {
	query := `
	UPDATE card
	SET
		title = $1,
		col_id = $2,
		updated_at = CURRENT_TIMESTAMP
	FROM kanban_column
	WHERE card.col_id = kanban_column.col_id
		AND kanban_column.board_id = $3
		AND card.card_id = $4
	RETURNING
		card.card_id, card.title, card.col_id, card.created_at, card.updated_at
	`
	updateCard = &models.Card{}

	err = r.db.QueryRow(ctx, query, data.NewTitle, data.NewColumnID, boardID, cardID).Scan(
		&updateCard.ID,
		&updateCard.Title,
		&updateCard.ColumnID,
		&updateCard.CreatedAt,
		&updateCard.UpdatedAt,
	)
	logging.Debug(ctx, "UpdateCard query has err: ", err)
	if err != nil {
		return nil, err
	}

	return updateCard, nil
}

// DeleteCard удаляет карточку
func (r *BoardRepository) DeleteCard(ctx context.Context, boardID int, cardID int) (err error) {
	query := `
		DELETE FROM card
		USING kanban_column
		WHERE card.col_id = kanban_column.col_id
			AND kanban_column.board_id = $1
			AND card.card_id = $2
	`
	_, err = r.db.Exec(ctx, query, boardID, cardID)
	logging.Debug(ctx, "DeleteCard query has err: ", err)
	if err != nil {
		return err
	}

	return nil
}
