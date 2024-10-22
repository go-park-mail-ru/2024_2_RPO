package boards

import (
	"RPO_back/database"
	"RPO_back/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func CreateBoardHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := database.GetUserId(w, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.CreateBoardRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		//TODO JSON
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Description == "" {
		http.Error(w, "Bad Request: Missing required fields", http.StatusBadRequest)
		return
	}

	var board models.Board
	fmt.Printf("Name: %s\n", req.Name)
	board.Name = req.Name
	board.Description = req.Description
	board.Background = "red" //TODO пересмотреть
	board.OwnerID = userId

	db, err := database.GetDbConnection()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}

	var boardId int
	err = db.QueryRow(context.Background(), `
		INSERT INTO Board (description, name, created_by, updated_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		RETURNING b_id`,
		board.Description, board.Name, board.OwnerID).Scan(&boardId)
	if err != nil {
		http.Error(w, "Failed to insert board", http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}
	board.ID = boardId

	_, err = db.Exec(context.Background(), `
		INSERT INTO User_to_Board (u_id, b_id, added_at, updated_at, can_edit, can_share, can_invite_members, is_admin, added_by, updated_by)
		VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, TRUE, TRUE, TRUE, TRUE, NULL, NULL)`,
		board.OwnerID, boardId)
	if err != nil {
		http.Error(w, "Failed to insert user to board", http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(board)
}
