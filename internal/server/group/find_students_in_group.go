package group

import (
	"context"

	"github.com/upassed/upassed-account-service/pkg/client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *groupServerAPI) FindStudentsInGroup(ctx context.Context, request *client.FindStudentsInGroupRequest) (*client.FindStudentsInGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindStudentsInGroup not implemented")
}
