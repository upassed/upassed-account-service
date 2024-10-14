package teacher

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logger"
	"github.com/upassed/upassed-account-service/internal/middleware"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

var (
	ErrorSearchingTeacher                error = errors.New("error while searching teacher")
	ErrorTeacherNotFound                 error = errors.New("teacher not found in database")
	ErrorFindTeacherByIDDeadlineExceeded error = errors.New("finding teacher by id in a database deadline exceeded")
)

func (repository *teacherRepositoryImpl) FindByID(ctx context.Context, teacherID uuid.UUID) (Teacher, error) {
	const op = "repository.TeacherRepositoryImpl.FindByID()"

	log := repository.log.With(
		slog.String("op", op),
		slog.String("teacherID", teacherID.String()),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	contextWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	resultChannel := make(chan Teacher)
	errorChannel := make(chan error)

	go func() {
		log.Debug("started searching teacher in a database")
		foundTeacher := Teacher{}
		searchResult := repository.db.First(&foundTeacher, teacherID)
		if searchResult.Error != nil {
			if errors.Is(searchResult.Error, gorm.ErrRecordNotFound) {
				log.Error("teacher was not found in the database", logger.Error(searchResult.Error))
				errorChannel <- handling.NewApplicationError(ErrorTeacherNotFound.Error(), codes.NotFound)
				return
			}

			log.Error("error while searching teacher in the database", logger.Error(searchResult.Error))
			errorChannel <- handling.NewApplicationError(ErrorSearchingTeacher.Error(), codes.Internal)
			return
		}

		log.Debug("teacher was successfully found in a database")
		resultChannel <- foundTeacher
	}()

	for {
		select {
		case <-contextWithTimeout.Done():
			return Teacher{}, ErrorFindTeacherByIDDeadlineExceeded
		case foundTeacher := <-resultChannel:
			return foundTeacher, nil
		case err := <-errorChannel:
			return Teacher{}, err
		}
	}
}
