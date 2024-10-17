package group

import (
	"context"

	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/pkg/client"
	"google.golang.org/grpc/codes"
)

func (server *groupServerAPI) FindStudentsInGroup(ctx context.Context, request *client.FindStudentsInGroupRequest) (*client.FindStudentsInGroupResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, handling.Wrap(err, handling.WithCode(codes.InvalidArgument))
	}

	response, err := server.service.FindStudentsInGroup(ctx, uuid.MustParse(request.GetGroupId()))
	if err != nil {
		return nil, err
	}

	return ConvertToFindStudentsInGroupResponse(response), nil
}
