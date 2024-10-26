package teacher

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
	errSearchingTeacherByID = errors.New("error while searching teacher by id")
	ErrTeacherNotFoundByID  = errors.New("teacher by id  not found in database")
)

func (repository *teacherRepositoryImpl) FindByID(ctx context.Context, teacherID uuid.UUID) (*domain.Teacher, error) {
	spanContext, span := otel.Tracer(repository.cfg.Tracing.TeacherTracerName).Start(ctx, "teacherRepository#FindByID")
	defer span.End()

	log := logging.Wrap(repository.log,
		logging.WithOp(repository.FindByID),
		logging.WithCtx(ctx),
		logging.WithAny("teacherID", teacherID),
	)

	log.Info("started searching teacher by id in redis cache")
	teacherFromCache, err := repository.cache.GetByID(spanContext, teacherID)
	if err == nil {
		log.Info("teacher was found in cache, not going to the database")
		return teacherFromCache, nil
	}

	log.Info("started searching teacher by id in a database")
	foundTeacher := domain.Teacher{}
	searchResult := repository.db.WithContext(ctx).First(&foundTeacher, teacherID)
	if searchResult.Error != nil {
		if errors.Is(searchResult.Error, gorm.ErrRecordNotFound) {
			log.Error("teacher was not found in the database", logging.Error(searchResult.Error))
			return nil, handling.New(ErrTeacherNotFoundByID.Error(), codes.NotFound)
		}

		log.Error("error while searching teacher in the database", logging.Error(searchResult.Error))
		return nil, handling.New(errSearchingTeacherByID.Error(), codes.Internal)
	}

	log.Info("teacher was successfully found in a database")
	log.Info("saving teacher to cache")
	if err := repository.cache.Save(spanContext, &foundTeacher); err != nil {
		log.Error("error while saving teacher to cache", logging.Error(err))
	}

	return &foundTeacher, nil
}
