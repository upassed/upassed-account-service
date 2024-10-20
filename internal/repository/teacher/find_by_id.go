package teacher

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
	errSearchingTeacherByID = errors.New("error while searching teacher by id")
	ErrTeacherNotFoundByID  = errors.New("teacher by id  not found in database")
)

func (repository *teacherRepositoryImpl) FindByID(ctx context.Context, teacherID uuid.UUID) (domain.Teacher, error) {
	op := runtime.FuncForPC(reflect.ValueOf(repository.FindByID).Pointer()).Name()

	log := repository.log.With(
		slog.String("op", op),
		slog.Any("teacherID", teacherID),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	log.Debug("started searching teacher in a database")
	foundTeacher := domain.Teacher{}
	searchResult := repository.db.WithContext(ctx).First(&foundTeacher, teacherID)
	if searchResult.Error != nil {
		if errors.Is(searchResult.Error, gorm.ErrRecordNotFound) {
			log.Error("teacher was not found in the database", logging.Error(searchResult.Error))
			return domain.Teacher{}, handling.New(ErrTeacherNotFoundByID.Error(), codes.NotFound)
		}

		log.Error("error while searching teacher in the database", logging.Error(searchResult.Error))
		return domain.Teacher{}, handling.New(errSearchingTeacherByID.Error(), codes.Internal)
	}

	log.Debug("teacher was successfully found in a database")
	return foundTeacher, nil
}
