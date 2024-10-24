package responses

import (
	"RPO_back/internal/models"
	"encoding/json"
	"net/http"
)

// Вернуть ответ с указанным статусом
func DoBadResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := models.BadResponse{
		Status: statusCode,
		Text:   message,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}
