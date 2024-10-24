package repository

import (
	"errors"
	"fmt"
	"github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/migration"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"log/slog"
	"reflect"
	"runtime"
)

var (
	errOpeningDbConnection     = errors.New("failed to open connection to a database")
	errPingingDatabase         = errors.New("failed to ping database")
	errRunningMigrationScripts = errors.New("error while running migration scripts")
)

func OpenGormDbConnection(cfg *config.Config, log *slog.Logger) (*gorm.DB, error) {
	op := runtime.FuncForPC(reflect.ValueOf(OpenGormDbConnection).Pointer()).Name()

	log = log.With(
		slog.String("op", op),
	)

	log.Info("started connecting to postgres database")
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  cfg.GetPostgresConnectionString(),
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})

	if err != nil {
		log.Error("error while opening connection to a database", logging.Error(err))
		return nil, fmt.Errorf("%s - %w", op, errOpeningDbConnection)
	}

	if postgresDB, err := db.DB(); err != nil || postgresDB.Ping() != nil {
		log.Error("error while pinging a database")
		return nil, fmt.Errorf("%s - %w", op, errPingingDatabase)
	}

	log.Info("database connection established successfully")
	if err := migration.RunMigrations(cfg, log); err != nil {
		return nil, errRunningMigrationScripts
	}

	return db, nil
}
