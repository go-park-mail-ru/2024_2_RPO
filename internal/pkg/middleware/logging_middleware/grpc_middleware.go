package logging_middleware

import (
	"context"

	"RPO_back/internal/pkg/utils/logging"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type GrpcLogMiddleware struct {
	grpcLogger *log.Logger
}

func CreateGrpcLogMiddleware(gL *log.Logger) *GrpcLogMiddleware {
	return &GrpcLogMiddleware{
		grpcLogger: gL,
	}
}

func (glm *GrpcLogMiddleware) InterceptorLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	logging.Infof(ctx, "gRPC method called: %s", info.FullMethod)
	logging.Debugf(ctx, "Request payload: %+v", req)

	resp, err := handler(ctx, req)

	if err != nil {
		logging.Errorf(ctx, "Error handling gRPC method %s: %v", info.FullMethod, err)
	} else {
		logging.Debugf(ctx, "Response payload for %s: %+v", info.FullMethod, resp)
	}

	return resp, err
}
