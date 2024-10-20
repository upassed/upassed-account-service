package group

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"runtime"

	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/migration"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var (
	ErrOpeningDbConnection     = errors.New("failed to open connection to a database")
	errPingingDatabase         = errors.New("failed to ping database")
	errRunningMigrationScripts = errors.New("error while running migration scripts")
)

type Repository interface {
	FindStudentsInGroup(context.Context, uuid.UUID) ([]domain.Student, error)
	FindByID(context.Context, uuid.UUID) (domain.Group, error)
	FindByFilter(context.Context, domain.GroupFilter) ([]domain.Group, error)
	Exists(context.Context, uuid.UUID) (bool, error)
}

type groupRepositoryImpl struct {
	db  *gorm.DB
	cfg *config.Config
	log *slog.Logger
}

func New(cfg *config.Config, log *slog.Logger) (Repository, error) {
	op := runtime.FuncForPC(reflect.ValueOf(New).Pointer()).Name()

	log = log.With(
		slog.String("op", op),
	)

	log.Info("started connecting to postgres database")
	postgresInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Storage.Host,
		cfg.Storage.Port,
		cfg.Storage.User,
		cfg.Storage.Password,
		cfg.Storage.DatabaseName,
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  postgresInfo,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})

	if err != nil {
		log.Error("error while opening connection to a database", logging.Error(err))
		return nil, fmt.Errorf("%s - %w", op, ErrOpeningDbConnection)
	}

	if postgresDB, err := db.DB(); err != nil || postgresDB.Ping() != nil {
		log.Error("error while pinging a database")
		return nil, fmt.Errorf("%s - %w", op, errPingingDatabase)
	}

	log.Info("database connection established successfully")
	if err := migration.RunMigrations(cfg, log); err != nil {
		return nil, errRunningMigrationScripts
	}

	return &groupRepositoryImpl{
		db:  db,
		cfg: cfg,
		log: log,
	}, nil
}
