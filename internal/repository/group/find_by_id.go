package group

import (
	"context"

	"github.com/google/uuid"
)

func (repository *groupRepositoryImpl) FindByID(context.Context, uuid.UUID) (Group, error) {
	panic("not implemented")
}
