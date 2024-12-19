package main

// Программа для обслуживания сервиса.
// Возможности:
// * Заполнить базу данных для нагрузочного тестирования;
// *

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	totalBoards      = 1000
	columnsPerBoard  = 10
	cardsPerColumn   = 10
	simpleCards      = 5
	interestingCards = 5

	attachmentsPerCard    = 10
	commentsPerCard       = 10
	checklistItemsPerCard = 10
)

func main() {
	switch os.Args[1] {
	case "elastic-migrator":
		break
	case "fill-db":
		fill_db()
	}
}

func fill_db() {
	ctx := context.Background()

	// Initialize database connection pool
	cfg, err := pgxpool.ParseConfig(dbConnectionString)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	defer pool.Close()

	// Begin the entire operation
	startTime := time.Now()
	log.Println("Starting user and data creation...")

	// Create user
	userID, err := createUser(ctx, pool)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	log.Printf("Created user with u_id: %d", userID)

	// Iterate through boards
	for b := 1; b <= totalBoards; b++ {
		boardName := fmt.Sprintf("board_%d", b)
		boardID, err := createBoard(ctx, pool, userID, boardName)
		if err != nil {
			log.Fatalf("Failed to create board %s: %v", boardName, err)
		}

		// Iterate through columns
		for c := 1; c <= columnsPerBoard; c++ {
			columnTitle := fmt.Sprintf("column_%d", c)
			orderIndex := c
			colID, err := createKanbanColumn(ctx, pool, boardID, columnTitle, orderIndex)
			if err != nil {
				log.Fatalf("Failed to create column %s for board %d: %v", columnTitle, boardID, err)
			}

			// Iterate through cards
			for card := 1; card <= cardsPerColumn; card++ {
				if card <= simpleCards {
					cardTitle := fmt.Sprintf("card_simple_%d", card)
					_, err := createCard(ctx, pool, colID, cardTitle, card, false, userID)
					if err != nil {
						log.Fatalf("Failed to create simple card %s: %v", cardTitle, err)
					}
				} else {
					cardTitle := fmt.Sprintf("card_interesting_%d", card-simpleCards)
					cardID, err := createCard(ctx, pool, colID, cardTitle, card, true, userID)
					if err != nil {
						log.Fatalf("Failed to create interesting card %s: %v", cardTitle, err)
					}

					for a := 1; a <= attachmentsPerCard; a++ {
						originalName := fmt.Sprintf("attachment_%d", a)
						fileID, err := createUploadedFile(ctx, pool, a)
						if err != nil {
							log.Fatalf("Failed to create uploaded file for attachment %s: %v", originalName, err)
						}
						attachmentTitle := fmt.Sprintf("attachment_%d", a)
						err = createAttachment(ctx, pool, cardID, fileID, originalName, userID)
						if err != nil {
							log.Fatalf("Failed to create attachment %s for card %d: %v", originalName, cardID, err)
						}
					}

					for m := 1; m <= commentsPerCard; m++ {
						commentTitle := fmt.Sprintf("comment_%d", m)
						err = createComment(ctx, pool, cardID, commentTitle, userID)
						if err != nil {
							log.Fatalf("Failed to create comment %s for card %d: %v", commentTitle, cardID, err)
						}
					}

					for chk := 1; chk <= checklistItemsPerCard; chk++ {
						checklistTitle := fmt.Sprintf("checklist_%d", chk)
						orderIdx := chk
						err = createChecklistField(ctx, pool, cardID, checklistTitle, orderIdx)
						if err != nil {
							log.Fatalf("Failed to create checklist item %s for card %d: %v", checklistTitle, cardID, err)
						}
					}
				}
			}

			if b%100 == 0 && c == columnsPerBoard {
				log.Printf("Created %d/%d boards", b, totalBoards)
			}
		}
	}

	elapsed := time.Since(startTime)
	log.Printf("Data creation completed in %s", elapsed)
}

func createUser(ctx context.Context, pool *pgxpool.Pool) (int64, error) {
	csatPollDT := time.Now().AddDate(0, 0, 30)

	query := `
		INSERT INTO "user" (nickname, email, password_hash, csat_poll_dt)
		VALUES ($1, $2, $3, $4)
		RETURNING u_id;
	`
	var userID int64
	err := pool.QueryRow(ctx, query, userNickname, userEmail, passwordHash, csatPollDT).Scan(&userID)
	return userID, err
}

func createBoard(ctx context.Context, pool *pgxpool.Pool, userID int64, boardName string) (int64, error) {
	query := `
		INSERT INTO board ("name", created_by)
		VALUES ($1, $2)
		RETURNING board_id;
	`
	var boardID int64
	err := pool.QueryRow(ctx, query, boardName, userID).Scan(&boardID)
	return boardID, err
}

func createKanbanColumn(ctx context.Context, pool *pgxpool.Pool, boardID int64, title string, orderIndex int) (int64, error) {
	query := `
		INSERT INTO kanban_column (board_id, title, order_index)
		VALUES ($1, $2, $3)
		RETURNING col_id;
	`
	var colID int64
	err := pool.QueryRow(ctx, query, boardID, title, orderIndex).Scan(&colID)
	return colID, err
}

func createCard(ctx context.Context, pool *pgxpool.Pool, colID int64, title string, orderIndex int, interesting bool, userID int64) (int64, error) {
	query := `
		INSERT INTO "card" (title, col_id, order_index, is_done)
		VALUES ($1, $2, $3, FALSE)
		RETURNING card_id;
	`
	var cardID int64
	err := pool.QueryRow(ctx, query, title, colID, orderIndex).Scan(&cardID)
	return cardID, err
}

func createUploadedFile(ctx context.Context, pool *pgxpool.Pool, fileNumber int) (int64, error) {
	query := `
		INSERT INTO user_uploaded_file (file_hash, file_extension, "size")
		VALUES ($1, $2, $3)
		RETURNING file_id;
	`
	fileHash := fmt.Sprintf("hash_%d", fileNumber)
	fileExtension := ".txt"
	size := 1024 // example size in bytes

	var fileID int64
	err := pool.QueryRow(ctx, query, fileHash, fileExtension, size).Scan(&fileID)
	return fileID, err
}

func createAttachment(ctx context.Context, pool *pgxpool.Pool, cardID, fileID int64, originalName string, attachedBy int64) error {
	query := `
		INSERT INTO card_attachment (card_id, file_id, original_name, attached_by)
		VALUES ($1, $2, $3, $4);
	`
	_, err := pool.Exec(ctx, query, cardID, fileID, originalName, attachedBy)
	return err
}

func createComment(ctx context.Context, pool *pgxpool.Pool, cardID int64, title string, createdBy int64) error {
	query := `
		INSERT INTO card_comment (card_id, title, created_by)
		VALUES ($1, $2, $3);
	`
	_, err := pool.Exec(ctx, query, cardID, title, createdBy)
	return err
}

func createChecklistField(ctx context.Context, pool *pgxpool.Pool, cardID int64, title string, orderIndex int) error {
	query := `
		INSERT INTO checklist_field (card_id, title, order_index)
		VALUES ($1, $2, $3);
	`
	_, err := pool.Exec(ctx, query, cardID, title, orderIndex)
	return err
}
