package boards

import (
	"RPO_back/internal/pkg/auth/repository"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func DeleteBoardHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := repository.GetUserId(w, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	params := mux.Vars(r)
	boardId := params["boardId"]

	if !strings.HasPrefix(boardId, "board_") || len(boardId) < 7 {
		http.Error(w, "board id should start with 'board_'", http.StatusBadRequest)
		return
	}
	boardId = strings.TrimPrefix(boardId, "board_")

	fmt.Printf("Board id: %s\n", boardId)
	fmt.Printf("UserId id: %d\n", userId)

	db, err := repository.GetDbConnection()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Проверяем, что пользователь имеет права на удаление доски
	var isAdmin bool
	err = db.QueryRow(context.Background(), `
		SELECT is_admin
		FROM User_to_Board
		WHERE u_id = $1 AND b_id = $2
	`, userId, boardId).Scan(&isAdmin)
	if err != nil || !isAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Удаляем записи в таблице User_to_Board
	_, err = db.Exec(context.Background(), `
		DELETE FROM User_to_Board
		WHERE b_id = $1
	`, boardId)
	if err != nil {
		http.Error(w, "Failed to delete related users", http.StatusInternalServerError)
		return
	}

	// Удаляем доску
	_, err = db.Exec(context.Background(), `
		DELETE FROM Board
		WHERE b_id = $1
	`, boardId)
	if err != nil {
		http.Error(w, "Failed to delete board", http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
