package group

import (
	"context"

	"github.com/google/uuid"
	business "github.com/upassed/upassed-account-service/internal/service/model"
)

func (service *groupServiceImpl) FindByID(ctx context.Context, groupID uuid.UUID) (business.Group, error) {
	panic("not implemented")
}
