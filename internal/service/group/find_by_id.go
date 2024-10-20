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
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"log/slog"
	"reflect"
	"runtime"
)

var (
	errFindGroupByIDDeadlineExceeded = errors.New("find group by id timeout exceeded")
)

func (service *groupServiceImpl) FindByID(ctx context.Context, groupID uuid.UUID) (business.Group, error) {
	op := runtime.FuncForPC(reflect.ValueOf(service.FindByID).Pointer()).Name()

	log := service.log.With(
		slog.String("op", op),
		slog.Any("groupID", groupID),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	spanContext, span := otel.Tracer(service.cfg.Tracing.GroupTracerName).Start(ctx, "groupService#FindByID")
	defer span.End()

	log.Info("started searching group by id")
	timeout := service.cfg.GetEndpointExecutionTimeout()
	foundGroup, err := async.ExecuteWithTimeout(spanContext, timeout, func(ctx context.Context) (domain.Group, error) {
		return service.repository.FindByID(ctx, groupID)
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Error("group searching by id deadline exceeded")
			return business.Group{}, handling.Wrap(errFindGroupByIDDeadlineExceeded, handling.WithCode(codes.DeadlineExceeded))
		}

		log.Error("error while searching group by id", logging.Error(err))
		return business.Group{}, handling.Process(err)
	}

	log.Info("group successfully found by id")
	return ConvertToServiceGroup(foundGroup), nil
}
