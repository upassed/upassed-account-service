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
)

var (
	ErrSavingTeacher = errors.New("error while saving teacher")
)

func (repository *repositoryImpl) Save(ctx context.Context, teacher *domain.Teacher) error {
	spanContext, span := otel.Tracer(repository.cfg.Tracing.TeacherTracerName).Start(ctx, "teacherRepository#Save")
	span.SetAttributes(attribute.String("username", teacher.Username))
	defer span.End()

	log := logging.Wrap(repository.log,
		logging.WithOp(repository.Save),
		logging.WithCtx(ctx),
		logging.WithAny("teacherUsername", teacher.Username),
	)

	log.Info("started saving teacher to a database")
	saveResult := repository.db.WithContext(ctx).Create(&teacher)
	if err := saveResult.Error; err != nil || saveResult.RowsAffected != 1 {
		log.Error("error while saving teacher data to a database", logging.Error(err))
		tracing.SetSpanError(span, err)
		return handling.New(ErrSavingTeacher.Error(), codes.Internal)
	}

	log.Info("teacher was successfully inserted into a database")
	log.Info("saving teacher data into the cache")
	if err := repository.cache.Save(spanContext, teacher); err != nil {
		log.Error("unable to insert teacher in cache", logging.Error(err))
		tracing.SetSpanError(span, err)
	}

	log.Info("teacher was saved to the cache")
	return nil
}
