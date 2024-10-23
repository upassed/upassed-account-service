package student

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
	ErrSavingStudent = errors.New("error while saving student")
)

func (repository *studentRepositoryImpl) Save(ctx context.Context, student domain.Student) error {
	op := runtime.FuncForPC(reflect.ValueOf(repository.Save).Pointer()).Name()

	spanContext, span := otel.Tracer(repository.cfg.Tracing.StudentTracerName).Start(ctx, "studentRepository#Save")
	defer span.End()

	log := repository.log.With(
		slog.String("op", op),
		slog.String("studentUsername", student.Username),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	log.Info("started saving student to a database")
	saveResult := repository.db.WithContext(ctx).Create(&student)
	if saveResult.Error != nil || saveResult.RowsAffected != 1 {
		log.Error("error while saving student data to a database", logging.Error(saveResult.Error))
		return handling.New(ErrSavingStudent.Error(), codes.Internal)
	}

	log.Info("student was successfully inserted into a database")
	log.Info("saving student data into the cache")

	if err := repository.cache.SaveStudent(spanContext, student); err != nil {
		log.Error("unable to insert student in cache", logging.Error(err))
	}

	return nil
}
