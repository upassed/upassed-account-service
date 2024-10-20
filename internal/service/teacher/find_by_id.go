package teacher

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/async"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/middleware"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"google.golang.org/grpc/codes"
	"log/slog"
	"reflect"
	"runtime"
)

var (
	errFindTeacherByIDDeadlineExceeded = errors.New("find teacher by id deadline exceeded")
)

func (service *teacherServiceImpl) FindByID(ctx context.Context, teacherID uuid.UUID) (business.Teacher, error) {
	op := runtime.FuncForPC(reflect.ValueOf(service.FindByID).Pointer()).Name()

	log := service.log.With(
		slog.String("op", op),
		slog.Any("teacherID", teacherID),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	log.Info("started finding teacher by id")
	timeout := service.cfg.GetEndpointExecutionTimeout()
	foundTeacher, err := async.ExecuteWithTimeout(ctx, timeout, func(ctx context.Context) (business.Teacher, error) {
		foundTeacher, err := service.repository.FindByID(ctx, teacherID)
		if err != nil {
			return business.Teacher{}, handling.Process(err)
		}

		return ConvertToServiceTeacher(foundTeacher), nil
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Error("find teacher by id deadline exceeded")
			return business.Teacher{}, handling.Wrap(errFindTeacherByIDDeadlineExceeded, handling.WithCode(codes.DeadlineExceeded))
		}

		log.Error("error while finding teacher by id", logging.Error(err))
		return business.Teacher{}, handling.Wrap(err)
	}

	log.Info("teacher successfully found by id")
	return foundTeacher, nil
}
