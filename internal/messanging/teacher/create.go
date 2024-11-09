package teacher

import (
	"context"
	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/middleware/grpc/requestid"
	"github.com/upassed/upassed-account-service/internal/tracing"
	"github.com/wagslane/go-rabbitmq"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"log/slog"
)

func (client *rabbitClient) CreateQueueConsumer() func(d rabbitmq.Delivery) rabbitmq.Action {
	return func(delivery rabbitmq.Delivery) rabbitmq.Action {
		requestID := uuid.New().String()
		ctx := context.WithValue(context.Background(), requestid.ContextKey, requestID)

		log := logging.Wrap(client.log,
			logging.WithOp(client.CreateQueueConsumer),
			logging.WithCtx(ctx),
		)

		log.Info("consumed teacher create message", slog.String("messageBody", string(delivery.Body)))
		spanContext, span := otel.Tracer(client.cfg.Tracing.TeacherTracerName).Start(ctx, "teacher#Create")
		span.SetAttributes(attribute.String(string(requestid.ContextKey), requestid.GetRequestIDFromContext(ctx)))
		defer span.End()

		log.Info("converting message body to teacher create request")
		request, err := ConvertToTeacherCreateRequest(delivery.Body)
		if err != nil {
			log.Error("unable to convert message body to techer create request", logging.Error(err))
			tracing.SetSpanError(span, err)
			return rabbitmq.NackDiscard
		}

		span.SetAttributes(attribute.String("username", request.Username))
		log.Info("validating teacher create request")
		if err := request.Validate(); err != nil {
			log.Error("teacher create request is invalid", logging.Error(err))
			tracing.SetSpanError(span, err)
			return rabbitmq.NackDiscard
		}

		log.Info("creating a teacher")
		response, err := client.service.Create(spanContext, ConvertToTeacher(request))
		if err != nil {
			log.Error("unable to create teacher", logging.Error(err))
			tracing.SetSpanError(span, err)
			return rabbitmq.NackDiscard
		}

		log.Info("successfully created teacher", slog.Any("createdTeacherID", response.CreatedTeacherID))
		return rabbitmq.Ack
	}
}
