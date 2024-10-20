package migration

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	"github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/logging"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(config *config.Config, log *slog.Logger) error {
	op := runtime.FuncForPC(reflect.ValueOf(RunMigrations).Pointer()).Name()

	migrator, err := migrate.New(
		fmt.Sprintf("file://%s", config.Migration.MigrationsPath),
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable&x-migrations-table=%s",
			config.Storage.User,
			config.Storage.Password,
			config.Storage.Host,
			config.Storage.Port,
			config.Storage.DatabaseName,
			config.Migration.MigrationsTableName,
		),
	)

	log = log.With(
		slog.String("op", op),
	)

	if err != nil {
		log.Error("error while creating migrator", logging.Error(err))
		return err
	}

	log.Info("starting sql migration scripts running")
	if err := migrator.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("no migrations to apply, nothing changed")
			return nil
		}

		log.Error("error while applying migrations", logging.Error(err))
		return err
	}

	log.Info("all migrations applied successfully")
	return nil
}
