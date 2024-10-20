package teacher

import (
	"context"
	"errors"
	"github.com/upassed/upassed-account-service/internal/async"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/middleware"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"google.golang.org/grpc/codes"
	"log/slog"
	"reflect"
	"runtime"
)

var (
	errCreateTeacherDeadlineExceeded = errors.New("create teacher deadline exceeded")
)

func (service *teacherServiceImpl) Create(ctx context.Context, teacherToCreate business.Teacher) (business.TeacherCreateResponse, error) {
	op := runtime.FuncForPC(reflect.ValueOf(service.Create).Pointer()).Name()

	log := service.log.With(
		slog.String("op", op),
		slog.String("teacherUsername", teacherToCreate.Username),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	log.Info("started creating teacher")
	timeout := service.cfg.GetEndpointExecutionTimeout()
	teacherCreateResponse, err := async.ExecuteWithTimeout(ctx, timeout, func(ctx context.Context) (business.TeacherCreateResponse, error) {
		duplicateExists, err := service.repository.CheckDuplicateExists(ctx, teacherToCreate.ReportEmail, teacherToCreate.Username)
		if err != nil {
			return business.TeacherCreateResponse{}, handling.Process(err)
		}

		if duplicateExists {
			log.Error("teacher with this username or report email already exists")
			return business.TeacherCreateResponse{}, handling.Wrap(errors.New("teacher duplicate found"), handling.WithCode(codes.AlreadyExists))
		}

		domainTeacher := ConvertToRepositoryTeacher(teacherToCreate)
		if err := service.repository.Save(ctx, domainTeacher); err != nil {
			return business.TeacherCreateResponse{}, handling.Process(err)
		}

		return business.TeacherCreateResponse{
			CreatedTeacherID: domainTeacher.ID,
		}, nil
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Error("creating teacher deadline exceeded")
			return business.TeacherCreateResponse{}, handling.Wrap(errCreateTeacherDeadlineExceeded, handling.WithCode(codes.DeadlineExceeded))
		}

		log.Error("error while creating a teacher", logging.Error(err))
		return business.TeacherCreateResponse{}, handling.Wrap(err)
	}

	log.Info("teacher successfully created", slog.Any("createdTeacherID", teacherCreateResponse.CreatedTeacherID))
	return teacherCreateResponse, nil
}
