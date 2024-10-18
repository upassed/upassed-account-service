package group

import (
	"context"

	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/pkg/client"
	"google.golang.org/grpc/codes"
)

func (server *groupServerAPI) FindByID(ctx context.Context, request *client.GroupFindByIDRequest) (*client.GroupFindByIDResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, handling.Wrap(err, handling.WithCode(codes.InvalidArgument))
	}

	foundGroup, err := server.service.FindByID(ctx, uuid.MustParse(request.GetGroupId()))
	if err != nil {
		return nil, err
	}

	return ConvertToFindByIDResponse(foundGroup), nil
}
