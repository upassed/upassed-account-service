package logging

import (
	"context"
	"fmt"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/messanging"
	"github.com/upassed/upassed-account-service/internal/middleware/amqp"
	"github.com/upassed/upassed-account-service/internal/middleware/grpc/requestid"
	"github.com/wagslane/go-rabbitmq"
	"log/slog"
	"time"
)

func Middleware(log *slog.Logger) amqp.Middleware {
	return func(ctx context.Context, next messanging.HandlerWithContext) messanging.HandlerWithContext {
		return func(ctx context.Context, d rabbitmq.Delivery) (action rabbitmq.Action) {
			log := logging.Wrap(log, logging.WithOp(Middleware))

			startTime := time.Now()
			resp := next(ctx, d)
			elapsedTime := time.Since(startTime)

			log.Info("handled amqp request",
				slog.String("requestID", requestid.GetRequestIDFromContext(ctx)),
				slog.String("routingKey", d.RoutingKey),
				slog.String("duration", fmt.Sprintf("%.2f ms", elapsedTime.Seconds()*1000)),
				slog.Any("resultAction", resp),
			)

			return resp
		}
	}
}
