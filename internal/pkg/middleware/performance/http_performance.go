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
		Name: serviceName + "_http_hits",
		Help: "Number of HTTP calls",
		ConstLabels: prometheus.Labels{
			"serviceName": serviceName,
		},
	}, []string{"method", "status"})

	times := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: serviceName + "_http_times",
		Help: "Histogram of http call times",
		ConstLabels: prometheus.Labels{
			"serviceName": serviceName,
		},
		Buckets: prometheus.DefBuckets, // возожно не понадобится, подумаем...
	}, []string{"method", "status"})

	errors := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: serviceName + "_http_errors",
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

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func CreateResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (mh *responseWriter) WriteHeader(code int) {
	mh.statusCode = code
	mh.ResponseWriter.WriteHeader(code)
}

func (mh *HTTPPerformanceMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()

		rw := CreateResponseWriter(w)
		next.ServeHTTP(rw, r)

		path, err := mux.CurrentRoute(r).GetPathTemplate()
		if err != nil || path == "" {
			path = "/unknown-route"
		}

		method := r.Method
		status := fmt.Sprintf("%d", rw.statusCode)

		mh.Hits.WithLabelValues(method, status).Inc()

		duration := time.Since(timeStart).Seconds()
		mh.Times.WithLabelValues(method, status).Observe(duration)

		mh.Errors.WithLabelValues(method, status).Inc()

	})
}
