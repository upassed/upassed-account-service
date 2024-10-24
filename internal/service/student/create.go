package student

import (
	"context"
	"errors"
	"github.com/upassed/upassed-account-service/internal/async"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/middleware"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"go.opentelemetry.io/otel"
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

	spanContext, span := otel.Tracer(service.cfg.Tracing.StudentTracerName).Start(ctx, "studentService#Create")
	defer span.End()

	log.Info("started creating student")
	timeout := service.cfg.GetEndpointExecutionTimeout()
	studentCreateResponse, err := async.ExecuteWithTimeout(spanContext, timeout, func(ctx context.Context) (business.StudentCreateResponse, error) {
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
		existingGroup, err := service.groupRepository.FindByID(ctx, student.Group.ID)
		if err != nil {
			log.Error("error while searching group by id", logging.Error(err))
			return business.StudentCreateResponse{}, handling.Wrap(errors.New("error searching group"), handling.WithCode(codes.Internal))
		}

		domainStudent.Group = existingGroup
		if err := service.studentRepository.Save(ctx, domainStudent); err != nil {
			return business.StudentCreateResponse{}, err
		}

		return business.StudentCreateResponse{
			CreatedStudentID: domainStudent.ID,
		}, nil
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Error("student creating deadline exceeded")
			return business.StudentCreateResponse{}, handling.Wrap(errCreateStudentDeadlineExceeded, handling.WithCode(codes.DeadlineExceeded))
		}

		log.Error("error while creating student", logging.Error(err))
		return business.StudentCreateResponse{}, handling.Process(err)
	}

	log.Info("student successfully created", slog.Any("createdStudentID", studentCreateResponse.CreatedStudentID))
	return studentCreateResponse, nil
}
