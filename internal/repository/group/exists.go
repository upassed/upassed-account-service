package group

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/middleware"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"google.golang.org/grpc/codes"
	"log/slog"
	"reflect"
	"runtime"
)

var (
	errCheckGroupExists = errors.New("error while checking if group exists in database")
)

func (repository *groupRepositoryImpl) Exists(ctx context.Context, groupID uuid.UUID) (bool, error) {
	op := runtime.FuncForPC(reflect.ValueOf(repository.Exists).Pointer()).Name()

	log := repository.log.With(
		slog.String("op", op),
		slog.Any("groupID", groupID),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
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
