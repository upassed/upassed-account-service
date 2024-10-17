package group

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logger"
	"github.com/upassed/upassed-account-service/internal/middleware"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

var (
	ErrorSearchingGroup                error = errors.New("error while searching group")
	ErrorGroupNotFound                 error = errors.New("group not found in database")
	ErrorFindGroupByIDDeadlineExceeded error = errors.New("finding group by id in a database deadline exceeded")
)

func (repository *groupRepositoryImpl) FindByID(ctx context.Context, groupID uuid.UUID) (domain.Group, error) {
	const op = "group.groupRepositoryImpl.FindByID()"

	log := repository.log.With(
		slog.String("op", op),
		slog.Any("groupID", groupID),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	contextWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	resultChannel := make(chan domain.Group)
	errorChannel := make(chan error)

	go func() {
		log.Debug("started searching group in a database")
		foundGroup := domain.Group{}
		searchResult := repository.db.First(&foundGroup, groupID)
		if searchResult.Error != nil {
			if errors.Is(searchResult.Error, gorm.ErrRecordNotFound) {
				log.Error("group was not found in the database", logger.Error(searchResult.Error))
				errorChannel <- handling.New(ErrorGroupNotFound.Error(), codes.NotFound)
				return
			}

			log.Error("error while searching group in the database", logger.Error(searchResult.Error))
			errorChannel <- handling.New(ErrorSearchingGroup.Error(), codes.Internal)
			return
		}

		log.Debug("group was successfully found in a database")
		resultChannel <- foundGroup
	}()

	for {
		select {
		case <-contextWithTimeout.Done():
			return domain.Group{}, ErrorFindGroupByIDDeadlineExceeded
		case foundGroup := <-resultChannel:
			return foundGroup, nil
		case err := <-errorChannel:
			return domain.Group{}, err
		}
	}
}
