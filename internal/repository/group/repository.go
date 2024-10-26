package group

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/upassed/upassed-account-service/internal/caching/group"
	"github.com/upassed/upassed-account-service/internal/config"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"gorm.io/gorm"
	"log/slog"
	"reflect"
	"runtime"
)

type Repository interface {
	FindStudentsInGroup(context.Context, uuid.UUID) ([]*domain.Student, error)
	FindByID(context.Context, uuid.UUID) (*domain.Group, error)
	FindByFilter(context.Context, *domain.GroupFilter) ([]*domain.Group, error)
	Exists(context.Context, uuid.UUID) (bool, error)
}

type groupRepositoryImpl struct {
	db    *gorm.DB
	cache *group.RedisClient
	cfg   *config.Config
	log   *slog.Logger
}

func New(db *gorm.DB, redisClient *redis.Client, cfg *config.Config, log *slog.Logger) Repository {
	op := runtime.FuncForPC(reflect.ValueOf(New).Pointer()).Name()

	log = log.With(
		slog.String("op", op),
	)

	cacheClient := group.New(redisClient, cfg, log)
	return &groupRepositoryImpl{
		db:    db,
		cache: cacheClient,
		cfg:   cfg,
		log:   log,
	}
}
