package student

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
	errCountingDuplicatesStudent = errors.New("error while counting duplicate students")
)

func (repository *studentRepositoryImpl) CheckDuplicateExists(ctx context.Context, educationalEmail, username string) (bool, error) {
	op := runtime.FuncForPC(reflect.ValueOf(repository.CheckDuplicateExists).Pointer()).Name()

	log := repository.log.With(
		slog.String("op", op),
		slog.String("studentUsername", username),
		slog.String("studentEducationalEmail", educationalEmail),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	log.Debug("started checking student duplicates")
	var studentCount int64
	countResult := repository.db.WithContext(ctx).Model(&domain.Student{}).Where("educational_email = ?", educationalEmail).Or("username = ?", username).Count(&studentCount)
	if countResult.Error != nil {
		log.Error("error while counting students with educational_email and username in database")
		return false, handling.New(errCountingDuplicatesStudent.Error(), codes.Internal)
	}

	if studentCount > 0 {
		log.Debug("found student duplicates in database", slog.Int64("studentDuplicatesCount", studentCount))
		return true, nil
	}

	log.Debug("student duplicates not found in database")
	return false, nil
}
