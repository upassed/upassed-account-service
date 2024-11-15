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
)

type Repository interface {
	FindStudentsInGroup(context.Context, uuid.UUID) ([]*domain.Student, error)
	FindByID(context.Context, uuid.UUID) (*domain.Group, error)
	FindByFilter(context.Context, *domain.GroupFilter) ([]*domain.Group, error)
	Exists(context.Context, uuid.UUID) (bool, error)
}

type repositoryImpl struct {
	db    *gorm.DB
	cache *group.RedisClient
	cfg   *config.Config
	log   *slog.Logger
}

func New(db *gorm.DB, redisClient *redis.Client, cfg *config.Config, log *slog.Logger) Repository {
	cacheClient := group.New(redisClient, cfg, log)
	return &repositoryImpl{
		db:    db,
		cache: cacheClient,
		cfg:   cfg,
		log:   log,
	}
}
