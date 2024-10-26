package group

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/upassed/upassed-account-service/internal/handling"
	"github.com/upassed/upassed-account-service/internal/logging"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"log/slog"
)

var (
	errSearchingStudentsInGroup = errors.New("error while searching students in group")
)

func (repository *groupRepositoryImpl) FindStudentsInGroup(ctx context.Context, groupID uuid.UUID) ([]*domain.Student, error) {
	log := logging.Wrap(repository.log,
		logging.WithOp(repository.Exists),
		logging.WithCtx(ctx),
		logging.WithAny("groupID", groupID),
	)

	_, span := otel.Tracer(repository.cfg.Tracing.GroupTracerName).Start(ctx, "groupRepository#FindStudentsInGroup")
	defer span.End()

	log.Info("started searching students in group in a database")
	foundStudents := make([]*domain.Student, 0)
	searchResult := repository.db.WithContext(ctx).Preload("Group").Where("group_id = ?", groupID).Find(&foundStudents)
	if searchResult.Error != nil {
		log.Error("error while searching students in group in the database", logging.Error(searchResult.Error))
		return nil, handling.New(errSearchingStudentsInGroup.Error(), codes.Internal)
	}

	log.Info("students in group were successfully found in a database", slog.Int("studentInGroup", len(foundStudents)))
	return foundStudents, nil
}
