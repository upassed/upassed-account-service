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
)

var (
	ErrSavingStudent = errors.New("error while saving student")
)

func (repository *repositoryImpl) Save(ctx context.Context, student *domain.Student) error {
	spanContext, span := otel.Tracer(repository.cfg.Tracing.StudentTracerName).Start(ctx, "studentRepository#Save")
	span.SetAttributes(attribute.String("username", student.Username))
	defer span.End()

	log := logging.Wrap(repository.log,
		logging.WithOp(repository.Save),
		logging.WithCtx(ctx),
		logging.WithAny("studentUsername", student.Username),
	)

	log.Info("started saving student to a database")
	saveResult := repository.db.WithContext(ctx).Create(&student)
	if err := saveResult.Error; err != nil || saveResult.RowsAffected != 1 {
		log.Error("error while saving student data to a database", logging.Error(err))
		tracing.SetSpanError(span, err)
		return handling.New(ErrSavingStudent.Error(), codes.Internal)
	}

	log.Info("student was successfully inserted into a database")
	log.Info("saving student data into the cache")
	if err := repository.cache.Save(spanContext, student); err != nil {
		log.Error("unable to insert student in cache", logging.Error(err))
		tracing.SetSpanError(span, err)
	}

	log.Info("student was saved to the cache")
	return nil
}
