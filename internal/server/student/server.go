package student

import (
	"context"

	"github.com/google/uuid"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"github.com/upassed/upassed-account-service/pkg/client"
	"google.golang.org/grpc"
)

type studentServerAPI struct {
	client.UnimplementedStudentServer
	service studentService
}

type studentService interface {
	Create(context.Context, business.Student) (business.StudentCreateResponse, error)
	FindByID(ctx context.Context, studentID uuid.UUID) (business.Student, error)
}

func Register(gRPC *grpc.Server, service studentService) {
	client.RegisterStudentServer(gRPC, &studentServerAPI{
		service: service,
	})
}
