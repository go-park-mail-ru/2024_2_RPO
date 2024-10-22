package boards

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/auth/repository"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetMyBoardsHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := repository.GetUserId(w, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	db, err := repository.GetDbConnection()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query(context.Background(), `
        SELECT b.b_id, b.description, b.created_at, b.updated_at, b."name"
        FROM Board b
        INNER JOIN User_to_Board ub ON b.b_id = ub.b_id
        WHERE ub.u_id = $1
    `, userId)
	if err != nil {
		http.Error(w, "Failed to fetch boards", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var boards []models.Board = make([]models.Board, 0)
	for rows.Next() {
		var board models.Board
		if err := rows.Scan(&board.ID, &board.Description, &board.CreatedAt, &board.UpdatedAt, &board.Name); err != nil {
			http.Error(w, "Failed to parse board data", http.StatusInternalServerError)
			fmt.Printf("Failed to parse board data: %s", err.Error())
			return
		}
		boards = append(boards, board)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Failed to fetch boards", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(boards)
}
