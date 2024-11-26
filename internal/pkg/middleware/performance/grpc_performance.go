package performance

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func mapStatusCodes(err error) int {
	switch status.Code(err) {
	case codes.OK:
		return http.StatusOK
	case codes.NotFound:
		return http.StatusNotFound
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func (mg *GRPCPerformanceMiddleware) GRPCMetricsInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	start := time.Now()
	resp, err = handler(ctx, req)

	statusCode := mapStatusCodes(err)
	grpcCode := status.Code(err).String()

	mg.Hits.WithLabelValues(info.FullMethod, fmt.Sprintf("%d", statusCode)).Inc()
	if err != nil {
		mg.Errors.WithLabelValues(info.FullMethod, fmt.Sprintf("%d", statusCode), grpcCode).Inc()
	}

	duration := time.Since(start).Seconds()
	mg.Times.WithLabelValues(info.FullMethod, fmt.Sprintf("%d", statusCode)).Observe(duration)

	return resp, err
}
