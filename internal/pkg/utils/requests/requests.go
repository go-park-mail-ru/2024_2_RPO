package requests

import (
	"RPO_back/internal/pkg/middleware/session"
	"RPO_back/internal/pkg/utils/responses"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Получить тело запроса. Прочитать данные из json, разместить в структуру
func GetRequestData(r *http.Request, requestData interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &requestData); err != nil {
		return err
	}

	return nil
}

// GetIDFromRequest получает id из строки с префиксом (e.g. 'board_')
func GetIDFromRequest(r *http.Request, requestVarName string, prefix string) (int, error) {
	vars := mux.Vars(r)
	rawID, isExist := vars[requestVarName]
	if !isExist {
		return 0, errors.New("there is no such parameter")
	}

	IDWithoutPrefix, found := strings.CutPrefix(rawID, prefix)
	if !found {
		return 0, errors.New("error in the parameters")
	}

	resultID, err := strconv.Atoi(IDWithoutPrefix)
	if err != nil {
		return 0, errors.New("failed to convert to Int")
	}

	return resultID, nil
}

// GetUserIDOrFail достаёт UserID из запроса. Если его нет, возвращает 401 и пишет в лог
func GetUserIDOrFail(w http.ResponseWriter, r *http.Request, prefix string) (userID int, ok bool) {
	userID, hasUserID := session.UserIDFromContext(r.Context())
	if !hasUserID {
		responses.DoBadResponse(w, http.StatusUnauthorized, "unauthorized")
		log.Warn(prefix, ": unauthorized")
		return 0, false
	}
	return userID, true
}
