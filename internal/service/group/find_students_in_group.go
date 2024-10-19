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
	errFindStudentsInGroupDeadlineExceeded = errors.New("find students in group deadline exceeded")
)

func (service *groupServiceImpl) FindStudentsInGroup(ctx context.Context, groupID uuid.UUID) ([]business.Student, error) {
	op := runtime.FuncForPC(reflect.ValueOf(service.FindStudentsInGroup).Pointer()).Name()

	log := service.log.With(
		slog.String("op", op),
		slog.Any("groupID", groupID),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	contextWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	resultChannel := make(chan []business.Student)
	errorChannel := make(chan error)

	go func() {
		log.Debug("started finding students in group")
		foundStudentsInGroup, err := service.repository.FindStudentsInGroup(contextWithTimeout, groupID)
		if err != nil {
			errorChannel <- handling.Process(err)
			return
		}

		log.Debug("successfully found students in group", slog.Int("studentsCount", len(foundStudentsInGroup)))
		resultChannel <- ConvertToServiceStudents(foundStudentsInGroup)
	}()

	for {
		select {
		case <-contextWithTimeout.Done():
			return []business.Student{}, errFindStudentsInGroupDeadlineExceeded
		case foundStudentsInGroup := <-resultChannel:
			return foundStudentsInGroup, nil
		case err := <-errorChannel:
			return []business.Student{}, err
		}
	}
}
