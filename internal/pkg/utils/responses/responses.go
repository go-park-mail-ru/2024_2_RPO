package responses

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/utils/logging"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// Вернуть ответ с указанным статусом
func DoBadResponseAndLog(r *http.Request, w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := models.BadResponse{
		Status: statusCode,
		Text:   message,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)

	logging.Error(r.Context(), "Bad response with status ", statusCode, " and message ", message)
}

func DoEmptyOkResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"status\":200,\"text\":\"success\"}"))
}

func DoJSONResponse(r *http.Request, w http.ResponseWriter, responseData interface{}, successStatusCode int) {
	body, err := json.Marshal(responseData)
	if err != nil {
		DoBadResponseAndLog(r, w, 500, "internal error")
		log.Error(fmt.Errorf("error in marshalling response body: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))

	w.WriteHeader(successStatusCode)
	_, _ = w.Write(body)
}

// ResponseErrorAndLog принимает ошибку, которая пришла из usecase, и делает ответ
// в соответствии с типом ошибки. Также он делает запись в log с типом WARN, если
// ошибка стандартная, и ERRO, если это 500.
//
// Типичная запись в логе: `UserToBoard: Not found`.
// В данном случае префикс - `UserToBoard`, двоеточие мы поставим сами.
//
// Поддерживаемые типы ошибок: 404, 403, 500
func ResponseErrorAndLog(r *http.Request, w http.ResponseWriter, err error, prefix string) {
	if errors.Is(err, errs.ErrNotFound) {
		DoBadResponseAndLog(r, w, http.StatusNotFound, "not found")
		log.Warn(prefix, ": ", err)
		return
	}
	if errors.Is(err, errs.ErrNotPermitted) {
		DoBadResponseAndLog(r, w, http.StatusForbidden, "forbidden")
		log.Warn(prefix, ": ", err)
		return
	}
	if errors.Is(err, errs.ErrValidation) {
		DoBadResponseAndLog(r, w, http.StatusBadRequest, err.Error())
		log.Warn(prefix, ": ", err)
		return
	}
	log.Error(prefix, ": ", err)
	DoBadResponseAndLog(r, w, http.StatusInternalServerError, "internal error")
}
