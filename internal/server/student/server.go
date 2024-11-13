package student

import (
	"context"
	"github.com/upassed/upassed-account-service/internal/config"

	business "github.com/upassed/upassed-account-service/internal/service/model"
	"github.com/upassed/upassed-account-service/pkg/client"
	"google.golang.org/grpc"
)

type studentServerAPI struct {
	client.UnimplementedStudentServer
	cfg     *config.Config
	service studentService
}

type studentService interface {
	Create(context.Context, *business.Student) (*business.StudentCreateResponse, error)
	FindByUsername(ctx context.Context, studentUsername string) (*business.Student, error)
}

func Register(gRPC *grpc.Server, cfg *config.Config, service studentService) {
	client.RegisterStudentServer(gRPC, &studentServerAPI{
		cfg:     cfg,
		service: service,
	})
}
