package group

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/async"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logging"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
)

var (
	errFindGroupByIDDeadlineExceeded = errors.New("find group by id timeout exceeded")
)

func (service *groupServiceImpl) FindByID(ctx context.Context, groupID uuid.UUID) (*business.Group, error) {
	log := logging.Wrap(service.log,
		logging.WithOp(service.FindByID),
		logging.WithCtx(ctx),
		logging.WithAny("groupID", groupID),
	)

	spanContext, span := otel.Tracer(service.cfg.Tracing.GroupTracerName).Start(ctx, "groupService#FindByID")
	defer span.End()

	log.Info("started searching group by id")
	timeout := service.cfg.GetEndpointExecutionTimeout()
	foundGroup, err := async.ExecuteWithTimeout(spanContext, timeout, func(ctx context.Context) (*domain.Group, error) {
		return service.repository.FindByID(ctx, groupID)
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Error("group searching by id deadline exceeded")
			return nil, handling.Wrap(errFindGroupByIDDeadlineExceeded, handling.WithCode(codes.DeadlineExceeded))
		}

		log.Error("error while searching group by id", logging.Error(err))
		return nil, handling.Process(err)
	}

	log.Info("group successfully found by id")
	return ConvertToServiceGroup(foundGroup), nil
}
