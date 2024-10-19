package group

import (
	"context"
	"errors"
	"log/slog"
	"reflect"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/middleware"
	business "github.com/upassed/upassed-account-service/internal/service/model"
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

	contextWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	resultChannel := make(chan business.Group)
	errorChannel := make(chan error)

	go func() {
		log.Debug("started finding group by id")
		foundGroup, err := service.repository.FindByID(contextWithTimeout, groupID)
		if err != nil {
			errorChannel <- handling.Process(err)
			return
		}

		log.Debug("group successfully found by id")
		resultChannel <- ConvertToServiceGroup(foundGroup)
	}()

	for {
		select {
		case <-contextWithTimeout.Done():
			return business.Group{}, errFindGroupByIDDeadlineExceeded
		case foundGroup := <-resultChannel:
			return foundGroup, nil
		case err := <-errorChannel:
			return business.Group{}, err
		}
	}
}
