package teacher

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
	errFindTeacherByIDDeadlineExceeded = errors.New("find teacher by id deadline exceeded")
)

func (service *teacherServiceImpl) FindByID(ctx context.Context, teacherID uuid.UUID) (business.Teacher, error) {
	const op = "teacher.teacherServiceImpl.FindByID()"

	log := service.log.With(
		slog.String("op", op),
		slog.Any("teacherID", teacherID),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	contextWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	resultChannel := make(chan business.Teacher)
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
			return business.Teacher{}, errFindTeacherByIDDeadlineExceeded
		case foundTeacher := <-resultChannel:
			return foundTeacher, nil
		case err := <-errorChannel:
			return business.Teacher{}, err
		}
	}
}
