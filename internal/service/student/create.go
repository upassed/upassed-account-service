package student

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
	errCreateStudentDeadlineExceeded = errors.New("create student deadline exceeded")
)

func (service *studentServiceImpl) Create(ctx context.Context, student business.Student) (business.StudentCreateResponse, error) {
	const op = "student.studentServiceImpl.Create()"

	log := service.log.With(
		slog.String("op", op),
		slog.String("studentUsername", student.Username),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	contextWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	resultChannel := make(chan business.StudentCreateResponse)
	errorChannel := make(chan error)

	go func() {
		log.Debug("started creating student")
		duplicateExists, err := service.studentRepository.CheckDuplicateExists(contextWithTimeout, student.EducationalEmail, student.Username)
		if err != nil {
			errorChannel <- handling.Process(err)
			return
		}

		if duplicateExists {
			log.Error("student with this username or educational email already exists")
			errorChannel <- handling.Wrap(errors.New("student duplicate found"), handling.WithCode(codes.AlreadyExists))
			return
		}

		groupExists, err := service.groupRepository.Exists(contextWithTimeout, student.Group.ID)
		if err != nil {
			errorChannel <- handling.Process(err)
			return
		}

		if !groupExists {
			log.Error("group with this id was not found in database", slog.Any("groupID", student.Group.ID))
			errorChannel <- handling.Wrap(errors.New("group does not exists by id"), handling.WithCode(codes.NotFound))
			return
		}

		domainStudent := ConvertToRepositoryStudent(student)
		if err := service.studentRepository.Save(contextWithTimeout, domainStudent); err != nil {
			errorChannel <- handling.Process(err)
			return
		}

		log.Debug("student successfully created", slog.Any("createdStudentID", domainStudent.ID))
		resultChannel <- business.StudentCreateResponse{
			CreatedStudentID: domainStudent.ID,
		}
	}()

	for {
		select {
		case <-contextWithTimeout.Done():
			return business.StudentCreateResponse{}, errCreateStudentDeadlineExceeded
		case createdStudentData := <-resultChannel:
			return createdStudentData, nil
		case err := <-errorChannel:
			return business.StudentCreateResponse{}, err
		}
	}
}
