package student

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/upassed/upassed-account-service/internal/caching/student"
	"log/slog"
	"reflect"
	"runtime"

	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/config"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"gorm.io/gorm"
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
	op := runtime.FuncForPC(reflect.ValueOf(New).Pointer()).Name()

	log = log.With(
		slog.String("op", op),
	)

	cacheClient := student.New(redisClient, cfg, log)
	return &studentRepositoryImpl{
		db:    db,
		cache: cacheClient,
		cfg:   cfg,
		log:   log,
	}
}
