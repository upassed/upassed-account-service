package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	config "github.com/upassed/upassed-account-service/internal/config/app"
	"github.com/upassed/upassed-account-service/internal/logger"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrorOpeningDbConnection error = errors.New("failed to open connection to a database")
	ErrorPingingDatabase     error = errors.New("failed to ping database")
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

func (repository *TeacherRepositoryImpl) Save(context.Context, domain.Teacher) error {
	panic("not implemented!")
}
