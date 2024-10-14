package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/middleware"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"github.com/upassed/upassed-account-service/internal/service/converter"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"google.golang.org/grpc/codes"
)

var (
	ErrorCreateTeacherDeadlineExceeded   error = errors.New("create teacher deadline exceeded")
	ErrorFindTeacherByIDDeadlineExceeded error = errors.New("find teacher by ud deadline exceeded")
)

type teacherServiceImpl struct {
	log        *slog.Logger
	repository teacherRepository
}

type teacherService interface {
	Create(ctx context.Context, teacher business.Teacher) (business.TeacherCreateResponse, error)
	FindByID(ctx context.Context, teacherID string) (business.Teacher, error)
}

func NewTeacherService(log *slog.Logger, repository teacherRepository) teacherService {
	return &teacherServiceImpl{
		log:        log,
		repository: repository,
	}
}

type teacherRepository interface {
	Save(context.Context, domain.Teacher) error
	FindByID(context.Context, uuid.UUID) (domain.Teacher, error)
	CheckDuplicateExists(ctx context.Context, reportEmail, username string) (bool, error)
}

func (service *teacherServiceImpl) Create(ctx context.Context, teacher business.Teacher) (business.TeacherCreateResponse, error) {
	const op = "TeacherServiceImpl.Create()"

	log := service.log.With(
		slog.String("op", op),
		slog.String("teacherUsername", teacher.Username),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	contextWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	resultChannel := make(chan business.TeacherCreateResponse)
	errorChannel := make(chan error)

	go func() {
		log.Debug("started creating teacher")
		reportEmailExists, err := service.repository.CheckDuplicateExists(contextWithTimeout, teacher.ReportEmail, teacher.Username)
		if err != nil {
			errorChannel <- handling.HandleApplicationError(err)
			return
		}

		if reportEmailExists {
			log.Error("teacher with this username or report email already exists")
			errorChannel <- handling.WrapAsApplicationError(errors.New("teacher duplicate found"), handling.WithCode(codes.AlreadyExists))
			return
		}

		domainTeacher := converter.ConvertTeacherToDomain(teacher)
		if err := service.repository.Save(contextWithTimeout, domainTeacher); err != nil {
			errorChannel <- handling.HandleApplicationError(err)
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
			return business.TeacherCreateResponse{}, ErrorCreateTeacherDeadlineExceeded
		case createdTeacherData := <-resultChannel:
			return createdTeacherData, nil
		case err := <-errorChannel:
			return business.TeacherCreateResponse{}, err
		}
	}
}

func (service *teacherServiceImpl) FindByID(ctx context.Context, teacherID string) (business.Teacher, error) {
	const op = "TeacherServiceImpl.FindByID()"

	log := service.log.With(
		slog.String("op", op),
		slog.String("teacherID", teacherID),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	contextWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	resultChannel := make(chan business.Teacher)
	errorChannel := make(chan error)

	go func() {
		log.Debug("started finding teacher by id")
		parsedID, err := uuid.Parse(teacherID)
		if err != nil {
			log.Error("error while parsing teacher id - wrong UUID passed")
			errorChannel <- handling.WrapAsApplicationError(err, handling.WithCode(codes.InvalidArgument))
			return
		}

		foundTeacher, err := service.repository.FindByID(contextWithTimeout, parsedID)
		if err != nil {
			errorChannel <- handling.HandleApplicationError(err)
			return
		}

		log.Debug("teacher successfully found by id")
		resultChannel <- converter.ConvertTeacherToBusiness(foundTeacher)
	}()

	for {
		select {
		case <-contextWithTimeout.Done():
			return business.Teacher{}, ErrorFindTeacherByIDDeadlineExceeded
		case foundTeacher := <-resultChannel:
			return foundTeacher, nil
		case err := <-errorChannel:
			return business.Teacher{}, err
		}
	}
}
