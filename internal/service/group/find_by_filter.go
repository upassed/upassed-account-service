package group

import (
	"context"
	"errors"
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
	errFindGroupsByFilterDeadlineExceeded = errors.New("find groups by filter timeout exceeded")
)

func (service *groupServiceImpl) FindByFilter(ctx context.Context, filter business.GroupFilter) ([]business.Group, error) {
	op := runtime.FuncForPC(reflect.ValueOf(service.FindByFilter).Pointer()).Name()

	log := service.log.With(
		slog.String("op", op),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	log.Info("started searching groups by filter")
	timeout := service.cfg.GetEndpointExecutionTimeout()
	foundGroups, err := async.ExecuteWithTimeout(ctx, timeout, func(ctx context.Context) ([]domain.Group, error) {
		return service.repository.FindByFilter(ctx, ConvertToGroupFilter(filter))
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Error("group searching by filter deadline exceeded")
			return make([]business.Group, 0), handling.Wrap(errFindGroupsByFilterDeadlineExceeded, handling.WithCode(codes.DeadlineExceeded))
		}

		log.Error("error while searching groups by filter", logging.Error(err))
		return make([]business.Group, 0), handling.Process(err)
	}

	log.Info("groups successfully found by filter")
	return ConvertToServiceGroups(foundGroups), nil
}
