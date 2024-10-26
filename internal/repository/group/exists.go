package group

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logging"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
)

var (
	errCheckGroupExists = errors.New("error while checking if group exists in database")
)

func (repository *groupRepositoryImpl) Exists(ctx context.Context, groupID uuid.UUID) (bool, error) {
	_, span := otel.Tracer(repository.cfg.Tracing.GroupTracerName).Start(ctx, "groupRepository#Exists")
	defer span.End()

	log := logging.Wrap(repository.log,
		logging.WithOp(repository.Exists),
		logging.WithCtx(ctx),
		logging.WithAny("groupID", groupID),
	)

	log.Info("started checking group exists")
	var groupCount int64
	countResult := repository.db.WithContext(ctx).Model(&domain.Group{}).Where("id = ?", groupID).Count(&groupCount)
	if countResult.Error != nil {
		log.Error("error while counting groups with id in database")
		return false, handling.New(errCheckGroupExists.Error(), codes.Internal)
	}

	if groupCount > 0 {
		log.Info("group exists in database")
		return true, nil
	}

	log.Info("group does not exists in database")
	return false, nil
}
