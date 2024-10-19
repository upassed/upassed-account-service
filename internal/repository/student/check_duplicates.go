package student

import (
	"context"
	"errors"
	"log/slog"
	"reflect"
	"runtime"
	"time"

	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/middleware"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"google.golang.org/grpc/codes"
)

var (
	errCountingDuplicatesStudent              = errors.New("error while counting duplicate students")
	errCheckStudentDuplicatesDeadlineExceeded = errors.New("checking student duplicates in a database deadline exceeded")
)

func (repository *studentRepositoryImpl) CheckDuplicateExists(ctx context.Context, educationalEmail, username string) (bool, error) {
	op := runtime.FuncForPC(reflect.ValueOf(repository.CheckDuplicateExists).Pointer()).Name()

	log := repository.log.With(
		slog.String("op", op),
		slog.String("studentUsername", username),
		slog.String("studentEducationalEmail", educationalEmail),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	contextWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	resultChannel := make(chan bool)
	errorChannel := make(chan error)

	go func() {
		log.Debug("started checking student duplicates")
		var studentCount int64
		countResult := repository.db.Model(&domain.Student{}).Where("educational_email = ?", educationalEmail).Or("username = ?", username).Count(&studentCount)
		if countResult.Error != nil {
			log.Error("error while counting students with educational_email and username in database")
			errorChannel <- handling.New(errCountingDuplicatesStudent.Error(), codes.Internal)
			return
		}

		if studentCount > 0 {
			log.Debug("found student duplicates in database", slog.Int64("studentDuplicatesCount", studentCount))
			resultChannel <- true
			return
		}

		log.Debug("student duplicates not found in database")
		resultChannel <- false
	}()

	for {
		select {
		case <-contextWithTimeout.Done():
			return false, errCheckStudentDuplicatesDeadlineExceeded
		case duplicatesFound := <-resultChannel:
			return duplicatesFound, nil
		case err := <-errorChannel:
			return false, err
		}
	}
}
