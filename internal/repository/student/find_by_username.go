package student

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
	errSearchingStudentByUsername = errors.New("error while searching student by username")
	ErrStudentNotFoundByUsername  = errors.New("student by username not found in database")
)

func (repository *studentRepositoryImpl) FindByUsername(ctx context.Context, studentUsername string) (*domain.Student, error) {
	spanContext, span := otel.Tracer(repository.cfg.Tracing.StudentTracerName).Start(ctx, "studentRepository#FindByUsername")
	span.SetAttributes(attribute.String("studentUsername", studentUsername))
	defer span.End()

	log := logging.Wrap(repository.log,
		logging.WithOp(repository.FindByUsername),
		logging.WithCtx(ctx),
		logging.WithAny("studentUsername", studentUsername),
	)

	log.Info("started searching student by username in redis cache")
	studentFromCache, err := repository.cache.GetByUsername(spanContext, studentUsername)
	if err == nil {
		log.Info("student was found in cache by username, not going to the database")
		return studentFromCache, nil
	}

	log.Info("started searching student by username in a database")
	foundStudent := domain.Student{}
	searchResult := repository.db.WithContext(ctx).Preload("Group").Where("username = ?", studentUsername).First(&foundStudent)
	if err := searchResult.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("student was not found in the database", logging.Error(err))
			tracing.SetSpanError(span, err)
			return nil, handling.New(ErrStudentNotFoundByUsername.Error(), codes.NotFound)
		}

		log.Error("error while searching student in the database", logging.Error(err))
		tracing.SetSpanError(span, err)
		return nil, handling.New(errSearchingStudentByUsername.Error(), codes.Internal)
	}

	log.Info("student was successfully found in a database")
	log.Info("saving student to cache")
	if err := repository.cache.Save(spanContext, &foundStudent); err != nil {
		log.Error("error while saving student to cache", logging.Error(err))
		tracing.SetSpanError(span, err)
	}

	log.Info("student was successfully saved to the cache")
	return &foundStudent, nil
}
