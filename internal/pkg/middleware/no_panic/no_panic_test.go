package no_panic

// import (
// 	"bytes"
// 	"io"
// 	"log"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"regexp"
// 	"testing"
// )

// func TestPanicMiddleware(t *testing.T) {
// 	handlerThatPanics := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		panic("unexpected error")
// 	})

// 	handler := PanicMiddleware(handlerThatPanics)

// 	t.Run("handler recovers from panic", func(t *testing.T) {
// 		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
// 		rr := httptest.NewRecorder()

// 		handler.ServeHTTP(rr, req)

// 		if rr.Code != http.StatusInternalServerError {
// 			t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, rr.Code)
// 		}

// 		expected := `{"error":"internal error"}`
// 		if rr.Body.String() != expected {
// 			t.Errorf("Expected response body %s, got %s", expected, rr.Body.String())
// 		}

// 		// Check that the error is logged
// 		var buf bytes.Buffer
// 		log.SetOutput(io.MultiWriter(&buf, os.Stderr))
// 		defer log.SetOutput(os.Stderr)

// 		matched, err := regexp.MatchString(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\+\d{2}:\d{2} level=error msg="Panic: unexpected error"$`, buf.String())
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		if !matched {
// 			t.Errorf("Expected log output to match %q, got %q", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\+\d{2}:\d{2} level=error msg="Panic: unexpected error"$`, buf.String())
// 		}
// 	})

// 	t.Run("handler works without panic", func(t *testing.T) {
// 		handlerThatDoesNotPanic := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			w.WriteHeader(http.StatusOK)
// 		})

// 		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
// 		rr := httptest.NewRecorder()

// 		PanicMiddleware(handlerThatDoesNotPanic).ServeHTTP(rr, req)

// 		if rr.Code != http.StatusOK {
// 			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
// 		}
// 	})
// }
