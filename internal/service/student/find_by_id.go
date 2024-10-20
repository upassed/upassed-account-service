package student

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/async"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/middleware"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"google.golang.org/grpc/codes"
	"log/slog"
	"reflect"
	"runtime"
)

var (
	errFindStudentByIDDeadlineExceeded = errors.New("find student by id deadline exceeded")
)

func (service *studentServiceImpl) FindByID(ctx context.Context, studentID uuid.UUID) (business.Student, error) {
	op := runtime.FuncForPC(reflect.ValueOf(service.FindByID).Pointer()).Name()

	log := service.log.With(
		slog.String("op", op),
		slog.Any("studentID", studentID),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	timeout := service.cfg.GetEndpointExecutionTimeout()
	foundStudent, err := async.ExecuteWithTimeout(ctx, timeout, func(ctx context.Context) (business.Student, error) {
		log.Debug("started finding student by id")
		foundStudent, err := service.studentRepository.FindByID(ctx, studentID)
		if err != nil {
			return business.Student{}, handling.Process(err)
		}

		log.Debug("student successfully found by id")
		return ConvertToServiceStudent(foundStudent), nil
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return business.Student{}, handling.Wrap(errFindStudentByIDDeadlineExceeded, handling.WithCode(codes.DeadlineExceeded))
		}

		return business.Student{}, handling.Wrap(err)
	}

	return foundStudent, nil
}
