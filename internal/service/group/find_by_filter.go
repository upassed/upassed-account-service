package group

import (
	"context"
	"errors"
	"github.com/upassed/upassed-account-service/internal/async"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logging"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
)

var (
	errFindGroupsByFilterDeadlineExceeded = errors.New("find groups by filter timeout exceeded")
)

func (service *groupServiceImpl) FindByFilter(ctx context.Context, filter *business.GroupFilter) ([]*business.Group, error) {
	spanContext, span := otel.Tracer(service.cfg.Tracing.GroupTracerName).Start(ctx, "groupService#FindByFilter")
	defer span.End()

	log := logging.Wrap(service.log,
		logging.WithOp(service.FindByFilter),
		logging.WithCtx(ctx),
	)

	log.Info("started searching groups by filter")
	timeout := service.cfg.GetEndpointExecutionTimeout()
	foundGroups, err := async.ExecuteWithTimeout(spanContext, timeout, func(ctx context.Context) ([]*domain.Group, error) {
		return service.repository.FindByFilter(ctx, ConvertToGroupFilter(filter))
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Error("group searching by filter deadline exceeded")
			return nil, handling.Wrap(errFindGroupsByFilterDeadlineExceeded, handling.WithCode(codes.DeadlineExceeded))
		}

		log.Error("error while searching groups by filter", logging.Error(err))
		return nil, handling.Process(err)
	}

	log.Info("groups successfully found by filter")
	return ConvertToServiceGroups(foundGroups), nil
}
