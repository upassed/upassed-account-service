package teacher

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	business "github.com/upassed/upassed-account-service/internal/service/model"
)

type teacherService interface {
	Create(ctx context.Context, teacher business.Teacher) (business.TeacherCreateResponse, error)
	FindByID(ctx context.Context, teacherID uuid.UUID) (business.Teacher, error)
}

type teacherServiceImpl struct {
	log        *slog.Logger
	repository teacherRepository
}

type teacherRepository interface {
	Save(context.Context, domain.Teacher) error
	FindByID(context.Context, uuid.UUID) (domain.Teacher, error)
	CheckDuplicateExists(ctx context.Context, reportEmail, username string) (bool, error)
}

func New(log *slog.Logger, repository teacherRepository) teacherService {
	return &teacherServiceImpl{
		log:        log,
		repository: repository,
	}
}
