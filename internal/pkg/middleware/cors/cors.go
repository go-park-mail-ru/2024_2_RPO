package cors

import (
	"RPO_back/internal/pkg/config"
	"net/http"

	"github.com/sirupsen/logrus"
)

const csp = "default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self'; base-uri 'self'; form-action 'self'"

// Middleware для CORS
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if config.CurrentConfig == nil {
			logrus.Warn("CorsMiddleware: current config is nil")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", config.CurrentConfig.CorsOriging)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Expose-Headers", "X-CSRF-Token")
		w.Header().Set("Content-Security-Policy", csp)

		// Для preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
