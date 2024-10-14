package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	config "github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logger"
	"github.com/upassed/upassed-account-service/internal/middleware"
	"github.com/upassed/upassed-account-service/internal/migration"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"google.golang.org/grpc/codes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var (
	ErrorOpeningDbConnection       error = errors.New("failed to open connection to a database")
	ErrorPingingDatabase           error = errors.New("failed to ping database")
	ErrorSavingTeacher             error = errors.New("error while saving teacher")
	ErrorTeacherNotFound           error = errors.New("teacher not found in database")
	ErrorSearchingTeacher          error = errors.New("error while searching teacher")
	ErrorCountingDuplicatesTeacher error = errors.New("error while counting duplicates teacher")
	ErrorRunningMigrationScripts   error = errors.New("error while running migration scripts")

	ErrorSaveTeacherDeadlineExceeded            error = errors.New("saving teacher into a database deadline exceeded")
	ErrorFindTeacherByIDDeadlineExceeded        error = errors.New("finding teacher by id in a database deadline exceeded")
	ErrorCheckTeacherDuplicatesDeadlineExceeded error = errors.New("checking teacher duplicates in a database deadline exceeded")
)

type teacherRepositoryImpl struct {
	log *slog.Logger
	db  *gorm.DB
}

type teacherRepository interface {
	Save(context.Context, domain.Teacher) error
	FindByID(context.Context, uuid.UUID) (domain.Teacher, error)
	CheckDuplicateExists(ctx context.Context, reportEmail, username string) (bool, error)
}

func NewTeacherRepository(config *config.Config, log *slog.Logger) (teacherRepository, error) {
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
	}), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})

	if err != nil {
		log.Error("error while opening connection to a database", logger.Error(err))
		return nil, fmt.Errorf("%s - %w", op, ErrorOpeningDbConnection)
	}

	if postgresDB, err := db.DB(); err != nil || postgresDB.Ping() != nil {
		log.Error("error while pinging a database")
		return nil, fmt.Errorf("%s - %w", op, ErrorPingingDatabase)
	}

	log.Debug("database connection established successfully")
	if err := migration.RunMigrations(config, log); err != nil {
		return nil, ErrorRunningMigrationScripts
	}

	return &teacherRepositoryImpl{
		db:  db,
		log: log,
	}, nil
}

func (repository *teacherRepositoryImpl) Save(ctx context.Context, teacher domain.Teacher) error {
	const op = "repository.TeacherRepositoryImpl.Save()"

	log := repository.log.With(
		slog.String("op", op),
		slog.String("teacherUsername", teacher.Username),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	contextWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	resultChannel := make(chan struct{})
	errorChannel := make(chan error)

	go func() {
		log.Debug("started saving teacher to a database")
		saveResult := repository.db.Create(&teacher)
		if saveResult.Error != nil || saveResult.RowsAffected != 1 {
			log.Error("error while saving teacher data to a database", logger.Error(saveResult.Error))
			errorChannel <- handling.NewApplicationError(ErrorSavingTeacher.Error(), codes.Internal)
			return
		}

		log.Debug("teacher was successfully inserted into a database")
		resultChannel <- struct{}{}
	}()

	for {
		select {
		case <-contextWithTimeout.Done():
			return ErrorSaveTeacherDeadlineExceeded
		case <-resultChannel:
			return nil
		case err := <-errorChannel:
			return err
		}
	}
}

func (repository *teacherRepositoryImpl) FindByID(ctx context.Context, teacherID uuid.UUID) (domain.Teacher, error) {
	const op = "repository.TeacherRepositoryImpl.FindByID()"

	log := repository.log.With(
		slog.String("op", op),
		slog.String("teacherID", teacherID.String()),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	contextWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	resultChannel := make(chan domain.Teacher)
	errorChannel := make(chan error)

	go func() {
		log.Debug("started searching teacher in a database")
		foundTeacher := domain.Teacher{}
		searchResult := repository.db.First(&foundTeacher, teacherID)
		if searchResult.Error != nil {
			if errors.Is(searchResult.Error, gorm.ErrRecordNotFound) {
				log.Error("teacher was not found in the database", logger.Error(searchResult.Error))
				errorChannel <- handling.NewApplicationError(ErrorTeacherNotFound.Error(), codes.NotFound)
				return
			}

			log.Error("error while searching teacher in the database", logger.Error(searchResult.Error))
			errorChannel <- handling.NewApplicationError(ErrorSearchingTeacher.Error(), codes.Internal)
			return
		}

		log.Debug("teacher was successfully found in a database")
		resultChannel <- foundTeacher
	}()

	for {
		select {
		case <-contextWithTimeout.Done():
			return domain.Teacher{}, ErrorFindTeacherByIDDeadlineExceeded
		case foundTeacher := <-resultChannel:
			return foundTeacher, nil
		case err := <-errorChannel:
			return domain.Teacher{}, err
		}
	}
}

func (repository *teacherRepositoryImpl) CheckDuplicateExists(ctx context.Context, reportEmail, username string) (bool, error) {
	const op = "repository.TeacherRepositoryImpl.CheckDuplicateExists()"

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
		countResult := repository.db.Model(&domain.Teacher{}).Where("report_email = ?", reportEmail).Or("username = ?", username).Count(&teacherCount)
		if countResult.Error != nil {
			log.Error("error while counting teachers with report_email and username in database")
			errorChannel <- handling.NewApplicationError(ErrorCountingDuplicatesTeacher.Error(), codes.Internal)
			return
		}

		if teacherCount > 0 {
			log.Debug("found teacher duplicates in database", slog.Int64("teacherDuplicatesCouint", teacherCount))
			resultChannel <- true
			return
		}

		log.Debug("teacher duplicates not found in database")
		resultChannel <- false
		return
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
