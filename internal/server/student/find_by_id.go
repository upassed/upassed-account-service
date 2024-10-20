package student

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

func (server *studentServerAPI) FindByID(ctx context.Context, request *client.StudentFindByIDRequest) (*client.StudentFindByIDResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, handling.Wrap(err, handling.WithCode(codes.InvalidArgument))
	}

	spanContext, span := otel.Tracer(server.cfg.Tracing.StudentTracerName).Start(ctx, "student#FindByID")
	span.SetAttributes(attribute.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)))
	defer span.End()

	response, err := server.service.FindByID(spanContext, uuid.MustParse(request.GetStudentId()))
	if err != nil {
		return nil, err
	}

	return ConvertToFindByIDResponse(response), nil
}
