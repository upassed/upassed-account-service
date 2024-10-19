package teacher

import (
	"context"
	"errors"
	"log/slog"
	"reflect"
	"runtime"
	"time"

	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/middleware"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"google.golang.org/grpc/codes"
)

var (
	ErrSavingTeacher               = errors.New("error while saving teacher")
	errSaveTeacherDeadlineExceeded = errors.New("saving teacher into a database deadline exceeded")
)

func (repository *teacherRepositoryImpl) Save(ctx context.Context, teacher domain.Teacher) error {
	op := runtime.FuncForPC(reflect.ValueOf(repository.Save).Pointer()).Name()

	log := repository.log.With(
		slog.String("op", op),
		slog.String("teacherUsername", teacher.Username),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	contextWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	resultChannel := make(chan struct{})
	errorChannel := make(chan error)

	go func() {
		log.Debug("started saving teacher to a database")
		saveResult := repository.db.Create(&teacher)
		if saveResult.Error != nil || saveResult.RowsAffected != 1 {
			log.Error("error while saving teacher data to a database", logging.Error(saveResult.Error))
			errorChannel <- handling.New(ErrSavingTeacher.Error(), codes.Internal)
			return
		}

		log.Debug("teacher was successfully inserted into a database")
		resultChannel <- struct{}{}
	}()

	for {
		select {
		case <-contextWithTimeout.Done():
			return errSaveTeacherDeadlineExceeded
		case <-resultChannel:
			return nil
		case err := <-errorChannel:
			return err
		}
	}
}
