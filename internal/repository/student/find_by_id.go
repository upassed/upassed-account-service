package student

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logger"
	"github.com/upassed/upassed-account-service/internal/middleware"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

var (
	ErrorSearchingStudentByID            error = errors.New("error while searching student by id")
	ErrorStudentNotFoundByID             error = errors.New("student by id not found in database")
	ErrorFindStudentByIDDeadlineExceeded error = errors.New("finding student by id in a database deadline exceeded")
)

func (repository *studentRepositoryImpl) FindByID(ctx context.Context, studentID uuid.UUID) (domain.Student, error) {
	const op = "student.studentRepositoryImpl.FindByID()"

	log := repository.log.With(
		slog.String("op", op),
		slog.Any("studentID", studentID),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	contextWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	resultChannel := make(chan domain.Student)
	errorChannel := make(chan error)

	go func() {
		log.Debug("started searching student in a database")
		foundStudent := domain.Student{}
		searchResult := repository.db.Preload("Group").First(&foundStudent, studentID)
		if searchResult.Error != nil {
			if errors.Is(searchResult.Error, gorm.ErrRecordNotFound) {
				log.Error("student was not found in the database", logger.Error(searchResult.Error))
				errorChannel <- handling.New(ErrorStudentNotFoundByID.Error(), codes.NotFound)
				return
			}

			log.Error("error while searching student in the database", logger.Error(searchResult.Error))
			errorChannel <- handling.New(ErrorSearchingStudentByID.Error(), codes.Internal)
			return
		}

		log.Debug("student was successfully found in a database")
		resultChannel <- foundStudent
	}()

	for {
		select {
		case <-contextWithTimeout.Done():
			return domain.Student{}, ErrorFindStudentByIDDeadlineExceeded
		case foundStudent := <-resultChannel:
			return foundStudent, nil
		case err := <-errorChannel:
			return domain.Student{}, err
		}
	}
}
