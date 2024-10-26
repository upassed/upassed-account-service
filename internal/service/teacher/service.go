package teacher

import (
	"context"
	"github.com/upassed/upassed-account-service/internal/config"
	"log/slog"

	"github.com/google/uuid"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	business "github.com/upassed/upassed-account-service/internal/service/model"
)

type Service interface {
	Create(ctx context.Context, teacher *business.Teacher) (*business.TeacherCreateResponse, error)
	FindByID(ctx context.Context, teacherID uuid.UUID) (*business.Teacher, error)
}

type teacherServiceImpl struct {
	cfg        *config.Config
	log        *slog.Logger
	repository teacherRepository
}

type teacherRepository interface {
	CheckDuplicateExists(ctx context.Context, reportEmail, username string) (bool, error)
	Save(context.Context, *domain.Teacher) error
	FindByID(context.Context, uuid.UUID) (*domain.Teacher, error)
}

func New(cfg *config.Config, log *slog.Logger, repository teacherRepository) Service {
	return &teacherServiceImpl{
		cfg:        cfg,
		log:        log,
		repository: repository,
	}
}
