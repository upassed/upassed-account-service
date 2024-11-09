package recovery

import (
	"context"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/messanging"
	"github.com/upassed/upassed-account-service/internal/middleware/amqp"
	"github.com/upassed/upassed-account-service/internal/middleware/grpc/requestid"
	"github.com/wagslane/go-rabbitmq"
	"log/slog"
)

func Middleware(log *slog.Logger) amqp.Middleware {
	return func(ctx context.Context, next messanging.HandlerWithContext) messanging.HandlerWithContext {
		return func(ctx context.Context, d rabbitmq.Delivery) (action rabbitmq.Action) {
			defer func() {
				if r := recover(); r != nil {
					log := logging.Wrap(log, logging.WithOp(Middleware))

					log.Error("panic recovered",
						slog.String("requestID", requestid.GetRequestIDFromContext(ctx)),
						slog.Any("message", r),
						slog.String("routingKey", d.RoutingKey),
					)

					action = rabbitmq.NackDiscard
				}
			}()

			return next(ctx, d)
		}
	}
}
