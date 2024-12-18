package group

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logging"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"github.com/upassed/upassed-account-service/internal/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

var (
	errSearchingGroupByID = errors.New("error while searching group by id")
	ErrGroupNotFoundByID  = errors.New("group by id not found in database")
)

func (repository *repositoryImpl) FindByID(ctx context.Context, groupID uuid.UUID) (*domain.Group, error) {
	spanContext, span := otel.Tracer(repository.cfg.Tracing.GroupTracerName).Start(ctx, "groupRepository#FindByID")
	span.SetAttributes(attribute.String("id", groupID.String()))
	defer span.End()

	log := logging.Wrap(repository.log,
		logging.WithOp(repository.FindByID),
		logging.WithCtx(ctx),
		logging.WithAny("groupID", groupID),
	)

	log.Info("started searching group by id in redis cache")
	groupFromCache, err := repository.cache.GetByID(spanContext, groupID)
	if err == nil {
		log.Info("group was found in cache, not going to the database")
		return groupFromCache, nil
	}

	log.Info("started searching group by id in a database")
	foundGroup := domain.Group{}
	searchResult := repository.db.WithContext(ctx).First(&foundGroup, groupID)
	if err := searchResult.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("group by id was not found in the database", logging.Error(err))
			tracing.SetSpanError(span, err)
			return nil, handling.New(ErrGroupNotFoundByID.Error(), codes.NotFound)
		}

		log.Error("error while searching group in the database", logging.Error(err))
		tracing.SetSpanError(span, err)
		return nil, handling.New(errSearchingGroupByID.Error(), codes.Internal)
	}

	log.Info("group by id was successfully found in a database")
	log.Info("saving group to cache")
	if err := repository.cache.Save(spanContext, &foundGroup); err != nil {
		log.Error("error while saving group to cache", logging.Error(err))
		tracing.SetSpanError(span, err)
	}

	log.Info("group was saved to the cache")
	return &foundGroup, nil
}
