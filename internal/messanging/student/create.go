package student

import (
	"context"
	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/middleware"
	"github.com/wagslane/go-rabbitmq"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"log/slog"
	"reflect"
	"runtime"
)

func (client *rabbitClient) CreateQueueConsumer(log *slog.Logger) func(d rabbitmq.Delivery) rabbitmq.Action {
	op := runtime.FuncForPC(reflect.ValueOf(client.CreateQueueConsumer).Pointer()).Name()

	return func(delivery rabbitmq.Delivery) rabbitmq.Action {
		requestID := uuid.New().String()
		ctx := context.WithValue(context.Background(), middleware.RequestIDKey, requestID)

		log = log.With(
			slog.String("op", op),
			slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
		)

		log.Info("consumed student create message", slog.String("messageBody", string(delivery.Body)))
		spanContext, span := otel.Tracer(client.cfg.Tracing.StudentTracerName).Start(ctx, "student#Create")
		span.SetAttributes(attribute.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)))
		defer span.End()

		request, err := ConvertToStudentCreateRequest(delivery.Body)
		if err != nil {
			return rabbitmq.NackDiscard
		}

		if err := request.Validate(); err != nil {
			return rabbitmq.NackDiscard
		}

		response, err := client.service.Create(spanContext, ConvertToStudent(request))
		if err != nil {
			return rabbitmq.NackDiscard
		}

		log.Info("successfully created student", slog.Any("createdStudentID", response.CreatedStudentID))
		return rabbitmq.Ack
	}
}
