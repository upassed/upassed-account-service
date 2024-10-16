package student

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/repository/group"
	"github.com/upassed/upassed-account-service/internal/repository/student"
)

type studentService interface {
	Create(ctx context.Context, student Student) (StudentCreateResponse, error)
	FindByID(ctx context.Context, studentID uuid.UUID) (Student, error)
}

type studentServiceImpl struct {
	log               *slog.Logger
	studentRepository studentRepository
	groupRepository   groupRepository
}

type studentRepository interface {
	Save(context.Context, student.Student) error
	CheckDuplicateExists(ctx context.Context, educationalEmail, username string) (bool, error)
	FindByID(context.Context, uuid.UUID) (student.Student, error)
}

type groupRepository interface {
	Exists(context.Context, uuid.UUID) (bool, error)
	FindByID(context.Context, uuid.UUID) (group.Group, error)
}

func New(log *slog.Logger, studentRepository studentRepository, groupgroupRepository groupRepository) studentService {
	return &studentServiceImpl{
		log:               log,
		studentRepository: studentRepository,
		groupRepository:   groupgroupRepository,
	}
}
