package requests

import (
	"RPO_back/internal/models"
	"bytes"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetRequestData(t *testing.T) {
	columnJSON := `{"id": 1, "title": "Test Column"}`
	req, err := http.NewRequest("POST", "/api/228", bytes.NewBufferString(columnJSON))
	assert.NoError(t, err)

	var column models.Column
	err = GetRequestData(req, &column)
	assert.NoError(t, err)
}

func TestGetIDFromRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/boards/board_123", nil)
	vars := map[string]string{
		"boardID": "board_123",
	}
	req = mux.SetURLVars(req, vars)

	_, err := GetIDFromRequest(req, "boardID", "board_")
	assert.NoError(t, err)
}

func TestGetIDFromRequestWithPrefixError(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/boards/123", nil)
	vars := map[string]string{
		"boardID": "123",
	}
	req = mux.SetURLVars(req, vars)

	_, err := GetIDFromRequest(req, "boardID", "board_")
	assert.Error(t, err)
}
