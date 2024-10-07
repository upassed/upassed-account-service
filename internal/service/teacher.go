package service

import (
	"context"
	"log/slog"

	"github.com/upassed/upassed-account-service/internal/handling"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"google.golang.org/grpc/codes"
)

type TeacherServiceImpl struct {
	log *slog.Logger
}

func NewTeacherService(log *slog.Logger) *TeacherServiceImpl {
	return &TeacherServiceImpl{
		log: log,
	}
}

func (service *TeacherServiceImpl) Create(context.Context, business.TeacherCreateRequest) (business.TeacherCreateResponse, error) {
	return business.TeacherCreateResponse{}, handling.NewServiceLayerError("some message", codes.InvalidArgument)
}
