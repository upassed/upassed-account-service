package service

import (
	"context"

	"github.com/google/uuid"
	business "github.com/upassed/upassed-account-service/internal/service/model"
)

type TeacherServiceImpl struct {
}

func NewTeacherService() *TeacherServiceImpl {
	return &TeacherServiceImpl{}
}

func (service *TeacherServiceImpl) Create(context.Context, business.TeacherCreateRequest) (business.TeacherCreateResponse, error) {
	return business.TeacherCreateResponse{
		CreatedTeacherID: uuid.New(),
	}, nil
}
