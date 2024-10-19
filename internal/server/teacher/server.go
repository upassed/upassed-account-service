package teacher

import (
	"context"

	"github.com/google/uuid"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"github.com/upassed/upassed-account-service/pkg/client"
	"google.golang.org/grpc"
)

type teacherServerAPI struct {
	client.UnimplementedTeacherServer
	service teacherService
}

type teacherService interface {
	Create(context.Context, business.Teacher) (business.TeacherCreateResponse, error)
	FindByID(ctx context.Context, teacherID uuid.UUID) (business.Teacher, error)
}

func Register(gRPC *grpc.Server, service teacherService) {
	client.RegisterTeacherServer(gRPC, &teacherServerAPI{
		service: service,
	})
}
