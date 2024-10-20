package group

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/middleware"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
	"log/slog"
	"reflect"
	"runtime"
)

var (
	errSearchingGroupByID = errors.New("error while searching group by id")
	ErrGroupNotFoundByID  = errors.New("group by id not found in database")
)

func (repository *groupRepositoryImpl) FindByID(ctx context.Context, groupID uuid.UUID) (domain.Group, error) {
	op := runtime.FuncForPC(reflect.ValueOf(repository.FindByID).Pointer()).Name()

	log := repository.log.With(
		slog.String("op", op),
		slog.Any("groupID", groupID),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	log.Debug("started searching group in a database")
	foundGroup := domain.Group{}
	searchResult := repository.db.WithContext(ctx).First(&foundGroup, groupID)
	if searchResult.Error != nil {
		if errors.Is(searchResult.Error, gorm.ErrRecordNotFound) {
			log.Error("group was not found in the database", logging.Error(searchResult.Error))
			return domain.Group{}, handling.New(ErrGroupNotFoundByID.Error(), codes.NotFound)
		}

		log.Error("error while searching group in the database", logging.Error(searchResult.Error))
		return domain.Group{}, handling.New(errSearchingGroupByID.Error(), codes.Internal)
	}

	log.Debug("group was successfully found in a database")
	return foundGroup, nil
}
