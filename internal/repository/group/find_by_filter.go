package group

import (
	"context"

	domain "github.com/upassed/upassed-account-service/internal/repository/model"
)

func (repository *groupRepositoryImpl) FindByFilter(ctx context.Context, filter domain.GroupFilter) ([]domain.Group, error) {
	panic("not implemented")
}
