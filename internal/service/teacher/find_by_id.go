package teacher

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/middleware"
	"google.golang.org/grpc/codes"
)

var (
	ErrorFindTeacherByIDDeadlineExceeded error = errors.New("find teacher by ud deadline exceeded")
)

func (service *teacherServiceImpl) FindByID(ctx context.Context, teacherID string) (Teacher, error) {
	const op = "TeacherServiceImpl.FindByID()"

	log := service.log.With(
		slog.String("op", op),
		slog.String("teacherID", teacherID),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	contextWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	resultChannel := make(chan Teacher)
	errorChannel := make(chan error)

	go func() {
		log.Debug("started finding teacher by id")
		parsedID, err := uuid.Parse(teacherID)
		if err != nil {
			log.Error("error while parsing teacher id - wrong UUID passed")
			errorChannel <- handling.Wrap(err, handling.WithCode(codes.InvalidArgument))
			return
		}

		foundTeacher, err := service.repository.FindByID(contextWithTimeout, parsedID)
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
