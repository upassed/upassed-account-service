package group

import (
	"context"

	"github.com/google/uuid"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"github.com/upassed/upassed-account-service/pkg/client"
	"google.golang.org/grpc"
)

type groupServerAPI struct {
	client.UnimplementedGroupServer
	service groupService
}

type groupService interface {
	FindStudentsInGroup(context.Context, uuid.UUID) ([]business.Student, error)
	FindByID(context.Context, uuid.UUID) (business.Group, error)
}

func Register(gRPC *grpc.Server, service groupService) {
	client.RegisterGroupServer(gRPC, &groupServerAPI{
		service: service,
	})
}
