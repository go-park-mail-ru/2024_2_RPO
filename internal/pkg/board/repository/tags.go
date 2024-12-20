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

// CreateNewTag создает тег на доске
func (r *BoardRepository) CreateNewTag(ctx context.Context, boardID int64, data *models.TagRequest) (newTag *models.Tag, err error) {
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
func (r *BoardRepository) UpdateTag(ctx context.Context, tagID int64, data *models.TagRequest) (updatedTag *models.Tag, err error) {
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
	tag, err := r.db.Exec(ctx, query, tagID)
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", funcName, errs.ErrNotFound)
	}

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
	tag, err := r.db.Exec(ctx, query, tagID, cardID)
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", funcName, errs.ErrNotFound)
	}

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
	tag, err := r.db.Exec(ctx, query, tagID, cardID)
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", funcName, errs.ErrNotFound)
	}

	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return fmt.Errorf("%s (query): %w", funcName, err)
	}

	return nil
}

func (r *BoardRepository) GetTagsForBoard(ctx context.Context, boardID int64) (tags []models.Tag, err error) {
	funcName := "GetTagsForBoard"
	query := `
		SELECT t.tag_id, t.color, t.title
		FROM tag AS t
		WHERE t.board_id = $1;
	`
	rows, err := r.db.Query(ctx, query, boardID)
	logging.Debug(ctx, funcName, " query has err: ", err)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return make([]models.Tag, 0), nil
		}
		return nil, fmt.Errorf("%s (query): %w", funcName, err)
	}

	defer rows.Close()

	tags = make([]models.Tag, 0)
	for rows.Next() {
		var tag models.Tag
		err := rows.Scan(
			&tag.ID,
			&tag.Color,
			&tag.Text,
		)
		if err != nil {
			return nil, fmt.Errorf("%s (scan): %w", funcName, err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
