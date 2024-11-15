package teacher

import (
	"context"
	"github.com/upassed/upassed-account-service/internal/config"
	"log/slog"

	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	business "github.com/upassed/upassed-account-service/internal/service/model"
)

type Service interface {
	Create(ctx context.Context, teacher *business.Teacher) (*business.TeacherCreateResponse, error)
	FindByUsername(ctx context.Context, teacherUsername string) (*business.Teacher, error)
}

type serviceImpl struct {
	cfg        *config.Config
	log        *slog.Logger
	repository repository
}

type repository interface {
	CheckDuplicateExists(ctx context.Context, reportEmail, username string) (bool, error)
	Save(context.Context, *domain.Teacher) error
	FindByUsername(context.Context, string) (*domain.Teacher, error)
}

func New(cfg *config.Config, log *slog.Logger, repository repository) Service {
	return &serviceImpl{
		cfg:        cfg,
		log:        log,
		repository: repository,
	}
}
