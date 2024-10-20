package group

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/async"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/middleware"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"google.golang.org/grpc/codes"
	"log/slog"
	"reflect"
	"runtime"
)

var (
	errFindStudentsInGroupDeadlineExceeded = errors.New("find students in group deadline exceeded")
)

func (service *groupServiceImpl) FindStudentsInGroup(ctx context.Context, groupID uuid.UUID) ([]business.Student, error) {
	op := runtime.FuncForPC(reflect.ValueOf(service.FindStudentsInGroup).Pointer()).Name()

	log := service.log.With(
		slog.String("op", op),
		slog.Any("groupID", groupID),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	log.Info("started searching students in group")
	timeout := service.cfg.GetEndpointExecutionTimeout()
	foundStudents, err := async.ExecuteWithTimeout(ctx, timeout, func(ctx context.Context) ([]domain.Student, error) {
		return service.repository.FindStudentsInGroup(ctx, groupID)
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Error("students in group searching deadline exceeded")
			return make([]business.Student, 0), handling.Wrap(errFindStudentsInGroupDeadlineExceeded, handling.WithCode(codes.DeadlineExceeded))
		}

		log.Error("error while searching students in group", logging.Error(err))
		return make([]business.Student, 0), handling.Process(err)
	}

	log.Info("successfully found students in group", slog.Int("studentsCount", len(foundStudents)))
	return ConvertToServiceStudents(foundStudents), nil
}
