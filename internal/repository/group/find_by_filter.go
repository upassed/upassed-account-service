package group

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logger"
	"github.com/upassed/upassed-account-service/internal/middleware"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"google.golang.org/grpc/codes"
)

var (
	ErrorSearchingGroupByFilter            error = errors.New("error while searching groups by filter")
	ErrorFindGroupByFilterDeadlineExceeded error = errors.New("finding groups by filter in a database deadline exceeded")
)

func (repository *groupRepositoryImpl) FindByFilter(ctx context.Context, filter domain.GroupFilter) ([]domain.Group, error) {
	const op = "group.groupRepositoryImpl.FindByFilter()"

	log := repository.log.With(
		slog.String("op", op),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	contextWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	resultChannel := make(chan []domain.Group)
	errorChannel := make(chan error)

	go func() {
		log.Debug("started searching groups by filter in a database")
		foundGroups := []domain.Group{}

		specializationCode := fmt.Sprintf("%%%s%%", filter.SpecializationCode)
		groupNumber := fmt.Sprintf("%%%s%%", filter.GroupNumber)
		searchResult := repository.db.Where("specialization_code LIKE ? AND group_number LIKE ?", specializationCode, groupNumber).Find(&foundGroups)

		if searchResult.Error != nil {
			log.Error("error while searching groups by filter in the database", logger.Error(searchResult.Error))
			errorChannel <- handling.New(ErrorSearchingGroupByFilter.Error(), codes.Internal)
			return
		}

		log.Debug("group was successfully found in a database")
		resultChannel <- foundGroups
	}()

	for {
		select {
		case <-contextWithTimeout.Done():
			return []domain.Group{}, ErrorFindGroupByFilterDeadlineExceeded
		case foundGroups := <-resultChannel:
			return foundGroups, nil
		case err := <-errorChannel:
			return []domain.Group{}, err
		}
	}
}
