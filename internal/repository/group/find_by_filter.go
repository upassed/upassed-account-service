package group

import (
	"context"
	"errors"
	"fmt"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/middleware"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"google.golang.org/grpc/codes"
	"log/slog"
	"reflect"
	"runtime"
)

var (
	errSearchingGroupByFilter = errors.New("error while searching groups by filter")
)

func (repository *groupRepositoryImpl) FindByFilter(ctx context.Context, filter domain.GroupFilter) ([]domain.Group, error) {
	op := runtime.FuncForPC(reflect.ValueOf(repository.FindByFilter).Pointer()).Name()

	log := repository.log.With(
		slog.String("op", op),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	log.Info("started searching groups by filter in a database")
	foundGroups := make([]domain.Group, 0)

	specializationCode := fmt.Sprintf("%%%s%%", filter.SpecializationCode)
	groupNumber := fmt.Sprintf("%%%s%%", filter.GroupNumber)
	searchResult := repository.db.WithContext(ctx).Where("specialization_code LIKE ? AND group_number LIKE ?", specializationCode, groupNumber).Find(&foundGroups)

	if searchResult.Error != nil {
		log.Error("error while searching groups by filter in the database", logging.Error(searchResult.Error))
		return make([]domain.Group, 0), handling.New(errSearchingGroupByFilter.Error(), codes.Internal)
	}

	log.Info("group was successfully found in a database")
	return foundGroups, nil
}
