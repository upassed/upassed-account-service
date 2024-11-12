package teacher

import (
	"context"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/middleware/amqp"
	loggingMiddleware "github.com/upassed/upassed-account-service/internal/middleware/amqp/logging"
	"github.com/upassed/upassed-account-service/internal/middleware/amqp/recovery"
	requestidMiddleware "github.com/upassed/upassed-account-service/internal/middleware/amqp/request_id"
	requestid "github.com/upassed/upassed-account-service/internal/middleware/common/request_id"
	"github.com/upassed/upassed-account-service/internal/tracing"
	"github.com/wagslane/go-rabbitmq"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"log/slog"
)

func (client *rabbitClient) CreateQueueConsumer() func(d rabbitmq.Delivery) rabbitmq.Action {
	baseHandler := func(ctx context.Context, delivery rabbitmq.Delivery) rabbitmq.Action {
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
			log.Error("unable to convert message body to teacher create request", logging.Error(err))
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

	handlerWithMiddleware := amqp.ChainMiddleware(
		baseHandler,
		requestidMiddleware.Middleware(),
		loggingMiddleware.Middleware(client.log),
		recovery.Middleware(client.log),
	)

	return func(d rabbitmq.Delivery) (action rabbitmq.Action) {
		ctx := context.Background()
		return handlerWithMiddleware(ctx, d)
	}
}
