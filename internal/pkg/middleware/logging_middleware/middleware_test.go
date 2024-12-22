package logging_middleware

import (
	"RPO_back/internal/pkg/utils/logging"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggingMiddleware(t *testing.T) {

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := logging.GetRequestID(r.Context()) + 1

		assert.NotZero(t, rid)
		w.WriteHeader(http.StatusOK)
	})

	handlerToTest := LoggingMiddleware(nextHandler)
	server := httptest.NewServer(handlerToTest)
	defer server.Close()

	resp, err := http.Get(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// type testWriter struct {
// 	content []byte
// }

// func (tw *testWriter) Write(p []byte) (n int, err error) {
// 	tw.content = append(tw.content, p...)
// 	return len(p), nil
// }

// func (tw *testWriter) String() string {
// 	return string(tw.content)
// }
