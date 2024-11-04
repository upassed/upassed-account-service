package student

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/upassed/upassed-account-service/internal/logging"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"github.com/upassed/upassed-account-service/internal/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var (
	ErrStudentIsNotPresentInCache     = errors.New("student id is not present as a key in the cache")
	errFetchingStudentFromCache       = errors.New("unable to get student by id from the cache")
	errUnmarshallingStudentDataToJson = errors.New("unable to unmarshall student data from the cache to json format")
)

func (client *RedisClient) GetByID(ctx context.Context, studentID uuid.UUID) (*domain.Student, error) {
	_, span := otel.Tracer(client.cfg.Tracing.StudentTracerName).Start(ctx, "redisClient#GetByID")
	span.SetAttributes(attribute.String("id", studentID.String()))
	defer span.End()

	log := logging.Wrap(client.log,
		logging.WithOp(client.GetByID),
		logging.WithCtx(ctx),
		logging.WithAny("studentID", studentID),
	)

	log.Info("getting student data from the cache")
	studentData, err := client.client.Get(ctx, fmt.Sprintf(keyFormat, studentID.String())).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Error("student by id was not found in cache")
			tracing.SetSpanError(span, err)
			return nil, ErrStudentIsNotPresentInCache
		}

		log.Error("error while fetching student by id from cache", logging.Error(err))
		tracing.SetSpanError(span, err)
		return nil, errFetchingStudentFromCache
	}

	log.Info("student data found in cache, unmarshalling data from json")
	var student domain.Student
	if err := json.Unmarshal([]byte(studentData), &student); err != nil {
		log.Error("error while unmarshalling student data to json", logging.Error(err))
		tracing.SetSpanError(span, err)
		return nil, errUnmarshallingStudentDataToJson
	}

	log.Info("student was successfully found in the cache and unmarshalled")
	return &student, nil
}
