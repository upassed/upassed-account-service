package student

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/middleware"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
	"log/slog"
	"reflect"
	"runtime"
)

var (
	errSearchingStudentByID = errors.New("error while searching student by id")
	ErrStudentNotFoundByID  = errors.New("student by id not found in database")
)

func (repository *studentRepositoryImpl) FindByID(ctx context.Context, studentID uuid.UUID) (domain.Student, error) {
	op := runtime.FuncForPC(reflect.ValueOf(repository.FindByID).Pointer()).Name()

	log := repository.log.With(
		slog.String("op", op),
		slog.Any("studentID", studentID),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	log.Info("started searching student in a database")
	foundStudent := domain.Student{}
	searchResult := repository.db.WithContext(ctx).Preload("Group").First(&foundStudent, studentID)
	if searchResult.Error != nil {
		if errors.Is(searchResult.Error, gorm.ErrRecordNotFound) {
			log.Error("student was not found in the database", logging.Error(searchResult.Error))
			return domain.Student{}, handling.New(ErrStudentNotFoundByID.Error(), codes.NotFound)
		}

		log.Error("error while searching student in the database", logging.Error(searchResult.Error))
		return domain.Student{}, handling.New(errSearchingStudentByID.Error(), codes.Internal)
	}

	log.Info("student was successfully found in a database")
	return foundStudent, nil
}
