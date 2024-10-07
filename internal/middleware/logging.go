package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggingMiddlewareInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		const op = "middleware.LoggingMiddlewareInterceptor()"

		startTime := time.Now()
		requestID, _ := GetRequestIDFromContext(ctx)
		resp, err := handler(ctx, req)
		elapsedTime := time.Since(startTime)

		log := log.With(
			slog.String("op", op),
		)

		log.Info("handled gRPC request",
			slog.String("requestID", requestID),
			slog.String("method", info.FullMethod),
			slog.String("duration", fmt.Sprintf("%.2f ms", elapsedTime.Seconds()*1000)),
			slog.String("status", status.Code(err).String()),
		)

		return resp, err
	}
}
