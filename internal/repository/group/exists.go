package group

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/middleware"
	"google.golang.org/grpc/codes"
)

var (
	ErrorCheckGroupExists                error = errors.New("error while checking if group exists in database")
	ErrorCheckGroupExistsTimeoutExceeded error = errors.New("echecking if group exists in database timeout exceeded")
)

func (repository *groupRepositoryImpl) Exists(ctx context.Context, groupID uuid.UUID) (bool, error) {
	const op = "group.groupRepositoryImpl.Exists()"

	log := repository.log.With(
		slog.String("op", op),
		slog.Any("groupID", groupID),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	contextWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	resultChannel := make(chan bool)
	errorChannel := make(chan error)

	go func() {
		log.Debug("started checking group exists")
		var groupCount int64
		countResult := repository.db.Model(&Group{}).Where("id = ?", groupID).Count(&groupCount)
		if countResult.Error != nil {
			log.Error("error while counting groups with id in database")
			errorChannel <- handling.New(ErrorCheckGroupExists.Error(), codes.Internal)
			return
		}

		if groupCount > 0 {
			log.Debug("group exists in database")
			resultChannel <- true
			return
		}

		log.Debug("group does not exists in database")
		resultChannel <- false
	}()

	for {
		select {
		case <-contextWithTimeout.Done():
			return false, ErrorCheckGroupExistsTimeoutExceeded
		case duplicatesFound := <-resultChannel:
			return duplicatesFound, nil
		case err := <-errorChannel:
			return false, err
		}
	}
}
