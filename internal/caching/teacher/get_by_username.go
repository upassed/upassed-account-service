package teacher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/upassed/upassed-account-service/internal/logging"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var (
	ErrTeacherUsernameIsNotPresentInCache = errors.New("teacher username is not present as a key in the cache")
	errFetchingTeacherByUsernameFromCache = errors.New("unable to get teacher by username from the cache")
	errUnmarshallingTeacherDataToJson     = errors.New(`unable to unmarshall teacher data from the cache to json format`)
)

func (client *RedisClient) GetByUsername(ctx context.Context, teacherUsername string) (*domain.Teacher, error) {
	_, span := otel.Tracer(client.cfg.Tracing.TeacherTracerName).Start(ctx, "redisClient#GetByUsername")
	span.SetAttributes(attribute.String("teacherUsername", teacherUsername))
	defer span.End()

	log := logging.Wrap(client.log,
		logging.WithOp(client.GetByUsername),
		logging.WithCtx(ctx),
		logging.WithAny("teacherUsername", teacherUsername),
	)

	log.Info("getting teacher data from the cache by username")
	teacherData, err := client.client.Get(ctx, fmt.Sprintf(usernameKeyFormat, teacherUsername)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Error("teacher by username was not found in cache")
			return nil, ErrTeacherUsernameIsNotPresentInCache
		}

		log.Error("error while fetching teacher by username from cache", logging.Error(err))
		return nil, errFetchingTeacherByUsernameFromCache
	}

	log.Info("teacher data by username found in the cache, unmarshalling from json")
	var teacher domain.Teacher
	if err := json.Unmarshal([]byte(teacherData), &teacher); err != nil {
		log.Error("error while unmarshalling teacher data to json", logging.Error(err))
		return nil, errUnmarshallingTeacherDataToJson
	}

	log.Info("teacher by username was successfully found in the cache and unmarshalled")
	return &teacher, nil
}
