package performance

import "github.com/prometheus/client_golang/prometheus"

type GRPCPerformanceMiddleware struct {
	Hits        *prometheus.CounterVec
	serviceName string
	Times       *prometheus.HistogramVec
	Errors      *prometheus.CounterVec
}

func CreateGRPCPerformanceMiddleware(serviceName string) (*GRPCPerformanceMiddleware, error) {
	panic("not implemented")
}
