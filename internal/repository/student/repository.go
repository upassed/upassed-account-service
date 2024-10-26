package student

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/upassed/upassed-account-service/internal/caching/student"
	"github.com/upassed/upassed-account-service/internal/config"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"gorm.io/gorm"
	"log/slog"
)

type Repository interface {
	CheckDuplicateExists(ctx context.Context, educationalEmail, username string) (bool, error)
	Save(context.Context, *domain.Student) error
	FindByID(context.Context, uuid.UUID) (*domain.Student, error)
}

type studentRepositoryImpl struct {
	db    *gorm.DB
	cache *student.RedisClient
	cfg   *config.Config
	log   *slog.Logger
}

func New(db *gorm.DB, redisClient *redis.Client, cfg *config.Config, log *slog.Logger) Repository {
	cacheClient := student.New(redisClient, cfg, log)
	return &studentRepositoryImpl{
		db:    db,
		cache: cacheClient,
		cfg:   cfg,
		log:   log,
	}
}
