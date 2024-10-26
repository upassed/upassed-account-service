package teacher

import (
	"context"
	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/middleware"
	"github.com/wagslane/go-rabbitmq"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"log/slog"
)

func (client *rabbitClient) CreateQueueConsumer() func(d rabbitmq.Delivery) rabbitmq.Action {
	return func(delivery rabbitmq.Delivery) rabbitmq.Action {
		requestID := uuid.New().String()
		ctx := context.WithValue(context.Background(), middleware.RequestIDKey, requestID)

		log := logging.Wrap(client.log,
			logging.WithOp(client.CreateQueueConsumer),
			logging.WithCtx(ctx),
		)

		log.Info("consumed teacher create message", slog.String("messageBody", string(delivery.Body)))

		spanContext, span := otel.Tracer(client.cfg.Tracing.TeacherTracerName).Start(ctx, "teacher#Create")
		span.SetAttributes(attribute.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)))
		defer span.End()

		request, err := ConvertToTeacherCreateRequest(delivery.Body)
		if err != nil {
			return rabbitmq.NackDiscard
		}

		if err := request.Validate(); err != nil {
			return rabbitmq.NackDiscard
		}

		response, err := client.service.Create(spanContext, ConvertToTeacher(request))
		if err != nil {
			return rabbitmq.NackDiscard
		}

		log.Info("successfully created teacher", slog.Any("createdTeacherID", response.CreatedTeacherID))
		return rabbitmq.Ack

	}
}
