package student

import (
	"context"
	"github.com/upassed/upassed-account-service/internal/config"
	"log/slog"

	"github.com/google/uuid"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	business "github.com/upassed/upassed-account-service/internal/service/model"
)

type Service interface {
	Create(ctx context.Context, student *business.Student) (*business.StudentCreateResponse, error)
	FindByUsername(ctx context.Context, studentUsername string) (*business.Student, error)
}

type serviceImpl struct {
	cfg               *config.Config
	log               *slog.Logger
	studentRepository repository
	groupRepository   groupRepository
}

type repository interface {
	CheckDuplicateExists(ctx context.Context, educationalEmail, username string) (bool, error)
	Save(context.Context, *domain.Student) error
	FindByUsername(context.Context, string) (*domain.Student, error)
}

type groupRepository interface {
	Exists(context.Context, uuid.UUID) (bool, error)
	FindByID(context.Context, uuid.UUID) (*domain.Group, error)
}

func New(cfg *config.Config, log *slog.Logger, studentRepository repository, groupRepository groupRepository) Service {
	return &serviceImpl{
		cfg:               cfg,
		log:               log,
		studentRepository: studentRepository,
		groupRepository:   groupRepository,
	}
}
