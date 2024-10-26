package migration

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/logging"
	"log/slog"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(cfg *config.Config, log *slog.Logger) error {
	log = logging.Wrap(log,
		logging.WithOp(RunMigrations),
	)

	log.Info("creating a new instance of migrator")
	migrator, err := migrate.New(
		fmt.Sprintf("file://%s", cfg.Migration.MigrationsPath),
		cfg.GetPostgresMigrationConnectionString(),
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
