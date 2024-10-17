package group

import (
	"context"

	"github.com/google/uuid"
	business "github.com/upassed/upassed-account-service/internal/service/model"
)

func (service *groupServiceImpl) FindStudentsInGroup(context.Context, uuid.UUID) ([]business.Student, error) {
	panic("not implemented")
}
