package teacher

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/middleware"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"google.golang.org/grpc/codes"
)

var (
	errCreateTeacherDeadlineExceeded = errors.New("create teacher deadline exceeded")
)

func (service *teacherServiceImpl) Create(ctx context.Context, teacherToCreate business.Teacher) (business.TeacherCreateResponse, error) {
	const op = "teacher.teacherServiceImpl.Create()"

	log := service.log.With(
		slog.String("op", op),
		slog.String("teacherUsername", teacherToCreate.Username),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	contextWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	resultChannel := make(chan business.TeacherCreateResponse)
	errorChannel := make(chan error)

	go func() {
		log.Debug("started creating teacher")
		duplicateExists, err := service.repository.CheckDuplicateExists(contextWithTimeout, teacherToCreate.ReportEmail, teacherToCreate.Username)
		if err != nil {
			errorChannel <- handling.Process(err)
			return
		}

		if duplicateExists {
			log.Error("teacher with this username or report email already exists")
			errorChannel <- handling.Wrap(errors.New("teacher duplicate found"), handling.WithCode(codes.AlreadyExists))
			return
		}

		domainTeacher := ConvertToRepositoryTeacher(teacherToCreate)
		if err := service.repository.Save(contextWithTimeout, domainTeacher); err != nil {
			errorChannel <- handling.Process(err)
			return
		}

		log.Debug("teacher successfully created", slog.Any("createdTeacherID", domainTeacher.ID))
		resultChannel <- business.TeacherCreateResponse{
			CreatedTeacherID: domainTeacher.ID,
		}
	}()

	for {
		select {
		case <-contextWithTimeout.Done():
			return business.TeacherCreateResponse{}, errCreateTeacherDeadlineExceeded
		case createdTeacherData := <-resultChannel:
			return createdTeacherData, nil
		case err := <-errorChannel:
			return business.TeacherCreateResponse{}, err
		}
	}
}
