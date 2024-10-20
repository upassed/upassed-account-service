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

func (server *groupServerAPI) FindByID(ctx context.Context, request *client.GroupFindByIDRequest) (*client.GroupFindByIDResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, handling.Wrap(err, handling.WithCode(codes.InvalidArgument))
	}

	spanContext, span := otel.Tracer(server.cfg.Tracing.GroupTracerName).Start(ctx, "group#FindByID")
	span.SetAttributes(attribute.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)))
	defer span.End()

	foundGroup, err := server.service.FindByID(spanContext, uuid.MustParse(request.GetGroupId()))
	if err != nil {
		return nil, err
	}

	return ConvertToFindByIDResponse(foundGroup), nil
}
