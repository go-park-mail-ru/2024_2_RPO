package repository

import (
	"RPO_back/internal/models"
	"context"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
)

func TestGetUserByEmail_Success(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("unable to create mock connection: %v", err)
	}
	defer mock.Close(ctx)
	email := "kaymekaydex@mail.ru"

	repo := &AuthRepository{db: mock}

	query := `SELECT u_id, nickname, email, description, joined_at, updated_at, password_hash FROM "user" WHERE email=\$1;`
	rows := pgxmock.NewRows([]string{"u_id", "nickname", "email", "description", "joined_at", "updated_at", "password_hash"}).
		AddRow(1, "testnickname", email, "test description", time.Now(), time.Now(), "hashedpassword")

	mock.ExpectQuery(query).WithArgs(email).WillReturnRows(rows)

	_, err = repo.GetUserByEmail(ctx, email)
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
}

func TestGetUserByID_Success(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("unable to create mock connection: %v", err)
	}
	defer mock.Close(ctx)
	email := "kaymekaydex@mail.ru"

	repo := &AuthRepository{db: mock}

	query := `SELECT u_id, nickname, email, description, joined_at, updated_at, password_hash FROM "user" WHERE u_id=\$1;`
	rows := pgxmock.NewRows([]string{"u_id", "nickname", "email", "description", "joined_at", "updated_at", "password_hash"}).
		AddRow(1337, "testnickname", email, "test description", time.Now(), time.Now(), "hashedpassword")

	mock.ExpectQuery(query).WithArgs(1337).WillReturnRows(rows)

	_, err = repo.GetUserByID(ctx, 1337)
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
}

func TestCreateUser_Success(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("unable to create mock connection: %v", err)
	}
	defer mock.Close(ctx)

	repo := &AuthRepository{db: mock}

	query := `INSERT INTO "user" \(nickname, email, password_hash, description, joined_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\) RETURNING u_id, nickname, email, password_hash, description, joined_at, updated_at`
	rows := pgxmock.NewRows([]string{"u_id", "nickname", "email", "password_hash", "description", "joined_at", "updated_at"}).
		AddRow(1, "testnickname", "testemail", "hashedpassword", "", time.Now(), time.Now())

	mock.ExpectQuery(query).WithArgs("testnickname", "testemail", "hashedpassword", "", pgxmock.AnyArg(), pgxmock.AnyArg()).WillReturnRows(rows)

	_, err = repo.CreateUser(ctx, &models.UserRegisterRequest{Name: "testnickname", Email: "testemail"}, "hashedpassword")
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
}

func TestCheckUniqueCredentials_Success(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("unable to create mock connection: %v", err)
	}
	defer mock.Close(ctx)

	repo := &AuthRepository{db: mock}

	query1 := `SELECT COUNT\(\*\) FROM "user" WHERE nickname = \$1;`
	query2 := `SELECT COUNT\(\*\) FROM "user" WHERE email = \$1;`

	mock.ExpectQuery(query1).WithArgs("testnickname").WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(query2).WithArgs("testemail").WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(0))

	err = repo.CheckUniqueCredentials(ctx, "testnickname", "testemail")
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
}

func TestSetNewPasswordHash_Success(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("unable to create mock connection: %v", err)
	}
	defer mock.Close(ctx)

	repo := &AuthRepository{db: mock}

	query := `UPDATE "user" SET password_hash=\$1 WHERE u_id=\$2;`

	mock.ExpectExec(query).WithArgs("newhashedpassword", 228).WillReturnResult(pgxmock.NewResult("1", 1))

	err = repo.SetNewPasswordHash(ctx, 228, "newhashedpassword")
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
}
