package student

import (
	"context"

	service "github.com/upassed/upassed-account-service/internal/service/student"
	"github.com/upassed/upassed-account-service/pkg/client"
	"google.golang.org/grpc"
)

type studentServerAPI struct {
	client.UnimplementedStudentServer
	service studentService
}

type studentService interface {
	Create(context.Context, service.Student) (service.StudentCreateResponse, error)
	FindByID(ctx context.Context, studentID string) (service.Student, error)
}

func Register(gRPC *grpc.Server, service studentService) {
	client.RegisterStudentServer(gRPC, &studentServerAPI{
		service: service,
	})
}
