package group

import (
	"context"
	"github.com/upassed/upassed-account-service/internal/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/pkg/client"
	"google.golang.org/grpc/codes"
)

func (server *groupServerAPI) FindStudentsInGroup(ctx context.Context, request *client.FindStudentsInGroupRequest) (*client.FindStudentsInGroupResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, handling.Wrap(err, handling.WithCode(codes.InvalidArgument))
	}

	spanContext, span := otel.Tracer(server.cfg.Tracing.GroupTracerName).Start(ctx, "group#FindStudentsInGroup")
	span.SetAttributes(attribute.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)))
	defer span.End()

	response, err := server.service.FindStudentsInGroup(spanContext, uuid.MustParse(request.GetGroupId()))
	if err != nil {
		return nil, err
	}

	return ConvertToFindStudentsInGroupResponse(response), nil
}
