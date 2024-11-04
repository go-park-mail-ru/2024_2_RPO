package repository

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/utils/uploads"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BoardRepository struct {
	db *pgxpool.Pool
}

func CreateBoardRepository(db *pgxpool.Pool) *BoardRepository {
	return &BoardRepository{db: db}
}

// CreateBoard creates a new board in the database.
func (r *BoardRepository) CreateBoard(name string, userID int) (*models.Board, error) {
	query := `
		INSERT INTO board (name, description, created_by)
		VALUES ($1, $2, $3)
		RETURNING board_id, name, description, created_at, updated_at
	`
	var board models.Board
	err := r.db.QueryRow(context.Background(), query, name, "", userID).Scan(
		&board.ID,
		&board.Name,
		&board.Description,
		&board.CreatedAt,
		&board.UpdatedAt,
	)
	board.BackgroundImageURL = uploads.DefaultBackgroundURL
	if err != nil {
		return nil, fmt.Errorf("CreateBoard: %w", err)
	}
	return &board, nil
}

// GetBoard retrieves a board by its ID.
func (r *BoardRepository) GetBoard(boardID int) (*models.Board, error) {
	query := `
		SELECT
			b.board_id,
			b.name,
			b.description,
			b.created_at,
			b.updated_at,
			COALESCE(file.file_uuid::text,''),
			COALESCE(file.file_extension,'')
		FROM board AS b
		LEFT JOIN user_uploaded_file AS file ON file.file_uuid=b.background_image_uuid
		WHERE b.board_id = $1;
	`
	var board models.Board
	var fileUUID string
	var fileExtension string
	err := r.db.QueryRow(context.Background(), query, boardID).Scan(
		&board.ID,
		&board.Name,
		&board.Description,
		&board.CreatedAt,
		&board.UpdatedAt,
		&fileUUID,
		&fileExtension,
	)
	board.BackgroundImageURL = uploads.JoinFileName(fileUUID, fileExtension, uploads.DefaultBackgroundURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetBoard: %w", errs.ErrNotFound)
		}
		return nil, err
	}
	return &board, nil
}

// UpdateBoard updates the specified fields of a board.
func (r *BoardRepository) UpdateBoard(boardID int, data *models.BoardPutRequest) (updatedBoard *models.Board, err error) {
	query := `
		UPDATE board
		SET name=$1, description=$2, updated_at = CURRENT_TIMESTAMP
		WHERE board_id = $3;
	`

	tag, err := r.db.Exec(context.Background(), query, data.NewName, data.NewDescription, boardID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("UpdateBoard: %w", errs.ErrNotFound)
		}
		return nil, fmt.Errorf("UpdateBoard: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return nil, fmt.Errorf("UpdateBoard: %w", errs.ErrNotFound)
	}
	return r.GetBoard(boardID)
}

// DeleteBoard удаляет доску по Id
func (r *BoardRepository) DeleteBoard(boardId int) error {
	query := `
		DELETE FROM board
		WHERE board_id = $1;
	`
	tag, err := r.db.Exec(context.Background(), query, boardId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("DeleteBoard: %w", errs.ErrNotFound)
		}
		return fmt.Errorf("DeleteBoard: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("DeleteBoard: %w", errs.ErrNotFound)
	}
	return nil
}

// GetBoardsForUser возвращает все доски, к которым пользователь имеет доступ
func (r *BoardRepository) GetBoardsForUser(userID int) (boardArray []models.Board, err error) {
	query := `
		SELECT b.board_id, b.name, b.description, b.created_at, b.updated_at,
		COALESCE(f.file_uuid::text, ''),
		COALESCE(f.file_extension, '')
		FROM user_to_board AS ub
		JOIN board AS b ON b.board_id = ub.board_id
		LEFT JOIN user_uploaded_file AS f ON f.file_uuid=b.background_image_uuid
		WHERE ub.u_id = $1
	`
	rows, err := r.db.Query(context.Background(), query, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("GetBoardsForUser: %w", err)
	}
	defer rows.Close()

	boardArray = []models.Board{}
	for rows.Next() {
		var board models.Board
		var fileUuid, fileExtension string
		err := rows.Scan(
			&board.ID,
			&board.Name,
			&board.Description,
			&board.CreatedAt,
			&board.UpdatedAt,
			&fileUuid,
			&fileExtension,
		)
		if err != nil {
			return nil, fmt.Errorf("GetBoardsForUser (for rows): %w", err)
		}
		if fileUuid != "" {
			board.BackgroundImageURL = fmt.Sprintf("%s.%s", fileUuid, fileExtension)
		} else {
			board.BackgroundImageURL = uploads.DefaultBackgroundURL
		}
		boardArray = append(boardArray, board)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetBoardsForUser: %w", err)
	}

	return boardArray, nil
}
