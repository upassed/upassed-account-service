package teacher

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
	errCreateTeacherDeadlineExceeded = errors.New("create teacher deadline exceeded")
)

func (service *teacherServiceImpl) Create(ctx context.Context, teacherToCreate *business.Teacher) (*business.TeacherCreateResponse, error) {
	spanContext, span := otel.Tracer(service.cfg.Tracing.TeacherTracerName).Start(ctx, "teacherService#Create")
	defer span.End()

	log := logging.Wrap(service.log,
		logging.WithOp(service.Create),
		logging.WithCtx(ctx),
		logging.WithAny("teacherUsername", teacherToCreate.Username),
	)

	log.Info("started creating teacher")
	timeout := service.cfg.GetEndpointExecutionTimeout()
	teacherCreateResponse, err := async.ExecuteWithTimeout(spanContext, timeout, func(ctx context.Context) (*business.TeacherCreateResponse, error) {
		duplicateExists, err := service.repository.CheckDuplicateExists(ctx, teacherToCreate.ReportEmail, teacherToCreate.Username)
		if err != nil {
			return nil, handling.Process(err)
		}

		if duplicateExists {
			log.Error("teacher with this username or report email already exists")
			return nil, handling.Wrap(errors.New("teacher duplicate found"), handling.WithCode(codes.AlreadyExists))
		}

		domainTeacher := ConvertToRepositoryTeacher(teacherToCreate)
		if err := service.repository.Save(ctx, domainTeacher); err != nil {
			return nil, handling.Process(err)
		}

		return &business.TeacherCreateResponse{
			CreatedTeacherID: domainTeacher.ID,
		}, nil
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Error("creating teacher deadline exceeded")
			return nil, handling.Wrap(errCreateTeacherDeadlineExceeded, handling.WithCode(codes.DeadlineExceeded))
		}

		log.Error("error while creating a teacher", logging.Error(err))
		return nil, handling.Wrap(err)
	}

	log.Info("teacher successfully created", slog.Any("createdTeacherID", teacherCreateResponse.CreatedTeacherID))
	return teacherCreateResponse, nil
}
