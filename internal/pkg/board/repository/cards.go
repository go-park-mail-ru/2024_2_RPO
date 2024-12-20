package repository

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/utils/logging"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// GetCardsForBoard возвращает все карточки, размещённые на доске
func (r *BoardRepository) GetCardsForBoard(ctx context.Context, boardID int64) (cards []models.Card, err error) {
	funcName := "GetCardsForBoard"
	query := `
	SELECT
		c.card_id,
		c.col_id,
		c.title,
		c.created_at,
		c.updated_at,
		c.deadline,
    	c.is_done,
		(SELECT (NOT COUNT(*)=0) FROM checklist_field AS f WHERE f.card_id=c.card_id),
    	(SELECT (NOT COUNT(*)=0) FROM card_attachment AS f WHERE f.card_id=c.card_id),
    	(SELECT (NOT COUNT(*)=0 )FROM card_user_assignment AS f WHERE f.card_id=c.card_id),
    	(SELECT (NOT COUNT(*)=0) FROM card_comment AS f WHERE f.card_id=c.card_id)
	FROM card c
	JOIN kanban_column kc ON c.col_id = kc.col_id
	WHERE kc.board_id = $1;
`
	rows, err := r.db.Query(ctx, query, boardID)
	logging.Debug(ctx, funcName, " query has err: ", err)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return make([]models.Card, 0), nil
		}
		return nil, fmt.Errorf("%s (query): %w", funcName, err)
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
			&card.Deadine,
			&card.IsDone,
			&card.HasCheckList,
			&card.HasAttachments,
			&card.HasAssignedUsers,
			&card.HasComments,
		)
		if err != nil {
			return nil, fmt.Errorf("%s (scan): %w", funcName, err)
		}
		cards = append(cards, card)
	}

	return cards, nil
}

// CreateNewCard создаёт новую карточку
func (r *BoardRepository) CreateNewCard(ctx context.Context, columnID int64, title string) (newCard *models.Card, err error) {
	funcName := "CreateNewCard"
	query := `
	WITH new_card AS (
		INSERT INTO card (col_id, order_index, title)
		VALUES ($1, (SELECT COUNT(*) FROM "card" WHERE col_id=$1), $2)
		RETURNING card_id, card_uuid, col_id, title, created_at, updated_at
	), update_board AS (
		UPDATE board
		SET updated_at=CURRENT_TIMESTAMP
		WHERE board_id=(
			SELECT board_id
			FROM kanban_column AS c
			JOIN board AS b USING(board_id)
			WHERE c.col_id=$1
		)
	)
	SELECT card_id, card_uuid::text, col_id, title, created_at, updated_at FROM new_card;
	`

	newCard = &models.Card{}
	err = r.db.QueryRow(ctx, query, columnID, title).Scan(
		&newCard.ID,
		&newCard.UUID,
		&newCard.ColumnID,
		&newCard.Title,
		&newCard.CreatedAt,
		&newCard.UpdatedAt,
	)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return nil, fmt.Errorf("%s (query): %w", funcName, err)
	}

	return newCard, nil
}

// UpdateCard обновляет карточку
func (r *BoardRepository) UpdateCard(ctx context.Context, cardID int64, data models.CardPatchRequest) (updateCard *models.Card, err error) {
	funcName := "UpdateCard"
	query := `
	WITH update_card AS (
		UPDATE card
		SET
		title = COALESCE($2,title),
		deadline = COALESCE($3, deadline),
		is_done = COALESCE($4, is_done),
		updated_at = CURRENT_TIMESTAMP
		WHERE card_id=$1
		RETURNING card_id
	), update_board AS (
		UPDATE board
		SET updated_at=CURRENT_TIMESTAMP
		WHERE board_id=(
			SELECT b.board_id
			FROM card AS c
			JOIN kanban_column AS cc ON cc.col_id=c.col_id
			JOIN board AS b ON b.board_id=cc.board_id
			WHERE c.card_id=$1
		)
	)
	SELECT
		c.card_id,
		c.col_id,
		c.title,
		c.created_at,
		c.updated_at,
		c.deadline,
		c.is_done,
		(SELECT (NOT COUNT(*)=0) FROM checklist_field AS f WHERE f.card_id=c.card_id),
		(SELECT (NOT COUNT(*)=0) FROM card_attachment AS f WHERE f.card_id=c.card_id),
		(SELECT (NOT COUNT(*)=0 )FROM card_user_assignment AS f WHERE f.card_id=c.card_id),
		(SELECT (NOT COUNT(*)=0) FROM card_comment AS f WHERE f.card_id=c.card_id)
	FROM card AS c
	WHERE c.card_id=$1;
	`
	updateCard = &models.Card{}

	err = r.db.QueryRow(ctx, query, cardID, data.NewTitle, data.NewDeadline, data.IsDone).Scan(
		&updateCard.ID,
		&updateCard.ColumnID,
		&updateCard.Title,
		&updateCard.CreatedAt,
		&updateCard.UpdatedAt,
		&updateCard.Deadine,
		&updateCard.IsDone,
		&updateCard.HasCheckList,
		&updateCard.HasAttachments,
		&updateCard.HasAssignedUsers,
		&updateCard.HasComments,
	)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return nil, err
	}

	return updateCard, nil
}

// DeleteCard удаляет карточку
func (r *BoardRepository) DeleteCard(ctx context.Context, cardID int64) (err error) {
	funcName := "DeleteCard"
	query := `
		DELETE FROM card
		WHERE card.card_id = $1;
	`
	_, err = r.db.Exec(ctx, query, cardID)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return err
	}

	return nil
}
