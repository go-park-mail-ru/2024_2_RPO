package csrf

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockGenerateCSRFToken() string {
	return "mockedCSRFToken"
}

func TestGetRequestSetsHeaderAndCookie(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	rr := httptest.NewRecorder()

	CSRFMiddleware(handler).ServeHTTP(rr, req)

	csrfTokenHeader := rr.Header().Get("X-CSRF-Token")
	cookie := rr.Result().Cookies()

	if csrfTokenHeader != cookie[0].Value {
		t.Errorf("expected X-CSRF-Token header and CSRF cookie to match, got header: '%s', cookie: '%s'", csrfTokenHeader, cookie[0].Value)
	}
}

func TestPOSTRequestWithDifferentHeaderAndCookie(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "http://example.com", nil)
	rr := httptest.NewRecorder()

	req.Header.Set("X-CSRF-Token", "invalidToken")
	req.AddCookie(&http.Cookie{
		Name:  "csrf_token",
		Value: mockGenerateCSRFToken(),
	})

	CSRFMiddleware(handler).ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status %v, got %v", http.StatusForbidden, rr.Code)
	}

	expectedMessage := `{"status":403,"text":"csrf tokens in cookie and header are different"}`
	if body := rr.Body.String(); body != expectedMessage {
		t.Errorf("expected body to contain %v, got %v", expectedMessage, body)
	}
}

func TestPOSTRequestWithValidHeaderAndCookie(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "http://example.com", nil)
	rr := httptest.NewRecorder()

	req.Header.Set("X-CSRF-Token", mockGenerateCSRFToken())
	req.AddCookie(&http.Cookie{
		Name:  "csrf_token",
		Value: mockGenerateCSRFToken(),
	})

	CSRFMiddleware(handler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}
}
