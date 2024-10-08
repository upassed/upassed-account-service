package service

import (
	"context"
	"log/slog"

	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logger"
	"github.com/upassed/upassed-account-service/internal/middleware"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"github.com/upassed/upassed-account-service/internal/service/converter"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"google.golang.org/grpc/codes"
)

type TeacherServiceImpl struct {
	log        *slog.Logger
	repository TeacherRepository
}

type TeacherRepository interface {
	Save(context.Context, domain.Teacher) error
}

func NewTeacherService(log *slog.Logger, repository TeacherRepository) *TeacherServiceImpl {
	return &TeacherServiceImpl{
		log:        log,
		repository: repository,
	}
}

// TODO work with ctx timeout. Add the timeout to ctx passing to repo layer.
func (service *TeacherServiceImpl) Create(ctx context.Context, teacher business.Teacher) (business.TeacherCreateResponse, error) {
	const op = "TeacherServiceImpl.Create()"

	log := service.log.With(
		slog.String("op", op),
		slog.String("teacherUsername", teacher.Username),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	log.Debug("started creating teacher")
	domainTeacher := converter.ConvertTeacher(teacher)
	if err := service.repository.Save(ctx, domainTeacher); err != nil {
		log.Error("error while saving a teacher to database", logger.Error(err))
		return business.TeacherCreateResponse{}, handling.NewServiceLayerError(err.Error(), codes.Internal)
	}

	log.Debug("teacher successfully created")
	return business.TeacherCreateResponse{
		CreatedTeacherID: domainTeacher.ID,
	}, nil
}
