package teacher

import (
	"context"
	"github.com/upassed/upassed-account-service/internal/config"

	business "github.com/upassed/upassed-account-service/internal/service/model"
	"github.com/upassed/upassed-account-service/pkg/client"
	"google.golang.org/grpc"
)

type teacherServerAPI struct {
	client.UnimplementedTeacherServer
	cfg     *config.Config
	service teacherService
}

type teacherService interface {
	Create(context.Context, *business.Teacher) (*business.TeacherCreateResponse, error)
	FindByUsername(ctx context.Context, teacherUsername string) (*business.Teacher, error)
}

func Register(gRPC *grpc.Server, cfg *config.Config, service teacherService) {
	client.RegisterTeacherServer(gRPC, &teacherServerAPI{
		cfg:     cfg,
		service: service,
	})
}
