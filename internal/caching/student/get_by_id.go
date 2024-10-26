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
	"go.opentelemetry.io/otel"
)

var (
	ErrStudentIsNotPresentInCache     = errors.New("student id is not present as a key in the cache")
	errFetchingStudentFromCache       = errors.New("unable to get student by id from the cache")
	errUnmarshallingStudentDataToJson = errors.New("unable to unmarshall student data from the cache to json format")
)

func (client *RedisClient) GetByID(ctx context.Context, studentID uuid.UUID) (*domain.Student, error) {
	_, span := otel.Tracer(client.cfg.Tracing.StudentTracerName).Start(ctx, "redisClient#GetByID")
	defer span.End()

	log := logging.Wrap(client.log,
		logging.WithOp(client.GetByID),
		logging.WithCtx(ctx),
		logging.WithAny("studentID", studentID),
	)

	studentData, err := client.client.Get(ctx, fmt.Sprintf(keyFormat, studentID.String())).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Error("student by id was not found in cache")
			return nil, ErrStudentIsNotPresentInCache
		}

		log.Error("error while fetching student by id from cache", logging.Error(err))
		return nil, errFetchingStudentFromCache
	}

	var student domain.Student
	if err := json.Unmarshal([]byte(studentData), &student); err != nil {
		log.Error("error while unmarshalling student data to json", logging.Error(err))
		return nil, errUnmarshallingStudentDataToJson
	}

	log.Info("student was successfully found by id")
	return &student, nil
}
