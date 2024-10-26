package student

import (
	"context"
	"errors"
	"github.com/upassed/upassed-account-service/internal/async"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logging"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"log/slog"
)

var (
	errCreateStudentDeadlineExceeded = errors.New("create student deadline exceeded")
)

func (service *studentServiceImpl) Create(ctx context.Context, student *business.Student) (*business.StudentCreateResponse, error) {
	spanContext, span := otel.Tracer(service.cfg.Tracing.StudentTracerName).Start(ctx, "studentService#Create")
	defer span.End()

	log := logging.Wrap(service.log,
		logging.WithOp(service.Create),
		logging.WithCtx(ctx),
		logging.WithAny("studentUsername", student.Username),
	)

	log.Info("started creating student")
	timeout := service.cfg.GetEndpointExecutionTimeout()
	studentCreateResponse, err := async.ExecuteWithTimeout(spanContext, timeout, func(ctx context.Context) (*business.StudentCreateResponse, error) {
		duplicateExists, err := service.studentRepository.CheckDuplicateExists(ctx, student.EducationalEmail, student.Username)
		if err != nil {
			return nil, err
		}

		if duplicateExists {
			log.Error("student with this username or educational email already exists")
			return nil, handling.Wrap(errors.New("student duplicate found"), handling.WithCode(codes.AlreadyExists))
		}

		groupExists, err := service.groupRepository.Exists(ctx, student.Group.ID)
		if err != nil {
			return nil, err
		}

		if !groupExists {
			log.Error("group with this id was not found in database", slog.Any("groupID", student.Group.ID))
			return nil, handling.Wrap(errors.New("group does not exists by id"), handling.WithCode(codes.NotFound))
		}

		domainStudent := ConvertToRepositoryStudent(student)
		existingGroup, err := service.groupRepository.FindByID(ctx, student.Group.ID)
		if err != nil {
			log.Error("error while searching group by id", logging.Error(err))
			return nil, handling.Wrap(errors.New("error searching group"), handling.WithCode(codes.Internal))
		}

		domainStudent.Group = *existingGroup
		if err := service.studentRepository.Save(ctx, domainStudent); err != nil {
			return nil, err
		}

		return &business.StudentCreateResponse{
			CreatedStudentID: domainStudent.ID,
		}, nil
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Error("student creating deadline exceeded")
			return nil, handling.Wrap(errCreateStudentDeadlineExceeded, handling.WithCode(codes.DeadlineExceeded))
		}

		log.Error("error while creating student", logging.Error(err))
		return nil, handling.Process(err)
	}

	log.Info("student successfully created", slog.Any("createdStudentID", studentCreateResponse.CreatedStudentID))
	return studentCreateResponse, nil
}
