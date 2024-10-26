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
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"log/slog"
)

var (
	errFindStudentsInGroupDeadlineExceeded = errors.New("find students in group deadline exceeded")
)

func (service *groupServiceImpl) FindStudentsInGroup(ctx context.Context, groupID uuid.UUID) ([]*business.Student, error) {
	spanContext, span := otel.Tracer(service.cfg.Tracing.GroupTracerName).Start(ctx, "groupService#FindStudentsInGroup")
	span.SetAttributes(attribute.String("id", groupID.String()))
	defer span.End()

	log := logging.Wrap(service.log,
		logging.WithOp(service.FindByID),
		logging.WithCtx(ctx),
		logging.WithAny("groupID", groupID),
	)

	log.Info("started searching students in group")
	timeout := service.cfg.GetEndpointExecutionTimeout()
	foundStudents, err := async.ExecuteWithTimeout(spanContext, timeout, func(ctx context.Context) ([]*domain.Student, error) {
		return service.repository.FindStudentsInGroup(ctx, groupID)
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Error("students in group searching deadline exceeded")
			span.SetAttributes(attribute.String("err", err.Error()))
			return nil, handling.Wrap(errFindStudentsInGroupDeadlineExceeded, handling.WithCode(codes.DeadlineExceeded))
		}

		log.Error("error while searching students in group", logging.Error(err))
		span.SetAttributes(attribute.String("err", err.Error()))
		return nil, handling.Process(err)
	}

	log.Info("successfully found students in group", slog.Int("studentsCount", len(foundStudents)))
	return ConvertToServiceStudents(foundStudents), nil
}
