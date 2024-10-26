package student

import (
	"context"
	"errors"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logging"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"log/slog"
)

var (
	errCountingDuplicatesStudent = errors.New("error while counting duplicate students")
)

func (repository *studentRepositoryImpl) CheckDuplicateExists(ctx context.Context, educationalEmail, username string) (bool, error) {
	_, span := otel.Tracer(repository.cfg.Tracing.StudentTracerName).Start(ctx, "studentRepository#CheckDuplicateExists")
	defer span.End()

	log := logging.Wrap(repository.log,
		logging.WithOp(repository.CheckDuplicateExists),
		logging.WithCtx(ctx),
		logging.WithAny("studentEducationalEmail", educationalEmail),
	)

	log.Info("started checking student duplicates")
	var studentCount int64
	countResult := repository.db.WithContext(ctx).Model(&domain.Student{}).Where("educational_email = ?", educationalEmail).Or("username = ?", username).Count(&studentCount)
	if countResult.Error != nil {
		log.Error("error while counting students with educational_email and username in database")
		return false, handling.New(errCountingDuplicatesStudent.Error(), codes.Internal)
	}

	if studentCount > 0 {
		log.Info("found student duplicates in database", slog.Int64("studentDuplicatesCount", studentCount))
		return true, nil
	}

	log.Info("student duplicates not found in database")
	return false, nil
}
