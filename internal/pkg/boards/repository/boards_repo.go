package repository

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/boards"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const defaultImageUrl = "/static/img/backgroundImage.png"

type BoardRepository struct {
	db *pgxpool.Pool
}

// NewBoardRepository creates a new instance of BoardRepository.
func NewBoardRepository(db *pgxpool.Pool) *BoardRepository {
	return &BoardRepository{db: db}
}

// CreateBoard creates a new board in the database.
func (r *BoardRepository) CreateBoard(name string, createdBy int64) (*models.Board, error) {
	query := `
		INSERT INTO board (name, description, created_by)
		VALUES ($1, $2, $3)
		RETURNING b_id, name, description, created_at, created_by, updated_at
	`
	var board models.Board
	err := r.db.QueryRow(context.Background(), query, name, "", createdBy).Scan(
		&board.Id,
		&board.Name,
		&board.Description,
		&board.BackgroundImageURL,
		&board.CreatedAt,
		&board.CreatedBy,
		&board.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("CreateBoard: %w", err)
	}
	return &board, nil
}

// GetBoard retrieves a board by its ID.
func (r *BoardRepository) GetBoard(boardID int64) (*Board, error) {
	query := `
		SELECT b_id, name, description, created_at, created_by, updated_at
		FROM board
		WHERE b_id = $1
	`
	var board models.Board
	err := r.db.QueryRow(context.Background(), query, boardID).Scan(
		&board.Id,
		&board.Name,
		&board.Description,
		&board.BackgroundImageURL,
		&board.CreatedAt,
		&board.CreatedBy,
		&board.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetBoard: %w", boards.ErrNotFound)
		}
		return nil, err
	}
	return &board, nil
}

// UpdateBoard updates the specified fields of a board.
func (r *BoardRepository) UpdateBoard(boardID int64, data *boards.BoardPutRequest) error {
	query := `
		UPDATE board
		SET name=$1, description=$2, updated_at = CURRENT_TIMESTAMP
		WHERE b_id = $3;
	`

	tag, err := r.db.Exec(context.Background(), query, data.NewName, data.NewDescription, boardID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { // Может ли эта ошибка появиться?
			return fmt.Errorf("UpdateBoard: %w", boards.ErrNotFound)
		}
		return fmt.Errorf("UpdateBoard: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("UpdateBoard: %w", boards.ErrNotFound)
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
			return fmt.Errorf("DeleteBoard: %w", boards.ErrNotFound)
		}
		return fmt.Errorf("DeleteBoard: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("DeleteBoard: %w", boards.ErrNotFound)
	}
	return nil
}

// GetBoardsForUser возвращает все доски, к которым пользователь имеет доступ
func (r *BoardRepository) GetBoardsForUser(userID int64) (boardArray []models.Board, err error) {
	query := `
		SELECT b.b_id, b.name, b.description, b.backgroundImageUrl, b.created_at, b.created_by, b.updated_at
		FROM board b
		JOIN user_to_board ub ON b.b_id = ub.b_id
		WHERE ub.u_id = $1
	`
	rows, err := r.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, fmt.Errorf("GetBoardsForUser: %w", err)
	}
	defer rows.Close()

	boardArray = []models.Board{}
	for rows.Next() {
		var board models.Board
		err := rows.Scan(
			&board.Id,
			&board.Name,
			&board.Description,
			&board.BackgroundImageURL,
			&board.CreatedAt,
			&board.CreatedBy,
			&board.UpdatedAt,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("GetBoardsForUser: %w", boards.ErrNotFound)
			}
			return nil, fmt.Errorf("GetBoardsForUser: %w", err)
		}
		boardArray = append(boardArray, board)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetBoardsForUser: %w", err)
	}

	return boardArray, nil
}
