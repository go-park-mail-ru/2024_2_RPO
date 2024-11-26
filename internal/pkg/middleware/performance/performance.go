package performance

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

type HTTPPerformanceMiddleware struct {
	Hits        *prometheus.CounterVec
	serviceName string
	Times       *prometheus.HistogramVec
	Errors      *prometheus.CounterVec
}

func CreateHTTPPerformanceMiddleware(serviceName string) (*HTTPPerformanceMiddleware, error) {
	hits := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: serviceName + "http_hits",
		Help: "Number of HTTP calls",
		ConstLabels: prometheus.Labels{
			"serviceName": serviceName,
		},
	}, []string{"method", "status"})

	times := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: serviceName + "http_times",
		Help: "Histogram of http call times",
		ConstLabels: prometheus.Labels{
			"serviceName": serviceName,
		},
		Buckets: prometheus.DefBuckets, // возожно не понадобится, подумаем...
	}, []string{"method", "status"})

	errors := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: serviceName + "http_errors",
		Help: "Number of http errors",
		ConstLabels: prometheus.Labels{
			"serviceName": serviceName,
		},
	}, []string{"method", "status"})

	if err := prometheus.Register(hits); err != nil {
		return nil, fmt.Errorf("create register hits (http): %w", err)
	}

	if err := prometheus.Register(times); err != nil {
		return nil, fmt.Errorf("create register times (http): %w", err)
	}

	if err := prometheus.Register(errors); err != nil {
		return nil, fmt.Errorf("create register errors (http): %w", err)
	}

	return &HTTPPerformanceMiddleware{
		Hits:        hits,
		serviceName: serviceName,
		Times:       times,
		Errors:      errors,
	}, nil
}

func (*HTTPPerformanceMiddleware) Middleware(next http.Handler) http.Handler {
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
