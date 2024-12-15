package repository

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/utils/logging"
	"RPO_back/internal/pkg/utils/pgxiface"
	"RPO_back/internal/pkg/utils/uploads"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type BoardRepository struct {
	db pgxiface.PgxIface
}

func CreateBoardRepository(db pgxiface.PgxIface) *BoardRepository {
	return &BoardRepository{db: db}
}

// CreateBoard создаёт новую доску и добавляет создателя на неё
func (r *BoardRepository) CreateBoard(ctx context.Context, name string, userID int64) (*models.Board, error) {
	funcName := "CreateBoard"
	query := `
	WITH inserted_board AS (
		INSERT INTO board (name, created_by)
		VALUES ($1, $2)
		RETURNING board_id
	),
	create_board AS (
		INSERT INTO user_to_board (u_id, board_id, added_by, updated_by, role)
		SELECT $2, board_id, $2, $2, 'admin'
		FROM inserted_board
	)
	SELECT board_id FROM inserted_board;
	`
	var boardID int64
	err := r.db.QueryRow(ctx, query, name, userID).Scan(
		&boardID,
	)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", funcName, err)
	}
	board, err := r.GetBoard(ctx, boardID, userID)
	return board, err
}

// GetBoard получает доску по ID
func (r *BoardRepository) GetBoard(ctx context.Context, boardID int64, userID int64) (*models.Board, error) {
	query := `
    SELECT
        b.board_id,
        b.name,
        b.created_at,
        b.updated_at,
        COALESCE(ub.last_visit_at, CURRENT_TIMESTAMP),
        COALESCE(file.file_uuid::text,''),
        COALESCE(file.file_extension,''),
		ub.invite_link_uuid::text
    FROM board AS b
    LEFT JOIN user_to_board AS ub ON ub.board_id = b.board_id AND ub.u_id = $1
    LEFT JOIN user_uploaded_file AS file ON file.file_id=b.background_image_id
    WHERE b.board_id = $2;
    `
	var board models.Board
	var fileUUID string
	var fileExtension string
	err := r.db.QueryRow(ctx, query, userID, boardID).Scan(
		&board.ID,
		&board.Name,
		&board.CreatedAt,
		&board.UpdatedAt,
		&board.LastVisitAt,
		&fileUUID,
		&fileExtension,
		&board.MyInviteUUID,
	)
	logging.Debug(ctx, "GetBoard query has err: ", err)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetBoard: %w", errs.ErrNotFound)
		}
		return nil, err
	}

	board.BackgroundImageURL = uploads.JoinFileURL(fileUUID, fileExtension, uploads.DefaultBackgroundURL)
	return &board, nil
}

// UpdateBoard обновляет информацию о доске
func (r *BoardRepository) UpdateBoard(ctx context.Context, boardID int64, userID int64, data *models.BoardRequest) (updatedBoard *models.Board, err error) {
	funcName := "UpdateBoard"
	query := `
		UPDATE board
		SET name=$1, updated_at=CURRENT_TIMESTAMP
		WHERE board_id=$2;
	`

	tag, err := r.db.Exec(ctx, query, data.NewName, boardID)
	logging.Debug(ctx, funcName, " query has err: ", err, " tag: ", tag)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", funcName, errs.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", funcName, err)
	}
	if tag.RowsAffected() == 0 {
		return nil, fmt.Errorf("%s: %w", funcName, errs.ErrNotFound)
	}
	return r.GetBoard(ctx, boardID, userID)
}

// DeleteBoard удаляет доску по ID
func (r *BoardRepository) DeleteBoard(ctx context.Context, boardID int64) error {
	query := `
		DELETE FROM board
		WHERE board_id = $1;
	`
	tag, err := r.db.Exec(ctx, query, boardID)
	logging.Debug(ctx, "DeleteBoard query has err: ", err, " tag: ", tag)
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
func (r *BoardRepository) GetBoardsForUser(ctx context.Context, userID int64) (boardArray []models.Board, err error) {
	query := `
		SELECT b.board_id, b.name, b.created_at, b.updated_at,
		COALESCE(f.file_uuid::text, ''),
		COALESCE(f.file_extension, ''), ub.invite_link_uuid
		FROM user_to_board AS ub
		JOIN board AS b ON b.board_id = ub.board_id
		LEFT JOIN user_uploaded_file AS f ON f.file_id=b.background_image_id
		WHERE ub.u_id = $1
	`
	rows, err := r.db.Query(ctx, query, userID)
	logging.Debug(ctx, "GetBoardsForUser query has err: ", err)
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
		var fileUUID, fileExtension string
		err := rows.Scan(
			&board.ID,
			&board.Name,
			&board.CreatedAt,
			&board.UpdatedAt,
			&fileUUID,
			&fileExtension,
			&board.MyInviteUUID,
		)
		if err != nil {
			return nil, fmt.Errorf("GetBoardsForUser (for rows): %w", err)
		}
		board.BackgroundImageURL = uploads.JoinFileURL(fileUUID, fileExtension, uploads.DefaultBackgroundURL)
		boardArray = append(boardArray, board)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetBoardsForUser: %w", err)
	}

	return boardArray, nil
}

// SetBoardBackground задаёт файл заднего фона доски
func (r *BoardRepository) SetBoardBackground(ctx context.Context, userID int64, boardID int64, fileID int64) (newBoard *models.Board, err error) {
	funcName := "SetBoardBackground"
	query := `
	WITH update_board AS (
		UPDATE board
		SET background_image_id=$1,
			updated_at=CURRENT_TIMESTAMP
		WHERE board_id=$2
		RETURNING board_id
	)
	SELECT b.board_id, b.name, b.created_at, b.updated_at,
	COALESCE(f.file_uuid::text, ''), COALESCE(f.file_extension, '')
	FROM board AS b
	LEFT JOIN user_uploaded_file AS f ON f.file_id=b.background_image_id
	WHERE b.board_id=$2;
	`

	newBoard = &models.Board{}
	var fileUUID, fileExt string

	row := r.db.QueryRow(ctx, query, fileID, boardID)
	err = row.Scan(
		&newBoard.ID,
		&newBoard.Name,
		&newBoard.CreatedAt,
		&newBoard.UpdatedAt,
		&fileUUID,
		&fileExt,
	)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return nil, fmt.Errorf("%s (query): %w", funcName, err)
	}
	newBoard.BackgroundImageURL = uploads.JoinFileURL(fileUUID, fileExt, uploads.DefaultBackgroundURL)
	return newBoard, nil
}
