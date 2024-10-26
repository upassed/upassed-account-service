package student

import (
	"github.com/redis/go-redis/v9"
	"github.com/upassed/upassed-account-service/internal/config"
	"log/slog"
)

const keyFormat = "student:%s"

type RedisClient struct {
	cfg    *config.Config
	log    *slog.Logger
	client *redis.Client
}

func New(client *redis.Client, cfg *config.Config, log *slog.Logger) *RedisClient {
	return &RedisClient{
		cfg:    cfg,
		log:    log,
		client: client,
	}
}
