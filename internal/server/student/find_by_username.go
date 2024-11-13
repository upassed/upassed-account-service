package student

import (
	"context"
	"github.com/upassed/upassed-account-service/internal/handling"
	requestid "github.com/upassed/upassed-account-service/internal/middleware/common/request_id"
	"github.com/upassed/upassed-account-service/internal/tracing"
	"github.com/upassed/upassed-account-service/pkg/client"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
)

func (server *studentServerAPI) FindByUsername(ctx context.Context, request *client.StudentFindByUsernameRequest) (*client.StudentFindByUsernameResponse, error) {
	spanContext, span := otel.Tracer(server.cfg.Tracing.StudentTracerName).Start(ctx, "student#FindByUsername")
	span.SetAttributes(
		attribute.String(string(requestid.ContextKey), requestid.GetRequestIDFromContext(ctx)),
		attribute.String("studentUsername", request.GetStudentUsername()),
	)
	defer span.End()

	if err := request.Validate(); err != nil {
		tracing.SetSpanError(span, err)
		return nil, handling.Wrap(err, handling.WithCode(codes.InvalidArgument))
	}

	response, err := server.service.FindByUsername(spanContext, request.GetStudentUsername())
	if err != nil {
		tracing.SetSpanError(span, err)
		return nil, err
	}

	return ConvertToFindByUsernameResponse(response), nil
}
