package group

import (
	"context"
	"github.com/upassed/upassed-account-service/internal/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/pkg/client"
	"google.golang.org/grpc/codes"
)

func (server *groupServerAPI) SearchByFilter(ctx context.Context, request *client.GroupSearchByFilterRequest) (*client.GroupSearchByFilterResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, handling.Wrap(err, handling.WithCode(codes.InvalidArgument))
	}

	spanContext, span := otel.Tracer(server.cfg.Tracing.GroupTracerName).Start(ctx, "group#SearchByFilter")
	span.SetAttributes(attribute.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)))
	defer span.End()

	matchedGroups, err := server.service.FindByFilter(spanContext, ConvertToGroupFilter(request))
	if err != nil {
		return nil, err
	}

	return ConvertToSearchByFilterResponse(matchedGroups), nil
}
