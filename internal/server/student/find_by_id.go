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
	spanContext, span := otel.Tracer(server.cfg.Tracing.StudentTracerName).Start(ctx, "student#FindByID")
	span.SetAttributes(
		attribute.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
		attribute.String("id", request.GetStudentId()),
	)
	defer span.End()

	if err := request.Validate(); err != nil {
		span.SetAttributes(attribute.String("err", err.Error()))
		return nil, handling.Wrap(err, handling.WithCode(codes.InvalidArgument))
	}

	response, err := server.service.FindByID(spanContext, uuid.MustParse(request.GetStudentId()))
	if err != nil {
		span.SetAttributes(attribute.String("err", err.Error()))
		return nil, err
	}

	return ConvertToFindByIDResponse(response), nil
}
