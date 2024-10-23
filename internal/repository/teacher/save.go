package teacher

import (
	"context"
	"errors"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/middleware"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"log/slog"
	"reflect"
	"runtime"
)

var (
	ErrSavingTeacher = errors.New("error while saving teacher")
)

func (repository *teacherRepositoryImpl) Save(ctx context.Context, teacher domain.Teacher) error {
	op := runtime.FuncForPC(reflect.ValueOf(repository.Save).Pointer()).Name()

	log := repository.log.With(
		slog.String("op", op),
		slog.String("teacherUsername", teacher.Username),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	spanContext, span := otel.Tracer(repository.cfg.Tracing.TeacherTracerName).Start(ctx, "teacherRepository#Save")
	defer span.End()

	log.Info("started saving teacher to a database")
	saveResult := repository.db.WithContext(ctx).Create(&teacher)
	if saveResult.Error != nil || saveResult.RowsAffected != 1 {
		log.Error("error while saving teacher data to a database", logging.Error(saveResult.Error))
		return handling.New(ErrSavingTeacher.Error(), codes.Internal)
	}

	log.Info("teacher was successfully inserted into a database")
	log.Info("saving teacher data into the cache")
	if err := repository.cache.SaveTeacher(spanContext, teacher); err != nil {
		log.Error("unable to insert teacher in cache", logging.Error(err))
	}

	return nil
}
