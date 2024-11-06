package repository

import (
	"RPO_back/internal/pkg/utils/uploads"
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateBoard(t *testing.T) {
	// Создаем mock SQL базу данных
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer mock.Close()

	repo := CreateBoardRepository(mock)

	ctx := context.Background()

	name := "Sample Board"
	userID := 1

	expectedBoardID := 1
	expectedCreatedAt := time.Now()
	expectedUpdatedAt := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`
	INSERT INTO board (name, description, created_by)
	VALUES ($1, $2, $3)
	RETURNING board_id, name, description, created_at, updated_at
`)).
		WithArgs(name, "", userID).
		WillReturnRows(pgxmock.NewRows([]string{"board_id", "name", "description", "created_at", "updated_at"}).
			AddRow(expectedBoardID, name, "", expectedCreatedAt, expectedUpdatedAt))

	board, err := repo.CreateBoard(ctx, name, userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedBoardID, board.ID)
	assert.Equal(t, name, board.Name)
	assert.Equal(t, "", board.Description)
	assert.Equal(t, expectedCreatedAt, board.CreatedAt)
	assert.Equal(t, expectedUpdatedAt, board.UpdatedAt)
	assert.Equal(t, uploads.DefaultBackgroundURL, board.BackgroundImageURL)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetBoard(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)

	ctx := context.Background()
	boardID := 1

	rows := pgxmock.NewRows([]string{
		"board_id", "name", "description", "created_at", "updated_at", "file_uuid", "file_extension",
	}).AddRow(
		boardID, "Test Board", "Description", time.Now(), time.Now(), "uuid-1234", ".jpg",
	)

	mock.ExpectQuery("^SELECT (.+) FROM board AS b LEFT JOIN user_uploaded_file AS file ON file.file_uuid=b.background_image_uuid WHERE b.board_id = \\$1;").
		WithArgs(boardID).
		WillReturnRows(rows)

	boardRepo := CreateBoardRepository(mock)

	board, err := boardRepo.GetBoard(ctx, boardID)
	assert.NoError(t, err)
	assert.NotNil(t, board)
}

func TestGetBoard_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)

	ctx := context.Background()
	boardID := 1

	mock.ExpectQuery("^SELECT (.+) FROM board AS b LEFT JOIN user_uploaded_file AS file ON file.file_uuid=b.background_image_uuid WHERE b.board_id = \\$1;").
		WithArgs(boardID).
		WillReturnError(pgx.ErrNoRows)

	boardRepo := CreateBoardRepository(mock)

	board, err := boardRepo.GetBoard(ctx, boardID)
	assert.Error(t, err)
	assert.Nil(t, board)
}

func TestGetBoard_QueryError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)

	ctx := context.Background()
	boardID := 1

	mock.ExpectQuery("^SELECT (.+) FROM board AS b LEFT JOIN user_uploaded_file AS file ON file.file_uuid=b.background_image_uuid WHERE b.board_id = \\$1;").
		WithArgs(boardID).
		WillReturnError(fmt.Errorf("some query error"))

	boardRepo := CreateBoardRepository(mock)

	board, err := boardRepo.GetBoard(ctx, boardID)
	assert.Error(t, err)
	assert.Nil(t, board)
}

