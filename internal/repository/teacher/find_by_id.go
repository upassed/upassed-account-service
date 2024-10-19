package teacher

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
	"gorm.io/gorm"
)

var (
	errSearchingTeacherByID            = errors.New("error while searching teacher by id")
	ErrTeacherNotFoundByID             = errors.New("teacher by id  not found in database")
	errFindTeacherByIDDeadlineExceeded = errors.New("finding teacher by id in a database deadline exceeded")
)

func (repository *teacherRepositoryImpl) FindByID(ctx context.Context, teacherID uuid.UUID) (domain.Teacher, error) {
	op := runtime.FuncForPC(reflect.ValueOf(repository.FindByID).Pointer()).Name()

	log := repository.log.With(
		slog.String("op", op),
		slog.Any("teacherID", teacherID),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	contextWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	resultChannel := make(chan domain.Teacher)
	errorChannel := make(chan error)

	go func() {
		log.Debug("started searching teacher in a database")
		foundTeacher := domain.Teacher{}
		searchResult := repository.db.First(&foundTeacher, teacherID)
		if searchResult.Error != nil {
			if errors.Is(searchResult.Error, gorm.ErrRecordNotFound) {
				log.Error("teacher was not found in the database", logging.Error(searchResult.Error))
				errorChannel <- handling.New(ErrTeacherNotFoundByID.Error(), codes.NotFound)
				return
			}

			log.Error("error while searching teacher in the database", logging.Error(searchResult.Error))
			errorChannel <- handling.New(errSearchingTeacherByID.Error(), codes.Internal)
			return
		}

		log.Debug("teacher was successfully found in a database")
		resultChannel <- foundTeacher
	}()

	for {
		select {
		case <-contextWithTimeout.Done():
			return domain.Teacher{}, errFindTeacherByIDDeadlineExceeded
		case foundTeacher := <-resultChannel:
			return foundTeacher, nil
		case err := <-errorChannel:
			return domain.Teacher{}, err
		}
	}
}
