package group

import (
	"context"
	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/config"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"log/slog"
)

type Service interface {
	FindByID(context.Context, uuid.UUID) (*business.Group, error)
	FindStudentsInGroup(context.Context, uuid.UUID) ([]*business.Student, error)
	FindByFilter(context.Context, *business.GroupFilter) ([]*business.Group, error)
}

type serviceImpl struct {
	cfg        *config.Config
	log        *slog.Logger
	repository repository
}

type repository interface {
	FindByID(context.Context, uuid.UUID) (*domain.Group, error)
	FindStudentsInGroup(context.Context, uuid.UUID) ([]*domain.Student, error)
	FindByFilter(context.Context, *domain.GroupFilter) ([]*domain.Group, error)
}

func New(cfg *config.Config, log *slog.Logger, repository repository) Service {
	return &serviceImpl{
		cfg:        cfg,
		log:        log,
		repository: repository,
	}
}
