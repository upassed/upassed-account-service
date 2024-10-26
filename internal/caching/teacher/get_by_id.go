package teacher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/upassed/upassed-account-service/internal/logging"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var (
	ErrTeacherIsNotPresentInCache     = errors.New("teacher id is not present as a key in the cache")
	errFetchingTeacherFromCache       = errors.New("unable to get teacher by id from the cache")
	errUnmarshallingTeacherDataToJson = errors.New(`unable to unmarshall teacher data from the cache to json format`)
)

func (client *RedisClient) GetByID(ctx context.Context, teacherID uuid.UUID) (*domain.Teacher, error) {
	_, span := otel.Tracer(client.cfg.Tracing.TeacherTracerName).Start(ctx, "redisClient#GetTeacherID")
	span.SetAttributes(attribute.String("id", teacherID.String()))
	defer span.End()

	log := logging.Wrap(client.log,
		logging.WithOp(client.GetByID),
		logging.WithCtx(ctx),
		logging.WithAny("teacherID", teacherID),
	)

	log.Info("getting teacher data from the cache")
	teacherData, err := client.client.Get(ctx, fmt.Sprintf(keyFormat, teacherID.String())).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Error("teacher by id was not found in cache")
			return nil, ErrTeacherIsNotPresentInCache
		}

		log.Error("error while fetching teacher by id from cache", logging.Error(err))
		return nil, errFetchingTeacherFromCache
	}

	log.Info("teacher data found in the cache, unmarshalling from json")
	var teacher domain.Teacher
	if err := json.Unmarshal([]byte(teacherData), &teacher); err != nil {
		log.Error("error while unmarshalling teacher data to json", logging.Error(err))
		return nil, errUnmarshallingTeacherDataToJson
	}

	log.Info("teacher was successfully found in the cache and unmarshalled")
	return &teacher, nil
}
