package group

import (
	"context"

	business "github.com/upassed/upassed-account-service/internal/service/model"
)

func (service *groupServiceImpl) FindByFilter(ctx context.Context, filter business.GroupFilter) ([]business.Group, error) {
	panic("not implemented")
}
