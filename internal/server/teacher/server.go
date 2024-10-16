package teacher

import (
	"context"

	"github.com/google/uuid"
	service "github.com/upassed/upassed-account-service/internal/service/teacher"
	"github.com/upassed/upassed-account-service/pkg/client"
	"google.golang.org/grpc"
)

type teacherServerAPI struct {
	client.UnimplementedTeacherServer
	service teacherService
}

type teacherService interface {
	Create(context.Context, service.Teacher) (service.TeacherCreateResponse, error)
	FindByID(ctx context.Context, teacherID uuid.UUID) (service.Teacher, error)
}

func Register(gRPC *grpc.Server, service teacherService) {
	client.RegisterTeacherServer(gRPC, &teacherServerAPI{
		service: service,
	})
}