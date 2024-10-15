package testcontainer

import (
	"context"
	"log/slog"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/migration"
)

type PostgresTestcontainer interface {
	Start(context.Context) (port int, err error)
	Migrate(*config.Config, *slog.Logger) error
	Stop(context.Context) error
}

type postgresTestcontainerImpl struct {
	container testcontainers.Container
}

func NewPostgresTestontainer(ctx context.Context) (PostgresTestcontainer, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	contextWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16.4",
		ExposedPorts: []string{"5432/tcp"},
		AutoRemove:   true,
		Env: map[string]string{
			"POSTGRES_USER":     cfg.Storage.User,
			"POSTGRES_PASSWORD": cfg.Storage.Password,
			"POSTGRES_DB":       cfg.Storage.DatabaseName,
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	postgres, err := testcontainers.GenericContainer(contextWithTimeout, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return nil, err
	}

	return &postgresTestcontainerImpl{
		container: postgres,
	}, nil
}

func (p *postgresTestcontainerImpl) Start(ctx context.Context) (int, error) {
	contextWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	port, err := p.container.MappedPort(contextWithTimeout, "5432")
	if err != nil {
		return 0, err
	}

	return port.Int(), nil
}

func (p *postgresTestcontainerImpl) Migrate(config *config.Config, log *slog.Logger) error {
	return migration.RunMigrations(config, log)
}

func (p *postgresTestcontainerImpl) Stop(ctx context.Context) error {
	contextWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err := p.container.Terminate(contextWithTimeout)
	if err != nil {
		return err
	}

	return nil

}
