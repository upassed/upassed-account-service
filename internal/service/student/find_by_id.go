package student

import (
	"context"

	"github.com/google/uuid"
)

func (service *studentServiceImpl) FindByID(ctx context.Context, studentID uuid.UUID) (Student, error) {
	panic("implement me")
}
