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
	"log/slog"
	"reflect"
	"runtime"
)

var (
	errSearchingStudentsInGroup = errors.New("error while searching students in group")
)

func (repository *groupRepositoryImpl) FindStudentsInGroup(ctx context.Context, groupID uuid.UUID) ([]domain.Student, error) {
	op := runtime.FuncForPC(reflect.ValueOf(repository.FindStudentsInGroup).Pointer()).Name()

	log := repository.log.With(
		slog.String("op", op),
		slog.Any("groupID", groupID),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	log.Info("started searching students in group in a database")
	foundStudents := make([]domain.Student, 0)
	searchResult := repository.db.WithContext(ctx).Preload("Group").Where("group_id = ?", groupID).Find(&foundStudents)
	if searchResult.Error != nil {
		log.Error("error while searching students in group in the database", logging.Error(searchResult.Error))
		return make([]domain.Student, 0), handling.New(errSearchingStudentsInGroup.Error(), codes.Internal)
	}

	log.Info("students in group were successfully found in a database", slog.Int("studentInGroup", len(foundStudents)))
	return foundStudents, nil
}
