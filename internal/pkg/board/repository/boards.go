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

// CreateBoard инсертит новую доску (ВНИМАНИЕ! ОН НЕ ДОБАВЛЯЕТ СОЗДАТЕЛЯ НА ДОСКУ!)
func (r *BoardRepository) CreateBoard(ctx context.Context, name string, userID int) (*models.Board, error) {
	query := `
		INSERT INTO board (name, description, created_by)
		VALUES ($1, $2, $3)
		RETURNING board_id, name, description, created_at, updated_at
	`
	var board models.Board
	err := r.db.QueryRow(ctx, query, name, "", userID).Scan(
		&board.ID,
		&board.Name,
		&board.Description,
		&board.CreatedAt,
		&board.UpdatedAt,
	)
	logging.Debug(ctx, "CreateBoard query has err: ", err)
	if err != nil {
		return nil, fmt.Errorf("CreateBoard: %w", err)
	}
	board.BackgroundImageURL = uploads.DefaultBackgroundURL
	return &board, nil
}

// GetBoard получает доску по ID
func (r *BoardRepository) GetBoard(ctx context.Context, boardID int, userID int) (*models.Board, error) {
	query := `
    SELECT
        b.board_id,
        b.name,
        b.description,
        b.created_at,
        b.updated_at,
        ub.last_visit_at,
        COALESCE(file.file_uuid::text,''),
        COALESCE(file.file_extension,'')
    FROM board AS b
    LEFT JOIN user_to_board AS ub ON ub.board_id = b.board_id AND ub.u_id = $1
    LEFT JOIN user_uploaded_file AS file ON file.file_uuid=b.background_image_uuid
    WHERE b.board_id = $2;
    `
	var board models.Board
	var fileUUID string
	var fileExtension string
	err := r.db.QueryRow(ctx, query, userID, boardID).Scan(
		&board.ID,
		&board.Name,
		&board.Description,
		&board.CreatedAt,
		&board.UpdatedAt,
		&board.LastVisitAt,
		&fileUUID,
		&fileExtension,
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
func (r *BoardRepository) UpdateBoard(ctx context.Context, boardID int, userID int, data *models.BoardPutRequest) (updatedBoard *models.Board, err error) {
	query := `
		UPDATE board
		SET name=$1, description=$2, updated_at = CURRENT_TIMESTAMP
		WHERE board_id = $3;
	`

	tag, err := r.db.Exec(ctx, query, data.NewName, data.NewDescription, boardID)
	logging.Debug(ctx, "UpdateBoard query has err: ", err, " tag: ", tag)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("UpdateBoard: %w", errs.ErrNotFound)
		}
		return nil, fmt.Errorf("UpdateBoard: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return nil, fmt.Errorf("UpdateBoard: %w", errs.ErrNotFound)
	}
	return r.GetBoard(ctx, boardID, userID)
}

// DeleteBoard удаляет доску по Id
func (r *BoardRepository) DeleteBoard(ctx context.Context, boardID int) error {
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
func (r *BoardRepository) GetBoardsForUser(ctx context.Context, userID int) (boardArray []models.Board, err error) {
	query := `
		SELECT b.board_id, b.name, b.description, b.created_at, b.updated_at,
		COALESCE(f.file_uuid::text, ''),
		COALESCE(f.file_extension, '')
		FROM user_to_board AS ub
		JOIN board AS b ON b.board_id = ub.board_id
		LEFT JOIN user_uploaded_file AS f ON f.file_uuid=b.background_image_uuid
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
			&board.Description,
			&board.CreatedAt,
			&board.UpdatedAt,
			&fileUUID,
			&fileExtension,
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

func (r *BoardRepository) SetBoardBackground(ctx context.Context, userID int, boardID int, fileExtension string, fileSize int) (fileName string, err error) {
	query1 := `
	INSERT INTO user_uploaded_file
	(file_extension, created_at, created_by, "size")
	VALUES ($1, CURRENT_TIMESTAMP, $2, $3)
	RETURNING file_uuid::text;
	`
	query2 := `
	UPDATE board
	SET background_image_uuid=to_uuid($1)
	WHERE board_id=$2;
	`
	var fileUUID string
	row := r.db.QueryRow(ctx, query1, fileExtension, userID, fileSize)
	err = row.Scan(&fileUUID)
	logging.Debug(ctx, "SetBoardBackground query 1 has err: ", err)
	if err != nil {
		return "", fmt.Errorf("SetBoardBackground (register file): %w", err)
	}
	tag, err := r.db.Exec(ctx, query2, fileUUID, boardID)
	logging.Debug(ctx, "SetBoardBackground query 2 has err: ", err, "tag: ", tag)
	if err != nil {
		return "", fmt.Errorf("SetBoardBackground (update board): %w", err)
	}
	if tag.RowsAffected() == 0 {
		return "", fmt.Errorf("SetBoardBackground (update board): no rows affected")
	}
	return uploads.JoinFilePath(fileUUID, fileExtension), nil
}

func (r *BoardRepository) UpdateLastVisit(ctx context.Context, userID int, boardID int) error {
	query := `
	UPDATE user_to_board
    SET last_visit_at = NOW()
    WHERE u_id = $1 AND board_id = $2;
	`

	_, err := r.db.Exec(ctx, query, userID, boardID)
	return err
}
