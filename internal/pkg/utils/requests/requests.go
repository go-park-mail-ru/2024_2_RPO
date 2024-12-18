package requests

import (
	"RPO_back/internal/pkg/middleware/session"
	"RPO_back/internal/pkg/utils/responses"
	"RPO_back/internal/pkg/utils/validate"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var (
	uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
)

// Получить тело запроса. Прочитать данные из json, разместить в структуру
func GetRequestData(r *http.Request, requestData json.Unmarshaler) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if err := requestData.UnmarshalJSON(body); err != nil {
		return err
	}

	err = validate.Validate(r.Context(), requestData)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// GetIDFromRequest получает id из строки с префиксом (e.g. 'board_')
func GetIDFromRequest(r *http.Request, requestVarName string, prefix string) (int64, error) {
	vars := mux.Vars(r)
	rawID, isExist := vars[requestVarName]
	if !isExist {
		return 0, errors.New("there is no such parameter: " + requestVarName)
	}

	IDWithoutPrefix, found := strings.CutPrefix(rawID, prefix)
	if !found {
		return 0, errors.New("error in the parameters")
	}

	resultID, err := strconv.ParseInt(IDWithoutPrefix, 10, 64)
	if err != nil {
		return 0, errors.New("failed to convert to Int")
	}

	return resultID, nil
}

// GetUUIDFromRequest получает UUID из параметра запроса.
func GetUUIDFromRequest(r *http.Request, requestVarName string) (string, error) {
	vars := mux.Vars(r)
	rawID, isExist := vars[requestVarName]
	if !isExist {
		return "", errors.New("there is no such parameter: " + requestVarName)
	}

	if !uuidRegex.MatchString(rawID) {
		return "", errors.New("invalid UUID format")
	}

	return rawID, nil
}

// GetUserIDOrFail достаёт UserID из запроса. Если его нет, возвращает 401 и пишет в лог
func GetUserIDOrFail(w http.ResponseWriter, r *http.Request, prefix string) (userID int64, ok bool) {
	userID, ok = session.UserIDFromContext(r.Context())

	if !ok {
		responses.DoBadResponseAndLog(r, w, http.StatusUnauthorized, "unauthorized")
		log.Warn(prefix, ": unauthorized")
		return 0, false
	}
	return userID, true
}
