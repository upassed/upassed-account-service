package teacher

import (
	"context"

	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/pkg/client"
	"google.golang.org/grpc/codes"
)

func (server *teacherServerAPI) FindByID(ctx context.Context, request *client.TeacherFindByIDRequest) (*client.TeacherFindByIDResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, handling.Wrap(err, handling.WithCode(codes.InvalidArgument))
	}

	teacher, err := server.service.FindByID(ctx, uuid.MustParse(request.GetTeacherId()))
	if err != nil {
		return nil, err
	}

	return ConvertToFindByIDResponse(teacher), nil
}
