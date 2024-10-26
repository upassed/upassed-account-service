package student

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/async"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logging"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
)

var (
	errFindStudentByIDDeadlineExceeded = errors.New("find student by id deadline exceeded")
)

func (service *studentServiceImpl) FindByID(ctx context.Context, studentID uuid.UUID) (*business.Student, error) {
	spanContext, span := otel.Tracer(service.cfg.Tracing.StudentTracerName).Start(ctx, "studentService#FindByID")
	defer span.End()

	log := logging.Wrap(service.log,
		logging.WithOp(service.FindByID),
		logging.WithCtx(ctx),
		logging.WithAny("studentID", studentID),
	)

	log.Info("started searching student by id")
	timeout := service.cfg.GetEndpointExecutionTimeout()
	foundStudent, err := async.ExecuteWithTimeout(spanContext, timeout, func(ctx context.Context) (*business.Student, error) {
		foundStudent, err := service.studentRepository.FindByID(ctx, studentID)
		if err != nil {
			return nil, handling.Process(err)
		}

		return ConvertToServiceStudent(foundStudent), nil
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Error("searching student by id deadline exceeded")
			return nil, handling.Wrap(errFindStudentByIDDeadlineExceeded, handling.WithCode(codes.DeadlineExceeded))
		}

		log.Error("error while searching student by id", logging.Error(err))
		return nil, handling.Wrap(err)
	}

	log.Info("student successfully found by id")
	return foundStudent, nil
}
