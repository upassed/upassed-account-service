package group

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/logger"
	"github.com/upassed/upassed-account-service/internal/migration"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var (
	ErrorOpeningDbConnection     error = errors.New("failed to open connection to a database")
	ErrorPingingDatabase         error = errors.New("failed to ping database")
	ErrorRunningMigrationScripts error = errors.New("error while running migration scripts")
)

type groupRepository interface {
	FindStudentsInGroup(context.Context, uuid.UUID) ([]domain.Student, error)
	FindByID(context.Context, uuid.UUID) (domain.Group, error)
	FindByFilter(context.Context, domain.GroupFilter) ([]domain.Group, error)
	Exists(context.Context, uuid.UUID) (bool, error)
}

type groupRepositoryImpl struct {
	log *slog.Logger
	db  *gorm.DB
}

func New(config *config.Config, log *slog.Logger) (groupRepository, error) {
	const op = "gorup.New()"

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

	return &groupRepositoryImpl{
		db:  db,
		log: log,
	}, nil
}
