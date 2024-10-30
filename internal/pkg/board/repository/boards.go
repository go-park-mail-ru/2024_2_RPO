package repository

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const defaultImageUrl = "/static/img/backgroundImage.png"

type BoardRepository struct {
	db *pgxpool.Pool
}

func CreateBoardRepository(db *pgxpool.Pool) *BoardRepository {
	return &BoardRepository{db: db}
}

// CreateBoard creates a new board in the database.
func (r *BoardRepository) CreateBoard(name string, createdBy int64) (*models.Board, error) {
	query := `
		INSERT INTO board (name, description, created_by)
		VALUES ($1, $2, $3)
		RETURNING b_id, name, description, created_at, updated_at
	`
	var board models.Board
	err := r.db.QueryRow(context.Background(), query, name, "", createdBy).Scan(
		&board.Id,
		&board.Name,
		&board.Description,
		&board.CreatedAt,
		&board.UpdatedAt,
	)
	board.BackgroundImageUrl = defaultImageUrl
	if err != nil {
		return nil, fmt.Errorf("CreateBoard: %w", err)
	}
	return &board, nil
}

// GetBoard retrieves a board by its ID.
func (r *BoardRepository) GetBoard(boardID int64) (*models.Board, error) {
	query := `
		SELECT
		b.b_id, b.name,
		b.description, b.created_at, b.updated_at,
		file.file_uuid, file.file_extension
		FROM board AS b
		LEFT JOIN user_uploaded_file AS file ON file.file_uuid=b.avatar_file_uuid
		WHERE b.b_id = $1;
	`
	var board models.Board
	var fileUuid string
	var fileExtension string
	err := r.db.QueryRow(context.Background(), query, boardID).Scan(
		&board.Id,
		&board.Name,
		&board.Description,
		&board.CreatedAt,
		&board.UpdatedAt,
		&fileUuid,
		&fileExtension,
	)
	if fileUuid == "" {
		board.BackgroundImageUrl = defaultImageUrl
	} else {
		board.BackgroundImageUrl = fmt.Sprintf("%s%s.%s", os.Getenv("USER_UPLOADS_URL"), fileUuid, fileExtension)
	}
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetBoard: %w", errs.ErrNotFound)
		}
		return nil, err
	}
	return &board, nil
}

// UpdateBoard updates the specified fields of a board.
func (r *BoardRepository) UpdateBoard(boardID int64, data *models.BoardPutRequest) error {
	query := `
		UPDATE board
		SET name=$1, description=$2, updated_at = CURRENT_TIMESTAMP
		WHERE b_id = $3;
	`

	tag, err := r.db.Exec(context.Background(), query, data.NewName, data.NewDescription, boardID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { // Может ли эта ошибка появиться?
			return fmt.Errorf("UpdateBoard: %w", errs.ErrNotFound)
		}
		return fmt.Errorf("UpdateBoard: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("UpdateBoard: %w", errs.ErrNotFound)
	}
	return nil
}

// DeleteBoard удаляет доску по Id
func (r *BoardRepository) DeleteBoard(boardId int64) error {
	query := `
		DELETE FROM board
		WHERE b_id = $1
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
func (r *BoardRepository) GetBoardsForUser(userID int64) (boardArray []models.Board, err error) {
	query := `
		SELECT b.b_id, b.name, b.description,
		b.created_at, b.updated_at,
		f.file_uuid, f.file_extension
		FROM board AS b
		JOIN user_to_board AS ub ON b.b_id = ub.b_id
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
			&board.Id,
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
		boardArray = append(boardArray, board)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetBoardsForUser: %w", err)
	}

	return boardArray, nil
}
