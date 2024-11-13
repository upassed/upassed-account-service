package teacher

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/upassed/upassed-account-service/internal/caching/teacher"
	"github.com/upassed/upassed-account-service/internal/config"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"gorm.io/gorm"
	"log/slog"
)

type Repository interface {
	CheckDuplicateExists(ctx context.Context, reportEmail, username string) (bool, error)
	Save(context.Context, *domain.Teacher) error
	FindByUsername(context.Context, string) (*domain.Teacher, error)
}

type teacherRepositoryImpl struct {
	db    *gorm.DB
	cache *teacher.RedisClient
	cfg   *config.Config
	log   *slog.Logger
}

func New(db *gorm.DB, redisClient *redis.Client, cfg *config.Config, log *slog.Logger) Repository {
	cacheClient := teacher.New(redisClient, cfg, log)
	return &teacherRepositoryImpl{
		db:    db,
		cache: cacheClient,
		cfg:   cfg,
		log:   log,
	}
}
