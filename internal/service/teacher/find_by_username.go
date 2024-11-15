package teacher

import (
	"context"
	"errors"
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
	errFindTeacherByUsernameDeadlineExceeded = errors.New("find teacher by username deadline exceeded")
)

func (service *serviceImpl) FindByUsername(ctx context.Context, teacherUsername string) (*business.Teacher, error) {
	spanContext, span := otel.Tracer(service.cfg.Tracing.TeacherTracerName).Start(ctx, "teacherService#FindByUsername")
	span.SetAttributes(attribute.String("teacherUsername", teacherUsername))
	defer span.End()

	log := logging.Wrap(service.log,
		logging.WithOp(service.FindByUsername),
		logging.WithCtx(ctx),
		logging.WithAny("teacherUsername", teacherUsername),
	)

	log.Info("started finding teacher by username")
	timeout := service.cfg.GetEndpointExecutionTimeout()
	foundTeacher, err := async.ExecuteWithTimeout(spanContext, timeout, func(ctx context.Context) (*business.Teacher, error) {
		log.Info("finding teacher data")
		foundTeacher, err := service.repository.FindByUsername(ctx, teacherUsername)
		if err != nil {
			log.Error("error while finding teacher data by username", logging.Error(err))
			tracing.SetSpanError(span, err)
			return nil, handling.Process(err)
		}

		return ConvertToServiceTeacher(foundTeacher), nil
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Error("find teacher by username deadline exceeded")
			tracing.SetSpanError(span, err)
			return nil, handling.Wrap(errFindTeacherByUsernameDeadlineExceeded, handling.WithCode(codes.DeadlineExceeded))
		}

		log.Error("error while finding teacher by username", logging.Error(err))
		tracing.SetSpanError(span, err)
		return nil, handling.Wrap(err)
	}

	log.Info("teacher successfully found by username")
	return foundTeacher, nil
}
