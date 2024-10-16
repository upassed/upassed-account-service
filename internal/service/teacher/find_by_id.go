package teacher

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/middleware"
)

var (
	ErrorFindTeacherByIDDeadlineExceeded error = errors.New("find teacher by ud deadline exceeded")
)

func (service *teacherServiceImpl) FindByID(ctx context.Context, teacherID uuid.UUID) (Teacher, error) {
	const op = "teacher.teacherServiceImpl.FindByID()"

	log := service.log.With(
		slog.String("op", op),
		slog.Any("teacherID", teacherID),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	contextWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	resultChannel := make(chan Teacher)
	errorChannel := make(chan error)

	go func() {
		log.Debug("started finding teacher by id")
		foundTeacher, err := service.repository.FindByID(contextWithTimeout, teacherID)
		if err != nil {
			errorChannel <- handling.Process(err)
			return
		}

		log.Debug("teacher successfully found by id")
		resultChannel <- ConvertToServiceTeacher(foundTeacher)
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
