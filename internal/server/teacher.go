package server

import (
	"context"

	business "github.com/upassed/upassed-account-service/internal/service/model"
	"github.com/upassed/upassed-account-service/pkg/client"
	"google.golang.org/grpc"
)

type teacherServerAPI struct {
	client.UnimplementedTeacherServer
	service teacherService
}

func registerTeacherServer(gRPC *grpc.Server, service teacherService) {
	client.RegisterTeacherServer(gRPC, &teacherServerAPI{
		service: service,
	})
}

type teacherService interface {
	Create(context.Context, business.TeacherCreateRequest) (business.TeacherCreateResponse, error)
}

func (server *teacherServerAPI) Create(ctx context.Context, request *client.TeacherCreateRequest) (*client.TeacherCreateResponse, error) {
	response, _ := server.service.Create(ctx, business.TeacherCreateRequest{})
	return &client.TeacherCreateResponse{
		CreatedTeacherId: response.CreatedTeacherID.String(),
	}, nil
}
