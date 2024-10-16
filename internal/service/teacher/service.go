package teacher

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	repository "github.com/upassed/upassed-account-service/internal/repository/teacher"
)

type teacherService interface {
	Create(ctx context.Context, teacher Teacher) (TeacherCreateResponse, error)
	FindByID(ctx context.Context, teacherID uuid.UUID) (Teacher, error)
}

type teacherServiceImpl struct {
	log        *slog.Logger
	repository teacherRepository
}

type teacherRepository interface {
	Save(context.Context, repository.Teacher) error
	FindByID(context.Context, uuid.UUID) (repository.Teacher, error)
	CheckDuplicateExists(ctx context.Context, reportEmail, username string) (bool, error)
}

func New(log *slog.Logger, repository teacherRepository) teacherService {
	return &teacherServiceImpl{
		log:        log,
		repository: repository,
	}
}
