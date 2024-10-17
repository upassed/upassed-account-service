package group

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	business "github.com/upassed/upassed-account-service/internal/service/model"
)

type groupService interface {
	FindStudentsInGroup(context.Context, uuid.UUID) ([]business.Student, error)
}

type groupServiceImpl struct {
	log        *slog.Logger
	repository groupRepository
}

type groupRepository interface {
}

func New(log *slog.Logger, repository groupRepository) groupService {
	return &groupServiceImpl{
		log:        log,
		repository: repository,
	}
}
