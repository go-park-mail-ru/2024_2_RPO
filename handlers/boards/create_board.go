package boards

import (
	"RPO_back/database"
	"RPO_back/models"
	"encoding/json"
	"log"
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
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Description == "" {
		http.Error(w, "Bad Request: Missing required fields", http.StatusBadRequest)
		return
	}

	var board models.Board
	board.Name = req.Name
	board.Description = req.Description
	board.Background = req.Background
	board.OwnerID = userId

	db, err := database.GetDbConnection()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := `INSERT INTO boards (name, description, background, owner_id) VALUES ($1, $2, $3, $4) RETURNING id`
	err = db.QueryRow(query, board.Name, board.Description, board.Background, board.OwnerID).Scan(&board.ID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error inserting board:", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(board)
}
