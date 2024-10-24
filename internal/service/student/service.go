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
	Create(ctx context.Context, student business.Student) (business.StudentCreateResponse, error)
	FindByID(ctx context.Context, studentID uuid.UUID) (business.Student, error)
}

type studentServiceImpl struct {
	cfg               *config.Config
	log               *slog.Logger
	studentRepository studentRepository
	groupRepository   groupRepository
}

type studentRepository interface {
	Save(context.Context, domain.Student) error
	CheckDuplicateExists(ctx context.Context, educationalEmail, username string) (bool, error)
	FindByID(context.Context, uuid.UUID) (domain.Student, error)
}

type groupRepository interface {
	Exists(context.Context, uuid.UUID) (bool, error)
	FindByID(context.Context, uuid.UUID) (domain.Group, error)
}

func New(cfg *config.Config, log *slog.Logger, studentRepository studentRepository, groupRepository groupRepository) Service {
	return &studentServiceImpl{
		cfg:               cfg,
		log:               log,
		studentRepository: studentRepository,
		groupRepository:   groupRepository,
	}
}
