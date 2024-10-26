package teacher

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/upassed/upassed-account-service/internal/caching/teacher"
	"log/slog"
	"reflect"
	"runtime"

	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/config"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"gorm.io/gorm"
)

type Repository interface {
	CheckDuplicateExists(ctx context.Context, reportEmail, username string) (bool, error)
	Save(context.Context, *domain.Teacher) error
	FindByID(context.Context, uuid.UUID) (*domain.Teacher, error)
}

type teacherRepositoryImpl struct {
	db    *gorm.DB
	cache *teacher.RedisClient
	cfg   *config.Config
	log   *slog.Logger
}

func New(db *gorm.DB, redisClient *redis.Client, cfg *config.Config, log *slog.Logger) Repository {
	op := runtime.FuncForPC(reflect.ValueOf(New).Pointer()).Name()

	log = log.With(
		slog.String("op", op),
	)

	cacheClient := teacher.New(redisClient, cfg, log)
	return &teacherRepositoryImpl{
		db:    db,
		cache: cacheClient,
		cfg:   cfg,
		log:   log,
	}
}
