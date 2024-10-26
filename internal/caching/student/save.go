package student

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/upassed/upassed-account-service/internal/logging"
	domain "github.com/upassed/upassed-account-service/internal/repository/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var (
	errMarshallingStudentData   = errors.New("unable to marshall student data to json format")
	errSavingStudentDataToCache = errors.New("unable to save student data to redis cache")
)

func (client *RedisClient) Save(ctx context.Context, student *domain.Student) error {
	_, span := otel.Tracer(client.cfg.Tracing.StudentTracerName).Start(ctx, "redisClient#Save")
	span.SetAttributes(attribute.String("username", student.Username))
	defer span.End()

	log := logging.Wrap(client.log,
		logging.WithOp(client.Save),
		logging.WithCtx(ctx),
		logging.WithAny("username", student.Username),
	)

	log.Info("marshalling student data to json format to save it to the cache")
	jsonStudentData, err := json.Marshal(student)
	if err != nil {
		log.Error("unable to marshall student data to json format")
		span.SetAttributes(attribute.String("err", err.Error()))
		return errMarshallingStudentData
	}

	if err := client.client.Set(ctx, fmt.Sprintf(keyFormat, student.ID.String()), jsonStudentData, client.cfg.GetRedisEntityTTL()).Err(); err != nil {
		log.Error("unable to save student data to the cache", logging.Error(err))
		span.SetAttributes(attribute.String("err", err.Error()))
		return errSavingStudentDataToCache
	}

	log.Info("student successfully saved to the cache")
	return nil
}
