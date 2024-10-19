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
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/middleware"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"google.golang.org/grpc/codes"
)

var (
	errSearchingStudentsInGroup            = errors.New("error while searching students in group")
	errFindStudentsInGroupDeadlineExceeded = errors.New("finding students in group in a database deadline exceeded")
)

func (repository *groupRepositoryImpl) FindStudentsInGroup(ctx context.Context, groupID uuid.UUID) ([]domain.Student, error) {
	op := runtime.FuncForPC(reflect.ValueOf(repository.FindStudentsInGroup).Pointer()).Name()

	log := repository.log.With(
		slog.String("op", op),
		slog.Any("groupID", groupID),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	contextWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	resultChannel := make(chan []domain.Student)
	errorChannel := make(chan error)

	go func() {
		log.Debug("started searching students in group in a database")
		foundStudents := make([]domain.Student, 0)
		searchResult := repository.db.Preload("Group").Where("group_id = ?", groupID).Find(&foundStudents)
		if searchResult.Error != nil {
			log.Error("error while searching students in group in the database", logging.Error(searchResult.Error))
			errorChannel <- handling.New(errSearchingStudentsInGroup.Error(), codes.Internal)
			return
		}

		log.Debug("students in group were successfully found in a database", slog.Int("studentInGroup", len(foundStudents)))
		resultChannel <- foundStudents
	}()

	for {
		select {
		case <-contextWithTimeout.Done():
			return []domain.Student{}, errFindStudentsInGroupDeadlineExceeded
		case foundStudentsInGroup := <-resultChannel:
			return foundStudentsInGroup, nil
		case err := <-errorChannel:
			return []domain.Student{}, err
		}
	}
}