package group

import (
	"context"
	"github.com/upassed/upassed-account-service/internal/config"
	"log/slog"

	"github.com/google/uuid"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	business "github.com/upassed/upassed-account-service/internal/service/model"
)

type Service interface {
	FindStudentsInGroup(context.Context, uuid.UUID) ([]business.Student, error)
	FindByID(context.Context, uuid.UUID) (business.Group, error)
	FindByFilter(context.Context, business.GroupFilter) ([]business.Group, error)
}

type groupServiceImpl struct {
	cfg        *config.Config
	log        *slog.Logger
	repository groupRepository
}

type groupRepository interface {
	FindStudentsInGroup(context.Context, uuid.UUID) ([]domain.Student, error)
	FindByID(context.Context, uuid.UUID) (domain.Group, error)
	FindByFilter(context.Context, domain.GroupFilter) ([]domain.Group, error)
}

func New(cfg *config.Config, log *slog.Logger, repository groupRepository) Service {
	return &groupServiceImpl{
		cfg:        cfg,
		log:        log,
		repository: repository,
	}
}
