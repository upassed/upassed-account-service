package student

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
	errFindStudentByUsernameDeadlineExceeded = errors.New("find student by username deadline exceeded")
)

func (service *serviceImpl) FindByUsername(ctx context.Context, studentUsername string) (*business.Student, error) {
	spanContext, span := otel.Tracer(service.cfg.Tracing.StudentTracerName).Start(ctx, "studentService#FindByUsername")
	span.SetAttributes(attribute.String("studentUsername", studentUsername))
	defer span.End()

	log := logging.Wrap(service.log,
		logging.WithOp(service.FindByUsername),
		logging.WithCtx(ctx),
		logging.WithAny("studentUsername", studentUsername),
	)

	log.Info("started searching student by username")
	timeout := service.cfg.GetEndpointExecutionTimeout()
	foundStudent, err := async.ExecuteWithTimeout(spanContext, timeout, func(ctx context.Context) (*business.Student, error) {
		log.Info("finding student data by username")
		foundStudent, err := service.studentRepository.FindByUsername(ctx, studentUsername)
		if err != nil {
			log.Error("unable to find student data by username", logging.Error(err))
			tracing.SetSpanError(span, err)
			return nil, handling.Process(err)
		}

		return ConvertToServiceStudent(foundStudent), nil
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Error("searching student by username deadline exceeded")
			tracing.SetSpanError(span, err)
			return nil, handling.Wrap(errFindStudentByUsernameDeadlineExceeded, handling.WithCode(codes.DeadlineExceeded))
		}

		log.Error("error while searching student by username", logging.Error(err))
		tracing.SetSpanError(span, err)
		return nil, handling.Wrap(err)
	}

	log.Info("student successfully found by username")
	return foundStudent, nil
}