func TestGetMembersWithPermissions(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)

	boardID := 1

	// Устанавливаем ожидания для запроса GetBoard.
	mock.ExpectQuery("SELECT (.+) FROM board").
		WithArgs(boardID).
		WillReturnRows(pgxmock.NewRows([]string{
			"board_id", "name", "description", "created_at", "updated_at", "file_uuid", "file_extension",
		}).AddRow(
			boardID, "Test Board", "Description", time.Now(), time.Now(), "uuid-1234", ".jpg",
		))

	// Настраиваем mock для основного SQL-запроса функции GetMembersWithPermissions.
	query := `
	SELECT (.+) FROM "user"
	`

	// Создаем строки, которые будут возвращены в ответ на вызов метода Query.
	rows := pgxmock.NewRows([]string{
		"u_id", "nickname", "email", "description", "joined_at", "updated_at",
		"role", "added_at", "updated_at",
		"adder_id", "adder_nickname", "adder_email", "adder_description", "adder_joined_at", "adder_updated_at",
		"updater_id", "updater_nickname", "updater_email", "updater_description", "updater_joined_at", "updater_updated_at",
		"member_avatar_uuid", "member_avatar_ext",
		"adder_avatar_uuid", "adder_avatar_ext",
		"updater_avatar_uuid", "updater_avatar_ext",
	}).AddRow(
		1, "User1", "user1@example.com", "Description1", time.Now(), time.Now(),
		"admin", time.Now(), time.Now(),
		2, "Adder1", "adder1@example.com", "Adder Description", time.Now(), time.Now(),
		3, "Updater1", "updater1@example.com", "Updater Description", time.Now(), time.Now(),
		"", "", "", "", "", "", // UUIDs and file extensions
	)

	mock.ExpectQuery(query).
		WithArgs(boardID).
		WillReturnRows(rows)

	repo := CreateBoardRepository(mock)
	_, err = repo.GetMembersWithPermissions(context.Background(), boardID)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSetMemberRole_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)

	ctx := context.Background()
	repo := CreateBoardRepository(mock)

	uid := 123
	bid := 456

	mock.ExpectExec(`
	UPDATE user_to_board`).
		WithArgs(uid, bid).
		WillReturnResult(pgxmock.NewResult("1", 1))

	mock.ExpectQuery(`SELECT (.*) FROM "user"`).
		WithArgs(pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"u_id", "nickname", "email",
			"description",
			"joined_at",
			"updated_at",
			"file_uuid",
			"file_extension"}).AddRow(2, "RVasily", "rvasily@yandex.ru", "Hello! I am user of Pumpkin", time.Now(), time.Now(), "", ""))

	mock.ExpectQuery("WITH board_check").
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"role", "added_at", "updated_at", "added_by", "updated_by"}).
			AddRow("admin", time.Now(), time.Now(), 2, 3))

	mock.ExpectQuery(`SELECT (.*) FROM "user"`).
		WithArgs(pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"u_id", "nickname", "email",
			"description",
			"joined_at",
			"updated_at",
			"file_uuid",
			"file_extension"}).AddRow(2, "RVasily", "rvasily@yandex.ru", "Hello! I am user of Pumpkin", time.Now(), time.Now(), "", ""))
	mock.ExpectQuery(`SELECT (.*) FROM "user"`).
		WithArgs(pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"u_id", "nickname", "email",
			"description",
			"joined_at",
			"updated_at",
			"file_uuid",
			"file_extension"}).AddRow(2, "RVasily", "rvasily@yandex.ru", "Hello! I am user of Pumpkin", time.Now(), time.Now(), "", ""))

	member, err := repo.SetMemberRole(ctx, bid, uid, "viewer")
	assert.NoError(t, err)
	assert.NotNil(t, member)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRemoveMember_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)

	ctx := context.Background()
	repo := CreateBoardRepository(mock)

	mock.ExpectExec(`DELETE FROM`).
		WithArgs(123, 456).
		WillReturnResult(pgxmock.NewResult("1", 1))

	err = repo.RemoveMember(ctx, 123, 456)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddMember_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	repo := CreateBoardRepository(mock)

	mock.ExpectQuery(`SELECT (.*) FROM "user"`).
		WithArgs(pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"u_id", "nickname", "email",
			"description",
			"joined_at",
			"updated_at",
			"file_uuid",
			"file_extension"}).AddRow(2, "RVasily", "rvasily@yandex.ru", "Hello! I am user of Pumpkin", time.Now(), time.Now(), "", ""))
	mock.ExpectQuery("WITH board_check").
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnError(pgx.ErrNoRows)

	mock.ExpectExec(`INSERT INTO user_to_board`).
		WithArgs(123, 456, 789).WillReturnResult(pgxmock.NewResult("1", 1))

	mock.ExpectQuery(`SELECT (.*) FROM "user"`).
		WithArgs(pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"u_id", "nickname", "email",
			"description",
			"joined_at",
			"updated_at",
			"file_uuid",
			"file_extension"}).AddRow(2, "RVasily", "rvasily@yandex.ru", "Hello! I am user of Pumpkin", time.Now(), time.Now(), "", ""))

	mock.ExpectQuery("WITH board_check").
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"role", "added_at", "updated_at", "added_by", "updated_by"}).
			AddRow("admin", time.Now(), time.Now(), 2, 3))

	mock.ExpectQuery(`SELECT (.*) FROM "user"`).
		WithArgs(pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"u_id", "nickname", "email",
			"description",
			"joined_at",
			"updated_at",
			"file_uuid",
			"file_extension"}).AddRow(2, "RVasily", "rvasily@yandex.ru", "Hello! I am user of Pumpkin", time.Now(), time.Now(), "", ""))

	mock.ExpectQuery(`SELECT (.*) FROM "user"`).
		WithArgs(pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"u_id", "nickname", "email",
			"description",
			"joined_at",
			"updated_at",
			"file_uuid",
			"file_extension"}).AddRow(2, "RVasily", "rvasily@yandex.ru", "Hello! I am user of Pumpkin", time.Now(), time.Now(), "", ""))

	member, err := repo.AddMember(context.Background(), 456, 789, 123)
	assert.NoError(t, err)
	assert.NotNil(t, member)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByNickname_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)

	ctx := context.Background()
	repo := CreateBoardRepository(mock)

	rows := pgxmock.NewRows([]string{"u_id", "nickname", "email", "description", "joined_at", "updated_at"}).
		AddRow(1, "testuser", "test@example.com", "description", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT u_id, nickname, email, description, joined_at, updated_at FROM "user" WHERE nickname=\$1`).
		WithArgs("testuser").
		WillReturnRows(rows)

	user, err := repo.GetUserByNickname(ctx, "testuser")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}
