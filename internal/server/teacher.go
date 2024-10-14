package server

import (
	"context"

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
	response, err := server.service.Create(ctx, convertedRequest)
	if err != nil {
		return nil, err
	}

	return converter.ConvertTeacherCreateResponse(response), nil
}

func (server *teacherServerAPI) FindByID(ctx context.Context, request *client.TeacherFindByIDRequest) (*client.TeacherFindByIDResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, handling.WrapAsApplicationError(err, handling.WithCode(codes.InvalidArgument))
	}

	teacher, err := server.service.FindByID(ctx, request.GetTeacherId())
	if err != nil {
		return nil, err
	}

	return converter.ConvertTeacher(teacher), nil
}
