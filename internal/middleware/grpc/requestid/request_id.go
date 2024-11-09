package requestid

import (
	"context"
	requestid "github.com/upassed/upassed-account-service/internal/middleware/common/request_id"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func MiddlewareInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		requestID := uuid.New().String()
		ctx = context.WithValue(ctx, requestid.ContextKey, requestID)
		return handler(ctx, req)
	}
}
