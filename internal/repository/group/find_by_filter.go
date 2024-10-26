package group

import (
	"context"
	"errors"
	"fmt"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logging"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
)

var (
	errSearchingGroupByFilter = errors.New("error while searching groups by filter")
)

const sqlContainsFormat = "%%%s%%"

func (repository *groupRepositoryImpl) FindByFilter(ctx context.Context, filter *domain.GroupFilter) ([]*domain.Group, error) {
	_, span := otel.Tracer(repository.cfg.Tracing.GroupTracerName).Start(ctx, "groupRepository#FindByFilter")
	span.SetAttributes(
		attribute.String("specializationCode", filter.SpecializationCode),
		attribute.String("groupNumber", filter.GroupNumber),
	)
	defer span.End()

	log := logging.Wrap(repository.log,
		logging.WithOp(repository.FindByFilter),
		logging.WithCtx(ctx),
	)

	log.Info("started searching groups by filter in a database")
	foundGroups := make([]*domain.Group, 0)

	specializationCode := fmt.Sprintf(sqlContainsFormat, filter.SpecializationCode)
	groupNumber := fmt.Sprintf(sqlContainsFormat, filter.GroupNumber)
	searchResult := repository.db.WithContext(ctx).Where("specialization_code LIKE ? AND group_number LIKE ?", specializationCode, groupNumber).Find(&foundGroups)

	if err := searchResult.Error; err != nil {
		log.Error("error while searching groups by filter in the database", logging.Error(err))
		span.SetAttributes(attribute.String("err", err.Error()))
		return nil, handling.New(errSearchingGroupByFilter.Error(), codes.Internal)
	}

	log.Info("group was successfully found in a database")
	return foundGroups, nil
}
