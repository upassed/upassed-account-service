package group

import (
	"context"
	"log/slog"
	"time"

	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/middleware"
	business "github.com/upassed/upassed-account-service/internal/service/model"
)

func (service *groupServiceImpl) FindByFilter(ctx context.Context, filter business.GroupFilter) ([]business.Group, error) {
	const op = "group.groupServiceImpl.FindByFilter()"

	log := service.log.With(
		slog.String("op", op),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	contextWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	resultChannel := make(chan []business.Group)
	errorChannel := make(chan error)

	go func() {
		log.Debug("started finding groups by filter")
		foundGroups, err := service.repository.FindByFilter(contextWithTimeout, ConvertToGroupFilter(filter))
		if err != nil {
			errorChannel <- handling.Process(err)
			return
		}

		log.Debug("groups successfully found by filter")
		resultChannel <- ConvertToServiceGroups(foundGroups)
	}()

	for {
		select {
		case <-contextWithTimeout.Done():
			return []business.Group{}, errFindGroupByIDDeadlineExceeded
		case foundGroups := <-resultChannel:
			return foundGroups, nil
		case err := <-errorChannel:
			return []business.Group{}, err
		}
	}
}
