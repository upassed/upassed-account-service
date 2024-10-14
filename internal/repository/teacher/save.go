package teacher

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logger"
	"github.com/upassed/upassed-account-service/internal/middleware"
	"google.golang.org/grpc/codes"
)

var (
	ErrorSavingTeacher               error = errors.New("error while saving teacher")
	ErrorSaveTeacherDeadlineExceeded error = errors.New("saving teacher into a database deadline exceeded")
)

func (repository *teacherRepositoryImpl) Save(ctx context.Context, teacher Teacher) error {
	const op = "repository.TeacherRepositoryImpl.Save()"

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
			log.Error("error while saving teacher data to a database", logger.Error(saveResult.Error))
			errorChannel <- handling.NewApplicationError(ErrorSavingTeacher.Error(), codes.Internal)
			return
		}

		log.Debug("teacher was successfully inserted into a database")
		resultChannel <- struct{}{}
	}()

	for {
		select {
		case <-contextWithTimeout.Done():
			return ErrorSaveTeacherDeadlineExceeded
		case <-resultChannel:
			return nil
		case err := <-errorChannel:
			return err
		}
	}
}
