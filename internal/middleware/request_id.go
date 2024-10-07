package middleware

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type contextKey string

const RequestIDKey = contextKey("requestID")

func RequestIDUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		requestID := uuid.New().String()
		ctx = context.WithValue(ctx, RequestIDKey, requestID)
		return handler(ctx, req)
	}
}

func GetRequestIDFromContext(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(RequestIDKey).(string)
	return requestID, ok
}
