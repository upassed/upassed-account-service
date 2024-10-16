package teacher

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/middleware"
	"google.golang.org/grpc/codes"
)

var (
	ErrorCountingDuplicatesTeacher              error = errors.New("error while counting duplicate teachers")
	ErrorCheckTeacherDuplicatesDeadlineExceeded error = errors.New("checking teacher duplicates in a database deadline exceeded")
)

func (repository *teacherRepositoryImpl) CheckDuplicateExists(ctx context.Context, reportEmail, username string) (bool, error) {
	const op = "teacher.teacherRepositoryImpl.CheckDuplicateExists()"

	log := repository.log.With(
		slog.String("op", op),
		slog.String("teacherUsername", username),
		slog.String("teacherReportEmail", reportEmail),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	contextWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	resultChannel := make(chan bool)
	errorChannel := make(chan error)

	go func() {
		log.Debug("started checking teacher duplicates")
		var teacherCount int64
		countResult := repository.db.Model(&Teacher{}).Where("report_email = ?", reportEmail).Or("username = ?", username).Count(&teacherCount)
		if countResult.Error != nil {
			log.Error("error while counting teachers with report_email and username in database")
			errorChannel <- handling.New(ErrorCountingDuplicatesTeacher.Error(), codes.Internal)
			return
		}

		if teacherCount > 0 {
			log.Debug("found teacher duplicates in database", slog.Int64("teacherDuplicatesCouint", teacherCount))
			resultChannel <- true
			return
		}

		log.Debug("teacher duplicates not found in database")
		resultChannel <- false
	}()

	for {
		select {
		case <-contextWithTimeout.Done():
			return false, ErrorCheckTeacherDuplicatesDeadlineExceeded
		case duplicatesFound := <-resultChannel:
			return duplicatesFound, nil
		case err := <-errorChannel:
			return false, err
		}
	}
}
