package student

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/middleware"
	business "github.com/upassed/upassed-account-service/internal/service/model"
)

var (
	ErrorFindStudentByIDDeadlineExceeded error = errors.New("find student by id deadline exceeded")
)

func (service *studentServiceImpl) FindByID(ctx context.Context, studentID uuid.UUID) (business.Student, error) {
	const op = "student.studentServiceImpl.FindByID()"

	log := service.log.With(
		slog.String("op", op),
		slog.Any("studentID", studentID),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	contextWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	resultChannel := make(chan business.Student)
	errorChannel := make(chan error)

	go func() {
		log.Debug("started finding student by id")
		foundStudent, err := service.studentRepository.FindByID(contextWithTimeout, studentID)
		if err != nil {
			errorChannel <- handling.Process(err)
			return
		}

		log.Debug("student successfully found by id")
		resultChannel <- ConvertToServiceStudent(foundStudent)
	}()

	for {
		select {
		case <-contextWithTimeout.Done():
			return business.Student{}, ErrorFindStudentByIDDeadlineExceeded
		case foundStudent := <-resultChannel:
			return foundStudent, nil
		case err := <-errorChannel:
			return business.Student{}, err
		}
	}
}
