package teacher

import (
	"context"
	"encoding/json"
	"errors"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"go.opentelemetry.io/otel"
	"log/slog"
	"reflect"
	"runtime"
)

var (
	errMarshallingTeacherData   = errors.New("unable to marshall teacher data to json format")
	errSavingTeacherDataToCache = errors.New("unable to save teacher data to redis cache")
)

func (client *RedisClient) SaveTeacher(ctx context.Context, teacher domain.Teacher) error {
	op := runtime.FuncForPC(reflect.ValueOf(client.SaveTeacher).Pointer()).Name()

	log := client.log.With(
		slog.String("op", op),
		slog.Any("teacherID", teacher.ID),
	)

	_, span := otel.Tracer(client.cfg.Tracing.TeacherTracerName).Start(ctx, "redisClient#SaveTeacher")
	defer span.End()

	jsonTeacherData, err := json.Marshal(teacher)
	if err != nil {
		log.Error("unable to marshall teacher data to json format")
		return errMarshallingTeacherData
	}

	if err := client.client.Set(ctx, "teacher:"+teacher.ID.String(), jsonTeacherData, client.cfg.GetRedisEntityTTL()).Err(); err != nil {
		return errSavingTeacherDataToCache
	}

	log.Info("teacher successfully saved to the cache")
	return nil
}
