package group

import (
	"context"
	"github.com/upassed/upassed-account-service/internal/config"

	"github.com/google/uuid"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"github.com/upassed/upassed-account-service/pkg/client"
	"google.golang.org/grpc"
)

type groupServerAPI struct {
	client.UnimplementedGroupServer
	cfg     *config.Config
	service groupService
}

type groupService interface {
	FindByID(context.Context, uuid.UUID) (*business.Group, error)
	FindStudentsInGroup(context.Context, uuid.UUID) ([]*business.Student, error)
	FindByFilter(context.Context, *business.GroupFilter) ([]*business.Group, error)
}

func Register(gRPC *grpc.Server, cfg *config.Config, service groupService) {
	client.RegisterGroupServer(gRPC, &groupServerAPI{
		cfg:     cfg,
		service: service,
	})
}
