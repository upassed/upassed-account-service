package student

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logging"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

var (
	errSearchingStudentByID = errors.New("error while searching student by id")
	ErrStudentNotFoundByID  = errors.New("student by id not found in database")
)

func (repository *studentRepositoryImpl) FindByID(ctx context.Context, studentID uuid.UUID) (*domain.Student, error) {
	log := logging.Wrap(repository.log,
		logging.WithOp(repository.FindByID),
		logging.WithCtx(ctx),
		logging.WithAny("studentID", studentID),
	)

	spanContext, span := otel.Tracer(repository.cfg.Tracing.StudentTracerName).Start(ctx, "studentRepository#FindByID")
	defer span.End()

	log.Info("started searching student by id in redis cache")
	studentFromCache, err := repository.cache.GetByID(spanContext, studentID)
	if err == nil {
		log.Info("student was found in cache, not going to the database")
		return studentFromCache, nil
	}

	log.Info("started searching student in a database")
	foundStudent := domain.Student{}
	searchResult := repository.db.WithContext(ctx).Preload("Group").First(&foundStudent, studentID)
	if searchResult.Error != nil {
		if errors.Is(searchResult.Error, gorm.ErrRecordNotFound) {
			log.Error("student was not found in the database", logging.Error(searchResult.Error))
			return nil, handling.New(ErrStudentNotFoundByID.Error(), codes.NotFound)
		}

		log.Error("error while searching student in the database", logging.Error(searchResult.Error))
		return nil, handling.New(errSearchingStudentByID.Error(), codes.Internal)
	}

	log.Info("student was successfully found in a database")
	log.Info("saving student to cache")
	if err := repository.cache.Save(spanContext, &foundStudent); err != nil {
		log.Error("error while saving student to cache", logging.Error(err))
	}

	return &foundStudent, nil
}
