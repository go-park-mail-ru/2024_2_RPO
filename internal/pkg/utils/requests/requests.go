package requests

import (
	"encoding/json"
	"io"
	"net/http"
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
