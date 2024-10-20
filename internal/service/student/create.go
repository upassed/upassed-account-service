package student

import (
	"context"
	"errors"
	"github.com/upassed/upassed-account-service/internal/async"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/middleware"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"google.golang.org/grpc/codes"
	"log/slog"
	"reflect"
	"runtime"
)

var (
	errCreateStudentDeadlineExceeded = errors.New("create student deadline exceeded")
)

func (service *studentServiceImpl) Create(ctx context.Context, student business.Student) (business.StudentCreateResponse, error) {
	op := runtime.FuncForPC(reflect.ValueOf(service.Create).Pointer()).Name()

	log := service.log.With(
		slog.String("op", op),
		slog.String("studentUsername", student.Username),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	timeout := service.cfg.GetEndpointExecutionTimeout()
	studentCreateResponse, err := async.ExecuteWithTimeout(ctx, timeout, func(ctx context.Context) (business.StudentCreateResponse, error) {
		log.Debug("started creating student")
		duplicateExists, err := service.studentRepository.CheckDuplicateExists(ctx, student.EducationalEmail, student.Username)
		if err != nil {
			return business.StudentCreateResponse{}, err
		}

		if duplicateExists {
			log.Error("student with this username or educational email already exists")
			return business.StudentCreateResponse{}, handling.Wrap(errors.New("student duplicate found"), handling.WithCode(codes.AlreadyExists))
		}

		groupExists, err := service.groupRepository.Exists(ctx, student.Group.ID)
		if err != nil {
			return business.StudentCreateResponse{}, err
		}

		if !groupExists {
			log.Error("group with this id was not found in database", slog.Any("groupID", student.Group.ID))
			return business.StudentCreateResponse{}, handling.Wrap(errors.New("group does not exists by id"), handling.WithCode(codes.NotFound))
		}

		domainStudent := ConvertToRepositoryStudent(student)
		if err := service.studentRepository.Save(ctx, domainStudent); err != nil {
			return business.StudentCreateResponse{}, err
		}

		log.Debug("student successfully created", slog.Any("createdStudentID", domainStudent.ID))
		return business.StudentCreateResponse{
			CreatedStudentID: domainStudent.ID,
		}, nil
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return business.StudentCreateResponse{}, handling.Wrap(errCreateStudentDeadlineExceeded, handling.WithCode(codes.DeadlineExceeded))
		}

		return business.StudentCreateResponse{}, handling.Process(err)
	}

	return studentCreateResponse, nil
}
