package logging_middleware

import (
	"context"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"
)

type contextKey string

const rIDKey = contextKey("requestID")

var (
	rIDCounter uint64 = 1
	mu         sync.Mutex
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		rIDCounter++
		rID := rIDCounter
		mu.Unlock()

		ctx := context.WithValue(r.Context(), rIDKey, rID)
		log.Infof("Запрос: %s %s, RequestID: %d", r.Method, r.RequestURI, rID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
