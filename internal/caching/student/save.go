package student

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/upassed/upassed-account-service/internal/middleware"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"go.opentelemetry.io/otel"
	"log/slog"
	"reflect"
	"runtime"
)

var (
	errMarshallingStudentData   = errors.New("unable to marshall student data to json format")
	errSavingStudentDataToCache = errors.New("unable to save student data to redis cache")
)

func (client *RedisClient) SaveStudent(ctx context.Context, student *domain.Student) error {
	op := runtime.FuncForPC(reflect.ValueOf(client.SaveStudent).Pointer()).Name()

	log := client.log.With(
		slog.String("op", op),
		slog.Any("studentID", student.ID),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	_, span := otel.Tracer(client.cfg.Tracing.StudentTracerName).Start(ctx, "redisClient#SaveStudent")
	defer span.End()

	jsonStudentData, err := json.Marshal(student)
	if err != nil {
		log.Error("unable to marshall student data to json format")
		return errMarshallingStudentData
	}

	if err := client.client.Set(ctx, fmt.Sprintf(keyFormat, student.ID.String()), jsonStudentData, client.cfg.GetRedisEntityTTL()).Err(); err != nil {
		return errSavingStudentDataToCache
	}

	log.Info("student successfully saved to the cache")
	return nil
}
