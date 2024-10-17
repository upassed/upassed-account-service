package group

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	business "github.com/upassed/upassed-account-service/internal/service/model"
)

type groupService interface {
	FindStudentsInGroup(context.Context, uuid.UUID) ([]business.Student, error)
	FindByID(context.Context, uuid.UUID) (business.Group, error)
}

type groupServiceImpl struct {
	log        *slog.Logger
	repository groupRepository
}

type groupRepository interface {
	FindStudentsInGroup(context.Context, uuid.UUID) ([]domain.Student, error)
	FindByID(context.Context, uuid.UUID) (domain.Group, error)
}

func New(log *slog.Logger, repository groupRepository) groupService {
	return &groupServiceImpl{
		log:        log,
		repository: repository,
	}
}
