package group

import (
	"context"

	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/pkg/client"
	"google.golang.org/grpc/codes"
)

func (server *groupServerAPI) SearchByFilter(ctx context.Context, request *client.GroupSearchByFilterRequest) (*client.GroupSearchByFilterResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, handling.Wrap(err, handling.WithCode(codes.InvalidArgument))
	}

	matchedGroups, err := server.service.FindByFilter(ctx, ConvertToGroupFilter(request))
	if err != nil {
		return nil, err
	}

	return ConvertToSearchByFilterResponse(matchedGroups), nil
}
