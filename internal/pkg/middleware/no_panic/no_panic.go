package no_panic

import (
	"RPO_back/internal/pkg/utils/responses"
	"net/http"
	"runtime/debug"

	log "github.com/sirupsen/logrus"
)

func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Паника для прерывания запроса обрабатывается в mux'e
				if err == http.ErrAbortHandler {
					log.Warn("Abort connection")
					panic(err)
				}
				log.Error("Panic: ", err)
				log.Error("Debug stack: ", string(debug.Stack()))
				responses.DoBadResponse(w, http.StatusInternalServerError, "internal error")
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}
