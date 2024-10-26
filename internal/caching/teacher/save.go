package teacher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/upassed/upassed-account-service/internal/logging"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"go.opentelemetry.io/otel"
)

var (
	errMarshallingTeacherData   = errors.New("unable to marshall teacher data to json format")
	errSavingTeacherDataToCache = errors.New("unable to save teacher data to redis cache")
)

func (client *RedisClient) Save(ctx context.Context, teacher *domain.Teacher) error {
	log := logging.Wrap(client.log,
		logging.WithOp(client.GetByID),
		logging.WithCtx(ctx),
		logging.WithAny("teacherID", teacher.ID),
	)

	_, span := otel.Tracer(client.cfg.Tracing.TeacherTracerName).Start(ctx, "redisClient#Save")
	defer span.End()

	jsonTeacherData, err := json.Marshal(teacher)
	if err != nil {
		log.Error("unable to marshall teacher data to json format")
		return errMarshallingTeacherData
	}

	if err := client.client.Set(ctx, fmt.Sprintf(keyFormat, teacher.ID.String()), jsonTeacherData, client.cfg.GetRedisEntityTTL()).Err(); err != nil {
		return errSavingTeacherDataToCache
	}

	log.Info("teacher successfully saved to the cache")
	return nil
}
