package teacher

import (
	"context"
	"errors"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logging"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
)

var (
	ErrSavingTeacher = errors.New("error while saving teacher")
)

func (repository *teacherRepositoryImpl) Save(ctx context.Context, teacher *domain.Teacher) error {
	spanContext, span := otel.Tracer(repository.cfg.Tracing.TeacherTracerName).Start(ctx, "teacherRepository#Save")
	defer span.End()

	log := logging.Wrap(repository.log,
		logging.WithOp(repository.Save),
		logging.WithCtx(ctx),
		logging.WithAny("teacherUsername", teacher.Username),
	)

	log.Info("started saving teacher to a database")
	saveResult := repository.db.WithContext(ctx).Create(&teacher)
	if saveResult.Error != nil || saveResult.RowsAffected != 1 {
		log.Error("error while saving teacher data to a database", logging.Error(saveResult.Error))
		return handling.New(ErrSavingTeacher.Error(), codes.Internal)
	}

	log.Info("teacher was successfully inserted into a database")
	log.Info("saving teacher data into the cache")
	if err := repository.cache.Save(spanContext, teacher); err != nil {
		log.Error("unable to insert teacher in cache", logging.Error(err))
	}

	return nil
}
