package teacher

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/async"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logging"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"github.com/upassed/upassed-account-service/internal/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
)

var (
	errFindTeacherByIDDeadlineExceeded = errors.New("find teacher by id deadline exceeded")
)

func (service *teacherServiceImpl) FindByID(ctx context.Context, teacherID uuid.UUID) (*business.Teacher, error) {
	spanContext, span := otel.Tracer(service.cfg.Tracing.TeacherTracerName).Start(ctx, "teacherService#FindByID")
	span.SetAttributes(attribute.String("id", teacherID.String()))
	defer span.End()

	log := logging.Wrap(service.log,
		logging.WithOp(service.FindByID),
		logging.WithCtx(ctx),
		logging.WithAny("teacherID", teacherID),
	)

	log.Info("started finding teacher by id")
	timeout := service.cfg.GetEndpointExecutionTimeout()
	foundTeacher, err := async.ExecuteWithTimeout(spanContext, timeout, func(ctx context.Context) (*business.Teacher, error) {
		log.Info("finding teacher data")
		foundTeacher, err := service.repository.FindByID(ctx, teacherID)
		if err != nil {
			log.Error("error while finding teacher data by id", logging.Error(err))
			tracing.SetSpanError(span, err)
			return nil, handling.Process(err)
		}

		return ConvertToServiceTeacher(foundTeacher), nil
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Error("find teacher by id deadline exceeded")
			tracing.SetSpanError(span, err)
			return nil, handling.Wrap(errFindTeacherByIDDeadlineExceeded, handling.WithCode(codes.DeadlineExceeded))
		}

		log.Error("error while finding teacher by id", logging.Error(err))
		tracing.SetSpanError(span, err)
		return nil, handling.Wrap(err)
	}

	log.Info("teacher successfully found by id")
	return foundTeacher, nil
}
