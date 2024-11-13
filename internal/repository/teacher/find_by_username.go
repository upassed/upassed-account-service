package teacher

import (
	"context"
	"errors"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logging"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"github.com/upassed/upassed-account-service/internal/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

var (
	errSearchingTeacherByUsername = errors.New("error while searching teacher by username")
	ErrTeacherNotFoundByUsername  = errors.New("teacher by username  not found in database")
)

func (repository *teacherRepositoryImpl) FindByUsername(ctx context.Context, teacherUsername string) (*domain.Teacher, error) {
	spanContext, span := otel.Tracer(repository.cfg.Tracing.TeacherTracerName).Start(ctx, "teacherRepository#FindByUsername")
	span.SetAttributes(attribute.String("teacherUsername", teacherUsername))
	defer span.End()

	log := logging.Wrap(repository.log,
		logging.WithOp(repository.FindByUsername),
		logging.WithCtx(ctx),
		logging.WithAny("teacherUsername", teacherUsername),
	)

	log.Info("started searching teacher by username in redis cache")
	teacherFromCache, err := repository.cache.GetByUsername(spanContext, teacherUsername)
	if err == nil {
		log.Info("teacher was found in cache, not going to the database")
		return teacherFromCache, nil
	}

	log.Info("started searching teacher by username in a database")
	foundTeacher := domain.Teacher{}
	searchResult := repository.db.WithContext(ctx).Where("username = ?", teacherUsername).First(&foundTeacher)
	if err := searchResult.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("teacher by username was not found in the database", logging.Error(err))
			tracing.SetSpanError(span, err)
			return nil, handling.New(ErrTeacherNotFoundByUsername.Error(), codes.NotFound)
		}

		log.Error("error while searching teacher by username the database", logging.Error(err))
		tracing.SetSpanError(span, err)
		return nil, handling.New(errSearchingTeacherByUsername.Error(), codes.Internal)
	}

	log.Info("teacher by username was successfully found in a database")
	log.Info("saving teacher by username to cache")
	if err := repository.cache.Save(spanContext, &foundTeacher); err != nil {
		log.Error("error while saving teacher to cache", logging.Error(err))
		tracing.SetSpanError(span, err)
	}

	log.Info("teacher was saved to the cache by username")
	return &foundTeacher, nil
}
