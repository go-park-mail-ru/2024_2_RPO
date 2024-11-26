package performance

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

type PerformanceMiddleware struct {
	serviceName string
	Hits        *prometheus.Metric
}

func NewPerformanceMiddleware() *PerformanceMiddleware {
	return &PerformanceMiddleware{}
}

func (*PerformanceMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()

		next.ServeHTTP(w, r)
		timeElapsed := time.Now().UnixMilli() - timeStart.UnixMilli()
		pathTemplate, err := mux.CurrentRoute(r).GetPathTemplate()
		if err != nil {
			pathTemplate = "/invalid-route/"
		}

		statPath := r.Method + pathTemplate

		fmt.Printf("Time elapsed: %dms Route: %s\n", timeElapsed, statPath)
	})
}
