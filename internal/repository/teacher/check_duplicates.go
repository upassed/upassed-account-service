package teacher

import (
	"context"
	"errors"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/middleware"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"google.golang.org/grpc/codes"
	"log/slog"
	"reflect"
	"runtime"
)

var (
	errCountingDuplicatesTeacher = errors.New("error while counting duplicate teachers")
)

func (repository *teacherRepositoryImpl) CheckDuplicateExists(ctx context.Context, reportEmail, username string) (bool, error) {
	op := runtime.FuncForPC(reflect.ValueOf(repository.CheckDuplicateExists).Pointer()).Name()

	log := repository.log.With(
		slog.String("op", op),
		slog.String("teacherUsername", username),
		slog.String("teacherReportEmail", reportEmail),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	log.Debug("started checking teacher duplicates")
	var teacherCount int64
	countResult := repository.db.WithContext(ctx).Model(&domain.Teacher{}).Where("report_email = ?", reportEmail).Or("username = ?", username).Count(&teacherCount)
	if countResult.Error != nil {
		log.Error("error while counting teachers with report_email and username in database")
		return false, handling.New(errCountingDuplicatesTeacher.Error(), codes.Internal)
	}

	if teacherCount > 0 {
		log.Debug("found teacher duplicates in database", slog.Int64("teacherDuplicatesCount", teacherCount))
		return true, nil
	}

	log.Debug("teacher duplicates not found in database")
	return false, nil
}
