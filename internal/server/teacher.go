package server

import (
	"context"
	"time"

	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/server/converter"
	business "github.com/upassed/upassed-account-service/internal/service/model"
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
	Create(context.Context, business.TeacherCreateRequest) (business.TeacherCreateResponse, error)
}

func (server *teacherServerAPI) Create(ctx context.Context, request *client.TeacherCreateRequest) (*client.TeacherCreateResponse, error) {
	convertedRequest := converter.ConvertTeacherCreateRequest(request)

	contextWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	response, err := server.service.Create(contextWithTimeout, convertedRequest)
	if err != nil {
		return nil, handling.HandleServiceLayerError(err)
	}

	convertedResponse := converter.TestConvertTeacherCreateResponse(response)
	return &convertedResponse, nil
}
