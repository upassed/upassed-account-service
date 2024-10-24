package teacher

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/upassed/upassed-account-service/internal/logging"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"go.opentelemetry.io/otel"
	"log/slog"
	"reflect"
	"runtime"
)

var (
	ErrTeacherIsNotPresentInCache     = errors.New("teacher id is not present as a key in the cache")
	errFetchingTeacherFromCache       = errors.New("unable to get teacher by id from the cache")
	errUnmarshallingTeacherDataToJson = errors.New(`unable to unmarshall teacher data from the cache to json format`)
)

func (client *RedisClient) GetTeacherByID(ctx context.Context, teacherID uuid.UUID) (domain.Teacher, error) {
	op := runtime.FuncForPC(reflect.ValueOf(client.GetTeacherByID).Pointer()).Name()

	log := client.log.With(
		slog.String("op", op),
		slog.Any("teacherID", teacherID),
	)

	_, span := otel.Tracer(client.cfg.Tracing.TeacherTracerName).Start(ctx, "redisClient#GetTeacherID")
	defer span.End()

	teacherData, err := client.client.Get(ctx, "teacher:"+teacherID.String()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Error("teacher by id was not found in cache")
			return domain.Teacher{}, ErrTeacherIsNotPresentInCache
		}

		log.Error("error while fetching teacher by id from cache", logging.Error(err))
		return domain.Teacher{}, errFetchingTeacherFromCache
	}

	var teacher domain.Teacher
	if err := json.Unmarshal([]byte(teacherData), &teacher); err != nil {
		log.Error("error while unmarshalling teacher data to json", logging.Error(err))
		return domain.Teacher{}, errUnmarshallingTeacherDataToJson
	}

	log.Info("teacher was successfully found by id")
	return teacher, nil
}
