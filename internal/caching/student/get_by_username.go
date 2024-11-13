package student

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/upassed/upassed-account-service/internal/logging"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"github.com/upassed/upassed-account-service/internal/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var (
	ErrStudentIsNotPresentInCache     = errors.New("student username is not present as a key in the cache")
	errFetchingStudentFromCache       = errors.New("unable to get student by username from the cache")
	errUnmarshallingStudentDataToJson = errors.New("unable to unmarshall student data from the cache to json format")
)

func (client *RedisClient) GetByUsername(ctx context.Context, studentUsername string) (*domain.Student, error) {
	_, span := otel.Tracer(client.cfg.Tracing.StudentTracerName).Start(ctx, "redisClient#GetByUsername")
	span.SetAttributes(attribute.String("studentUsername", studentUsername))
	defer span.End()

	log := logging.Wrap(client.log,
		logging.WithOp(client.GetByUsername),
		logging.WithCtx(ctx),
		logging.WithAny("studentUsername", studentUsername),
	)

	log.Info("getting student data from the cache by username")
	studentData, err := client.client.Get(ctx, fmt.Sprintf(usernameKeyFormat, studentUsername)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Error("student by username was not found in cache")
			tracing.SetSpanError(span, err)
			return nil, ErrStudentIsNotPresentInCache
		}

		log.Error("error while fetching student by username from cache", logging.Error(err))
		tracing.SetSpanError(span, err)
		return nil, errFetchingStudentFromCache
	}

	log.Info("student data found in cache by username, unmarshalling data from json")
	var student domain.Student
	if err := json.Unmarshal([]byte(studentData), &student); err != nil {
		log.Error("error while unmarshalling student data to json", logging.Error(err))
		tracing.SetSpanError(span, err)
		return nil, errUnmarshallingStudentDataToJson
	}

	log.Info("student was successfully found in the cache and unmarshalled")
	return &student, nil
}
