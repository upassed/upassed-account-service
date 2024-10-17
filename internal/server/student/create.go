package student

import (
	"context"

	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/pkg/client"
	"google.golang.org/grpc/codes"
)

func (server *studentServerAPI) Create(ctx context.Context, request *client.StudentCreateRequest) (*client.StudentCreateResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, handling.Wrap(err, handling.WithCode(codes.InvalidArgument))
	}

	student := ConvertToStudent(request)
	response, err := server.service.Create(ctx, student)
	if err != nil {
		return nil, err
	}

	return ConvertToStudentCreateResponse(response), nil
}
