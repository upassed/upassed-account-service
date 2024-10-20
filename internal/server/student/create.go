package student

import (
	"context"
	"github.com/upassed/upassed-account-service/internal/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/pkg/client"
	"google.golang.org/grpc/codes"
)

func (server *studentServerAPI) Create(ctx context.Context, request *client.StudentCreateRequest) (*client.StudentCreateResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, handling.Wrap(err, handling.WithCode(codes.InvalidArgument))
	}

	spanContext, span := otel.Tracer(server.cfg.Tracing.StudentTracerName).Start(ctx, "student#Create")
	span.SetAttributes(attribute.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)))
	defer span.End()

	student := ConvertToStudent(request)
	response, err := server.service.Create(spanContext, student)
	if err != nil {
		return nil, err
	}

	return ConvertToStudentCreateResponse(response), nil
}
