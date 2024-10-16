package group

import (
	"context"

	"github.com/google/uuid"
)

func (repository *groupRepositoryImpl) Exists(context.Context, uuid.UUID) (bool, error) {
	panic("not implemented")
}
