package student

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
	ErrorSavingStudent               error = errors.New("error while saving student")
	ErrorSaveStudentDeadlineExceeded error = errors.New("saving student into a database deadline exceeded")
)

func (repository *studentRepositoryImpl) Save(ctx context.Context, student Student) error {
	const op = "student.studentRepositoryImpl.Save()"

	log := repository.log.With(
		slog.String("op", op),
		slog.String("studentUsername", student.Username),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	contextWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	resultChannel := make(chan struct{})
	errorChannel := make(chan error)

	go func() {
		log.Debug("started saving student to a database")
		saveResult := repository.db.Create(&student)
		if saveResult.Error != nil || saveResult.RowsAffected != 1 {
			log.Error("error while saving student data to a database", logger.Error(saveResult.Error))
			errorChannel <- handling.New(ErrorSavingStudent.Error(), codes.Internal)
			return
		}

		log.Debug("student was successfully inserted into a database")
		resultChannel <- struct{}{}
	}()

	for {
		select {
		case <-contextWithTimeout.Done():
			return ErrorSaveStudentDeadlineExceeded
		case <-resultChannel:
			return nil
		case err := <-errorChannel:
			return err
		}
	}
}
