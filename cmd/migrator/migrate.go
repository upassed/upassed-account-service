package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	appConfig "github.com/upassed/upassed-account-service/internal/config/app"
	config "github.com/upassed/upassed-account-service/internal/config/migrator"
	"github.com/upassed/upassed-account-service/internal/logger"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Fatalf("error while loading migrator config: %s", err.Error())
	}

	migrator, err := migrate.New(
		fmt.Sprintf("file://%s", config.MigrationsPath),
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable&x-migrations-table=%s",
			config.Storage.User,
			config.Storage.Password,
			config.Storage.Host,
			config.Storage.Port,
			config.Storage.DatabaseName,
			config.MigrationsTableName,
		),
	)

	log := logger.New(appConfig.EnvLocal)

	if err != nil {
		log.Error("error while creating migrator", logger.Error(err))
		os.Exit(1)
	}

	if err := migrator.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("no migrations to apply, nothing changed")
			return
		}

		log.Error("error while applying migrations", logger.Error(err))
		os.Exit(1)
	}

	log.Info("all migrations applied successfully")
}
