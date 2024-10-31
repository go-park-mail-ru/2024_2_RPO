package csrf

import (
	"RPO_back/internal/pkg/utils/encrypt"
	"RPO_back/internal/pkg/utils/responses"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func checkCSRFToken(w http.ResponseWriter, r *http.Request) bool {
	csrfHeader := r.Header.Get("X-CSRF-Token")
	if csrfHeader == "" {
		log.Warn(r.URL.Path, "no X-CSRF-Token header")
		responses.DoBadResponse(w, http.StatusForbidden, "no X-CSRF-Token header")
		return false
	}

	csrfCookie, err := r.Cookie("csrf_token")
	if err != nil {
		log.Warn(r.URL.Path, "no csrf cookie: ", err)
		responses.DoBadResponse(w, http.StatusForbidden, "no csrf cookie")
		return false
	}

	if csrfCookie.Value != csrfHeader {
		log.Warn(r.URL.Path, "tokens in cookie and header are different")
		responses.DoBadResponse(w, http.StatusForbidden, "tokens in cookie and header are different")
		return false
	}

	return true
}

func SetCSRFToken(w http.ResponseWriter) {
	csrfToken := encrypt.GenerateCSRFToken()

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		MaxAge:   int(1 * time.Hour.Seconds()),
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})
	w.Header().Set("X-Csrf-Token", csrfToken)
}

func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if (r.Method == http.MethodPatch) || (r.Method == http.MethodPost) ||
			(r.Method == http.MethodPut) || (r.Method == http.MethodDelete) {

			if !checkCSRFToken(w, r) {
				return
			}
		}
		SetCSRFToken(w)

		log.Info("CSRF check PASS")
		next.ServeHTTP(w, r)
	})
}
