package server

import (
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
}
