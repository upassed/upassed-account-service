package teacher

import (
	"github.com/redis/go-redis/v9"
	"github.com/upassed/upassed-account-service/internal/config"
	"log/slog"
)

type RedisClient struct {
	cfg    *config.Config
	log    *slog.Logger
	client *redis.Client
}

const (
	usernameKeyFormat = "teacherUsername:%s"
)

func New(client *redis.Client, cfg *config.Config, log *slog.Logger) *RedisClient {
	return &RedisClient{
		cfg:    cfg,
		log:    log,
		client: client,
	}
}
