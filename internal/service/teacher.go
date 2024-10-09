package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/handling"
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
	FindByID(context.Context, uuid.UUID) (domain.Teacher, error)
	CheckDuplicateExists(ctx context.Context, reportEmail, username string) (bool, error)
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
	reportEmailExists, err := service.repository.CheckDuplicateExists(ctx, teacher.ReportEmail, teacher.Username)
	if err != nil {
		return business.TeacherCreateResponse{}, handling.HandleApplicationError(err)
	}

	if reportEmailExists {
		log.Error("teacher with this username or report email already exists")
		return business.TeacherCreateResponse{}, handling.NewApplicationError("teacher duplicate found", codes.AlreadyExists)
	}

	domainTeacher := converter.ConvertTeacherToDomain(teacher)
	if err := service.repository.Save(ctx, domainTeacher); err != nil {
		return business.TeacherCreateResponse{}, handling.HandleApplicationError(err)
	}

	log.Debug("teacher successfully created", slog.Any("createdTeacherID", domainTeacher.ID))
	return business.TeacherCreateResponse{
		CreatedTeacherID: domainTeacher.ID,
	}, nil
}

// TODO work with ctx timeout. Add the timeout to ctx passing to repo layer.
func (service *TeacherServiceImpl) FindByID(ctx context.Context, teacherID string) (business.Teacher, error) {
	const op = "TeacherServiceImpl.FindByID()"

	log := service.log.With(
		slog.String("op", op),
		slog.String("teacherID", teacherID),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	log.Debug("started finding teacher by id")
	parsedID, err := uuid.Parse(teacherID)
	if err != nil {
		log.Error("error while parsing teacher id - wrong UUID passed")
		return business.Teacher{}, handling.NewApplicationError(err.Error(), codes.InvalidArgument)
	}

	foundTeacher, err := service.repository.FindByID(ctx, parsedID)
	if err != nil {
		return business.Teacher{}, handling.HandleApplicationError(err)
	}

	log.Debug("teacher successfully found by id")
	return converter.ConvertTeacherToBusiness(foundTeacher), nil
}
