package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	config "github.com/upassed/upassed-account-service/internal/config/app"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logger"
	"github.com/upassed/upassed-account-service/internal/middleware"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"google.golang.org/grpc/codes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrorOpeningDbConnection error = errors.New("failed to open connection to a database")
	ErrorPingingDatabase     error = errors.New("failed to ping database")
	ErrorSavingTeacher       error = errors.New("error while saving teacher")
	ErrorTeacherNotFound     error = errors.New("teacher not found in database")
	ErrorSearchingTeacher    error = errors.New("error while searching teacher")
)

type TeacherRepositoryImpl struct {
	log *slog.Logger
	db  *gorm.DB
}

func NewTeacherRepository(config *config.Config, log *slog.Logger) (*TeacherRepositoryImpl, error) {
	const op = "repository.NewTeacherRepository()"

	log = log.With(
		slog.String("op", op),
	)

	log.Info("started connecting to postgres database")
	postgresInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Storage.Host,
		config.Storage.Port,
		config.Storage.User,
		config.Storage.Password,
		config.Storage.DatabaseName,
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  postgresInfo,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		log.Error("error while opening connection to a database", logger.Error(err))
		return nil, fmt.Errorf("%s - %w", op, ErrorOpeningDbConnection)
	}

	if postgresDB, err := db.DB(); err != nil || postgresDB.Ping() != nil {
		log.Error("error while pinging a database")
		return nil, fmt.Errorf("%s - %w", op, ErrorPingingDatabase)
	}

	log.Debug("database connection established successfully")
	return &TeacherRepositoryImpl{
		db:  db,
		log: log,
	}, nil
}

// TODO work with context
func (repository *TeacherRepositoryImpl) Save(ctx context.Context, teacher domain.Teacher) error {
	const op = "repository.TeacherRepositoryImpl.Save()"

	log := repository.log.With(
		slog.String("op", op),
		slog.String("teacherUsername", teacher.Username),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	log.Debug("started saving teacher to a database")
	saveResult := repository.db.Create(&teacher)
	if saveResult.Error != nil || saveResult.RowsAffected != 1 {
		log.Error("error while saving teacher data to a database", logger.Error(saveResult.Error))
		return handling.NewApplicationError(ErrorSavingTeacher.Error(), codes.Internal)
	}

	log.Debug("teacher was successfully inserted into a database")
	return nil
}

func (repository *TeacherRepositoryImpl) FindByID(ctx context.Context, teacherID uuid.UUID) (domain.Teacher, error) {
	const op = "repository.TeacherRepositoryImpl.FindByID()"

	log := repository.log.With(
		slog.String("op", op),
		slog.String("teacherID", teacherID.String()),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	log.Debug("started searching teacher in a database")
	foundTeacher := domain.Teacher{}
	searchResult := repository.db.First(&foundTeacher, teacherID)
	if searchResult.Error != nil {
		if errors.Is(searchResult.Error, gorm.ErrRecordNotFound) {
			log.Error("teacher was not found in the database", logger.Error(searchResult.Error))
			return domain.Teacher{}, handling.NewApplicationError(ErrorTeacherNotFound.Error(), codes.NotFound)
		}

		log.Error("error while searching teacher in the database", logger.Error(searchResult.Error))
		return domain.Teacher{}, handling.NewApplicationError(ErrorSearchingTeacher.Error(), codes.Internal)
	}

	log.Debug("teacher was successfully found in a database")
	return foundTeacher, nil
}

func (repository *TeacherRepositoryImpl) CheckDuplicateExists(ctx context.Context, reportEmail, username string) (bool, error) {
	panic("not implemented")
}
