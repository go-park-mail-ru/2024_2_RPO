package repository

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/utils/logging"
	"context"
	"fmt"
)

// CreateNewTag создает тег на доске
func (r *BoardRepository) CreateNewTag(ctx context.Context, boardID int64, data *models.TagPatchRequest) (newTag *models.Tag, err error) {
	funcName := "CreateNewTag"
	query := `
		WITH new_tag AS (
			INSERT INTO tag (board_id, title, color)
			VALUES ($1, $2, $3)
			RETURNING tag_id, color, title
		), update_board AS (
			UPDATE board
			SET updated_at=CURRENT_TIMESTAMP
			WHERE board_id=$1
		)
	SELECT tag_id, color, title FROM new_tag;
	`

	newTag = &models.Tag{}
	err = r.db.QueryRow(ctx, query, boardID, data.Text, data.Color).Scan(
		&newTag.ID,
		&newTag.Color,
		&newTag.Text,
	)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return nil, fmt.Errorf("%s (query): %w", funcName, err)
	}

	return newTag, nil
}

// UpdateTag обновляет тег
func (r *BoardRepository) UpdateTag(ctx context.Context, tagID int64, data *models.TagPatchRequest) (updatedTag *models.Tag, err error) {
	funcName := "UpdateTag"
	query := `
		WITH update_tag AS (
			UPDATE tag
			SET
			title=$2,
			color=$3
			WHERE tag_id=$1
			RETURNING tag_id, color, title
		), update_board AS (
			UPDATE board
			SET updated_at=CURRENT_TIMESTAMP
			WHERE board_id=(
				SELECT t.board_id
				FROM tag AS t
				WHERE t.tag_id=$1
			)
		)
		SELECT 
			t.tag_id, t.color, t.title
		FROM update_tag AS t;
	`

	updatedTag = &models.Tag{}

	err = r.db.QueryRow(ctx, query, tagID, data.Text, data.Color).Scan(
		&updatedTag.ID,
		&updatedTag.Text,
		&updatedTag.Color,
	)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return nil, err
	}

	return updatedTag, nil
}

// DeleteTag удаляет тег
func (r *BoardRepository) DeleteTag(ctx context.Context, tagID int64) (err error) {
	funcName := "DeleteTag"
	query := `
		DELETE FROM tag
		WHERE tag.tag_id=$1;
	`
	_, err = r.db.Exec(ctx, query, tagID)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return err
	}

	return nil
}

// AssignTagToCard назначает тег на карточку
func (r *BoardRepository) AssignTagToCard(ctx context.Context, tagID int64, cardID int64) (err error) {
	funcName := "AssignTagToCard"
	query := `
		INSERT INTO tag_to_card (tag_id, card_id)
		VALUES ($1, $2);
	`
	_, err = r.db.Exec(ctx, query, tagID, cardID)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return fmt.Errorf("%s (query): %w", funcName, err)
	}

	return nil
}

// DeassignTagFromCard убирает назначение тега с карточки
func (r *BoardRepository) DeassignTagFromCard(ctx context.Context, tagID int64, cardID int64) (err error) {
	funcName := "DeassignTagFromCard"
	query := `
		DELETE FROM tag_to_card
		WHERE tag_id = $1 AND card_id = $2;
	`
	_, err = r.db.Exec(ctx, query, tagID, cardID)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return fmt.Errorf("%s (query): %w", funcName, err)
	}

	return nil
}
