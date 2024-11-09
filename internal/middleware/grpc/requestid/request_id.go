package requestid

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type contextKey string

const ContextKey = contextKey("requestID")

func MiddlewareInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		requestID := uuid.New().String()
		ctx = context.WithValue(ctx, ContextKey, requestID)
		return handler(ctx, req)
	}
}

func GetRequestIDFromContext(ctx context.Context) string {
	requestID, ok := ctx.Value(ContextKey).(string)
	if !ok {
		return ""
	}

	return requestID
}
