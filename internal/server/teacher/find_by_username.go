package teacher

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

func (server *teacherServerAPI) FindByUsername(ctx context.Context, request *client.TeacherFindByUsernameRequest) (*client.TeacherFindByUsernameResponse, error) {
	spanContext, span := otel.Tracer(server.cfg.Tracing.TeacherTracerName).Start(ctx, "teacher#FindByUsername")
	span.SetAttributes(
		attribute.String(string(requestid.ContextKey), requestid.GetRequestIDFromContext(ctx)),
		attribute.String("teacherUsername", request.GetTeacherUsername()),
	)
	defer span.End()

	if err := request.Validate(); err != nil {
		tracing.SetSpanError(span, err)
		return nil, handling.Wrap(err, handling.WithCode(codes.InvalidArgument))
	}

	teacher, err := server.service.FindByUsername(spanContext, request.GetTeacherUsername())
	if err != nil {
		tracing.SetSpanError(span, err)
		return nil, err
	}

	return ConvertToFindByUsernameResponse(teacher), nil
}
