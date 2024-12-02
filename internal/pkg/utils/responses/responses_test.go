package responses_test

import (
	"RPO_back/internal/pkg/utils/responses"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDoBadResponse(t *testing.T) {
	recorder := httptest.NewRecorder()
	statusCode := http.StatusBadRequest
	message := "test error message"

	responses.DoBadResponseAndLog(recorder, statusCode, message)

	result := recorder.Result()
	defer result.Body.Close()

	if result.StatusCode != statusCode {
		t.Errorf("expected status %d, got %d", statusCode, result.StatusCode)
	}

	expectedContentType := "application/json"
	if result.Header.Get("Content-Type") != expectedContentType {
		t.Errorf("expected Content-Type %s, got %s", expectedContentType, result.Header.Get("Content-Type"))
	}

	expectedBody := `{"status":400,"text":"test error message"}`
	responseData := recorder.Body.String()
	if responseData != expectedBody {
		t.Errorf("expected body %s, got %s", expectedBody, responseData)
	}
}
