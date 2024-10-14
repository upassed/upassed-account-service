package teacher

import (
	"context"

	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/pkg/client"
	"google.golang.org/grpc/codes"
)

func (server *teacherServerAPI) Create(ctx context.Context, request *client.TeacherCreateRequest) (*client.TeacherCreateResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, handling.WrapAsApplicationError(err, handling.WithCode(codes.InvalidArgument))
	}

	teacher := ConvertToTeacher(request)
	response, err := server.service.Create(ctx, teacher)
	if err != nil {
		return nil, err
	}

	return ConvertToTeacherCreateResponse(response), nil
}
