package group

import (
	"context"

	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/service/student"
	"github.com/upassed/upassed-account-service/pkg/client"
	"google.golang.org/grpc"
)

type groupServerAPI struct {
	client.UnimplementedGroupServer
	service groupService
}

type groupService interface {
	FindStudentsInGroup(context.Context, uuid.UUID) ([]student.Student, error)
}

func Register(gRPC *grpc.Server, service groupService) {
	client.RegisterGroupServer(gRPC, &groupServerAPI{
		service: service,
	})
}
