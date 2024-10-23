package caching

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/logging"
	"log/slog"
	"net"
	"reflect"
	"runtime"
	"strconv"
)

type RedisClient struct {
	cfg    *config.Config
	log    *slog.Logger
	client *redis.Client
}

var (
	errCreatingRedisClient = errors.New("unable to create a redis client")
)

func New(cfg *config.Config, log *slog.Logger) (*RedisClient, error) {
	op := runtime.FuncForPC(reflect.ValueOf(New).Pointer()).Name()

	log = log.With(
		slog.String("op", op),
	)

	databaseNumber, err := strconv.Atoi(cfg.Redis.DatabaseNumber)
	if err != nil {
		log.Error("unable to parse redis database number", logging.Error(err))
		return nil, err
	}

	redisDatabase := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(cfg.Redis.Host, cfg.Redis.Port),
		Username: cfg.Redis.User,
		Password: cfg.Redis.DatabaseNumber,
		DB:       databaseNumber,
	})

	if _, err := redisDatabase.Ping(context.Background()).Result(); err != nil {
		log.Error("unable to ping redis database", logging.Error(err))
		return nil, errCreatingRedisClient
	}

	log.Info("redis client successfully created")
	return &RedisClient{
		cfg:    cfg,
		log:    log,
		client: redisDatabase,
	}, nil
}
