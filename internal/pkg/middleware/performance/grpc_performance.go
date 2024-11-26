package performance

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type GRPCPerformanceMiddleware struct {
	Hits        *prometheus.CounterVec
	serviceName string
	Times       *prometheus.HistogramVec
	Errors      *prometheus.CounterVec
}

func CreateGRPCPerformanceMiddleware(serviceName string) (*GRPCPerformanceMiddleware, error) {
	hits := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:        serviceName + "grps_hits",
		Help:        "Number of gRPC calls",
		ConstLabels: prometheus.Labels{"serviceName": serviceName},
	},
		[]string{"method", "status"})

	times := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: serviceName + "grpc_times",
		Help: "Histogram of gRPC call times",
		ConstLabels: prometheus.Labels{
			"serviceName": serviceName,
		},
		Buckets: prometheus.DefBuckets, // возожно не понадобится, подумаем...
	},
		[]string{"method", "status"})

	errors := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: serviceName + "grpc_errors",
		Help: "Number of gRPC errors",
		ConstLabels: prometheus.Labels{
			"serviceName": serviceName,
		},
	}, []string{"method", "status"})

	if err := prometheus.Register(hits); err != nil {
		return nil, fmt.Errorf("create register hits (gRPC): %w", err)
	}

	if err := prometheus.Register(times); err != nil {
		return nil, fmt.Errorf("create register times (gRPC): %w", err)
	}

	if err := prometheus.Register(errors); err != nil {
		return nil, fmt.Errorf("create register errors (gRPC): %w", err)
	}

	return &GRPCPerformanceMiddleware{
		Hits:        hits,
		serviceName: serviceName,
		Times:       times,
		Errors:      errors,
	}, nil
}
