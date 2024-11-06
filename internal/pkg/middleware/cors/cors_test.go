package cors

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCorsMiddleware(t *testing.T) {
	os.Setenv("CORS_ORIGIN", "https://example.com")

	// Функция, которую будет оборачивать middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		method          string
		expectedStatus  int
		expectedHeaders map[string]string
	}{
		{
			method:         "GET",
			expectedStatus: http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "https://example.com",
				"Access-Control-Allow-Methods":     "GET, POST, OPTIONS, PUT, PATCH, DELETE",
				"Access-Control-Allow-Headers":     "Content-Type, Authorization, X-CSRF-Token",
				"Access-Control-Allow-Credentials": "true",
				"Access-Control-Expose-Headers":    "X-CSRF-Token",
				"Content-Security-Policy":          csp,
			},
		},
		{
			method:         "OPTIONS",
			expectedStatus: http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "https://example.com",
				"Access-Control-Allow-Methods":     "GET, POST, OPTIONS, PUT, PATCH, DELETE",
				"Access-Control-Allow-Headers":     "Content-Type, Authorization, X-CSRF-Token",
				"Access-Control-Allow-Credentials": "true",
				"Access-Control-Expose-Headers":    "X-CSRF-Token",
				"Content-Security-Policy":          csp,
			},
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest(test.method, "http://example.com", nil)
		rr := httptest.NewRecorder()

		CorsMiddleware(handler).ServeHTTP(rr, req)

		if status := rr.Code; status != test.expectedStatus {
			t.Errorf("handler returned wrong status code: got %v want %v", status, test.expectedStatus)
		}

		for header, expectedValue := range test.expectedHeaders {
			if value := rr.Header().Get(header); value != expectedValue {
				t.Errorf("header %s = %v, want %v", header, value, expectedValue)
			}
		}
	}
}
