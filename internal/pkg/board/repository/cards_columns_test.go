package repository

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/utils/misc"
	"context"
	"errors"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetCardsForBoard(t *testing.T) {
	mock, err := pgxmock.NewConn()
	assert.NoError(t, err)
	defer mock.Close(context.Background())

	repo := CreateBoardRepository(mock)

	mock.ExpectQuery("SELECT c.card_id,").
		WithArgs(1).
		WillReturnError(errors.New("some error"))

	_, err = repo.GetCardsForBoard(context.Background(), 1)
	assert.Error(t, err)
}

func TestCreateNewCard(t *testing.T) {
	mock, err := pgxmock.NewConn()
	assert.NoError(t, err)
	defer mock.Close(context.Background())

	repo := CreateBoardRepository(mock)

	mock.ExpectQuery("WITH col_check AS").
		WithArgs(1, 1, "New Card").
		WillReturnError(errors.New("some error"))

	_, err = repo.CreateNewCard(context.Background(), 1, "New Card")
	assert.Error(t, err)
}

func TestUpdateCard(t *testing.T) {
	mock, err := pgxmock.NewConn()
	assert.NoError(t, err)
	defer mock.Close(context.Background())

	repo := CreateBoardRepository(mock)

	mock.ExpectQuery("UPDATE card").
		WithArgs("Updated Title", 2, 1, 1).
		WillReturnError(errors.New("some error"))

	data := models.CardPatchRequest{
		NewTitle: misc.StringPtr("Updated Title"),
	}

	_, err = repo.UpdateCard(context.Background(), 1, data)
	assert.Error(t, err)
}

func TestDeleteCard(t *testing.T) {
	mock, err := pgxmock.NewConn()
	assert.NoError(t, err)
	defer mock.Close(context.Background())

	repo := CreateBoardRepository(mock)

	mock.ExpectExec("DELETE FROM card").
		WithArgs(1, 1).
		WillReturnError(errors.New("some error"))

	err = repo.DeleteCard(context.Background(), 1)
	assert.Error(t, err)
}

func TestGetColumnsForBoard_Error(t *testing.T) {
	dbMock, err := pgxmock.NewPool()
	assert.NoError(t, err)

	boardRepo := CreateBoardRepository(dbMock)

	dbMock.ExpectQuery("SELECT col_id, title FROM kanban_column WHERE board_id = \\$1;").
		WithArgs(1).
		WillReturnError(errors.New("some error"))

	_, err = boardRepo.GetColumnsForBoard(context.Background(), 1)

	assert.Error(t, err)
	assert.NoError(t, dbMock.ExpectationsWereMet())
}

func TestCreateColumn_Error(t *testing.T) {
	dbMock, err := pgxmock.NewPool()
	assert.NoError(t, err)

	boardRepo := CreateBoardRepository(dbMock)

	dbMock.ExpectQuery("INSERT INTO kanban_column \\(board_id, title, created_at, updated_at\\) VALUES \\(\\$1, \\$2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP\\) RETURNING col_id, title;").
		WithArgs(1, "Test Column").
		WillReturnError(errors.New("some error"))

	_, err = boardRepo.CreateColumn(context.Background(), 1, "Test Column")

	assert.Error(t, err)
	assert.NoError(t, dbMock.ExpectationsWereMet())
}

func TestUpdateColumn_Error(t *testing.T) {
	dbMock, err := pgxmock.NewPool()
	assert.NoError(t, err)

	boardRepo := CreateBoardRepository(dbMock)

	dbMock.ExpectQuery("UPDATE kanban_column SET title = \\$1, updated_at = CURRENT_TIMESTAMP WHERE col_id = \\$2 AND board_id = \\$3 RETURNING col_id, title;").
		WithArgs("Updated Title", 1, 1).
		WillReturnError(errors.New("some error"))

	_, err = boardRepo.UpdateColumn(context.Background(), 1, models.ColumnRequest{NewTitle: "Updated Title"})

	assert.Error(t, err)
	assert.NoError(t, dbMock.ExpectationsWereMet())
}

func TestDeleteColumn_Error(t *testing.T) {
	dbMock, err := pgxmock.NewPool()
	assert.NoError(t, err)

	boardRepo := CreateBoardRepository(dbMock)

	dbMock.ExpectExec("DELETE FROM kanban_column WHERE col_id = \\$1 AND board_id = \\$2;").
		WithArgs(1, 1).
		WillReturnError(errors.New("some error"))

	err = boardRepo.DeleteColumn(context.Background(), 1)

	assert.Error(t, err)
	assert.NoError(t, dbMock.ExpectationsWereMet())
}
