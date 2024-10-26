package teacher

import (
	"context"
	"errors"
	"github.com/upassed/upassed-account-service/internal/async"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logging"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"log/slog"
)

var (
	errCreateTeacherDeadlineExceeded = errors.New("create teacher deadline exceeded")
)

func (service *teacherServiceImpl) Create(ctx context.Context, teacherToCreate *business.Teacher) (*business.TeacherCreateResponse, error) {
	spanContext, span := otel.Tracer(service.cfg.Tracing.TeacherTracerName).Start(ctx, "teacherService#Create")
	span.SetAttributes(attribute.String("username", teacherToCreate.Username))
	defer span.End()

	log := logging.Wrap(service.log,
		logging.WithOp(service.Create),
		logging.WithCtx(ctx),
		logging.WithAny("teacherUsername", teacherToCreate.Username),
	)

	log.Info("started creating teacher")
	timeout := service.cfg.GetEndpointExecutionTimeout()
	teacherCreateResponse, err := async.ExecuteWithTimeout(spanContext, timeout, func(ctx context.Context) (*business.TeacherCreateResponse, error) {
		log.Info("checking teacher duplicates")
		duplicateExists, err := service.repository.CheckDuplicateExists(ctx, teacherToCreate.ReportEmail, teacherToCreate.Username)
		if err != nil {
			log.Error("error while checking teacher duplicates", logging.Error(err))
			span.SetAttributes(attribute.String("err", err.Error()))
			return nil, handling.Process(err)
		}

		if duplicateExists {
			log.Error("teacher with this username or report email already exists")
			span.SetAttributes(attribute.String("err", "teacher duplicate found"))
			return nil, handling.Wrap(errors.New("teacher duplicate found"), handling.WithCode(codes.AlreadyExists))
		}

		domainTeacher := ConvertToRepositoryTeacher(teacherToCreate)
		log.Info("saving teacher data to the database")
		if err := service.repository.Save(ctx, domainTeacher); err != nil {
			log.Error("error while saving teacher data to the database", logging.Error(err))
			span.SetAttributes(attribute.String("err", err.Error()))
			return nil, handling.Process(err)
		}

		return &business.TeacherCreateResponse{
			CreatedTeacherID: domainTeacher.ID,
		}, nil
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Error("creating teacher deadline exceeded")
			span.SetAttributes(attribute.String("err", err.Error()))
			return nil, handling.Wrap(errCreateTeacherDeadlineExceeded, handling.WithCode(codes.DeadlineExceeded))
		}

		log.Error("error while creating a teacher", logging.Error(err))
		span.SetAttributes(attribute.String("err", err.Error()))
		return nil, handling.Wrap(err)
	}

	log.Info("teacher successfully created", slog.Any("createdTeacherID", teacherCreateResponse.CreatedTeacherID))
	return teacherCreateResponse, nil
}
