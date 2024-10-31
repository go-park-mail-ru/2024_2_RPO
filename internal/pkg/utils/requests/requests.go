package requests

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
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

// GetIdFromRequest получает id из строки с префиксом (e.g. 'board_')
func GetIdFromRequest(r *http.Request, requestVarName string, prefix string) (int, error) {
	vars := mux.Vars(r)
	rawID, isExist := vars[requestVarName]
	if !isExist {
		return 0, errors.New("there is no such parameter")
	}

	IDWithoutPrefix, found := strings.CutPrefix(rawID, prefix)
	if !found {
		return 0, errors.New("error in the parameters")
	}

	resultId, err := strconv.Atoi(IDWithoutPrefix)
	if err != nil {
		return 0, errors.New("failed to convert to Int")
	}

	return resultId, nil
}
