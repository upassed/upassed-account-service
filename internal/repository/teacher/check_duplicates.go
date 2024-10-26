package teacher

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
	errCountingDuplicatesTeacher = errors.New("error while counting duplicate teachers")
)

func (repository *teacherRepositoryImpl) CheckDuplicateExists(ctx context.Context, reportEmail, username string) (bool, error) {
	log := logging.Wrap(repository.log,
		logging.WithOp(repository.CheckDuplicateExists),
		logging.WithCtx(ctx),
		logging.WithAny("teacherReportEmail", reportEmail),
	)

	_, span := otel.Tracer(repository.cfg.Tracing.TeacherTracerName).Start(ctx, "teacherRepository#CheckDuplicateExists")
	defer span.End()

	log.Info("started checking teacher duplicates")
	var teacherCount int64
	countResult := repository.db.WithContext(ctx).Model(&domain.Teacher{}).Where("report_email = ?", reportEmail).Or("username = ?", username).Count(&teacherCount)
	if countResult.Error != nil {
		log.Error("error while counting teachers with report_email and username in database")
		return false, handling.New(errCountingDuplicatesTeacher.Error(), codes.Internal)
	}

	if teacherCount > 0 {
		log.Info("found teacher duplicates in database", slog.Int64("teacherDuplicatesCount", teacherCount))
		return true, nil
	}

	log.Info("teacher duplicates not found in database")
	return false, nil
}
