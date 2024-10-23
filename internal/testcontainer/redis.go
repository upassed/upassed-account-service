package testcontainer

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/upassed/upassed-account-service/internal/config"
	"time"
)

type RedisTestcontainer interface {
	Start(context.Context) (port int, err error)
	Stop(context.Context) error
}

type redisTestcontainerImpl struct {
	container testcontainers.Container
}

func NewRedisTestcontainer(ctx context.Context, cfg *config.Config) (RedisTestcontainer, error) {
	contextWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	req := testcontainers.ContainerRequest{
		Image:        "redis:7.4",
		ExposedPorts: []string{"6379/tcp"},
		HostConfigModifier: func(cfg *container.HostConfig) {
			cfg.AutoRemove = true
		},
		Env: map[string]string{
			"REDIS_PASSWORD": cfg.Redis.Password,
		},
		WaitingFor: wait.ForListeningPort("6379/tcp"),
	}

	redis, err := testcontainers.GenericContainer(contextWithTimeout, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return nil, err
	}

	return &redisTestcontainerImpl{
		container: redis,
	}, nil
}

func (r *redisTestcontainerImpl) Start(ctx context.Context) (int, error) {
	contextWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	port, err := r.container.MappedPort(contextWithTimeout, "6379")
	if err != nil {
		return 0, err
	}

	return port.Int(), nil
}

func (r *redisTestcontainerImpl) Stop(ctx context.Context) error {
	contextWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err := r.container.Terminate(contextWithTimeout)
	if err != nil {
		return err
	}

	return nil
}
