package student

import (
	"context"

	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/pkg/client"
	"google.golang.org/grpc/codes"
)

func (server *studentServerAPI) FindByID(ctx context.Context, request *client.StudentFindByIDRequest) (*client.StudentFindByIDResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, handling.Wrap(err, handling.WithCode(codes.InvalidArgument))
	}

	response, err := server.service.FindByID(ctx, uuid.MustParse(request.GetStudentId()))
	if err != nil {
		return nil, err
	}

	return ConvertToFindByIDResponse(response), nil
}
