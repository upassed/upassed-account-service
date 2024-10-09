package server

import (
	"context"
	"time"

	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/server/converter"
	business "github.com/upassed/upassed-account-service/internal/service/model"
	"github.com/upassed/upassed-account-service/pkg/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
	Create(context.Context, business.Teacher) (business.TeacherCreateResponse, error)
	FindByID(ctx context.Context, teacherID string) (business.Teacher, error)
}

func (server *teacherServerAPI) Create(ctx context.Context, request *client.TeacherCreateRequest) (*client.TeacherCreateResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, handling.WrapAsApplicationError(err, handling.WithCode(codes.InvalidArgument))
	}

	convertedRequest := converter.ConvertTeacherCreateRequest(request)
	contextWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	response, err := server.service.Create(contextWithTimeout, convertedRequest)
	if err != nil {
		return nil, handling.HandleApplicationError(err)
	}

	convertedResponse := converter.ConvertTeacherCreateResponse(response)
	return &convertedResponse, nil
}

func (server *teacherServerAPI) FindByID(ctx context.Context, request *client.TeacherFindByIDRequest) (*client.TeacherFindByIDResponse, error) {
	contextWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	teacher, err := server.service.FindByID(contextWithTimeout, request.GetTeacherId())
	if err != nil {
		return nil, handling.HandleApplicationError(err)
	}

	convertedTeacherResponse := converter.ConvertTeacher(teacher)
	return &convertedTeacherResponse, nil
}
