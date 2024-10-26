package student

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/middleware"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"go.opentelemetry.io/otel"
	"log/slog"
	"reflect"
	"runtime"
)

var (
	ErrStudentIsNotPresentInCache     = errors.New("student id is not present as a key in the cache")
	errFetchingStudentFromCache       = errors.New("unable to get student by id from the cache")
	errUnmarshallingStudentDataToJson = errors.New("unable to unmarshall student data from the cache to json format")
)

func (client *RedisClient) GetStudentByID(ctx context.Context, studentID uuid.UUID) (*domain.Student, error) {
	op := runtime.FuncForPC(reflect.ValueOf(client.GetStudentByID).Pointer()).Name()

	log := client.log.With(
		slog.String("op", op),
		slog.Any("studentID", studentID),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	_, span := otel.Tracer(client.cfg.Tracing.StudentTracerName).Start(ctx, "redisClient#GetStudentByID")
	defer span.End()

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
