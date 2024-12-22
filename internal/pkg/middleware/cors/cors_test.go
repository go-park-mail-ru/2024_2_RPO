package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCorsMiddleware(t *testing.T) {
	handler := CorsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fail()
	}
}
