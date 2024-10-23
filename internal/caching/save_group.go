package caching

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
	errMarshallingGroupData   = errors.New("unable to marshall group data to json format")
	errSavingGroupDataToCache = errors.New("unable to save group data to redis cache")
)

func (client *RedisClient) SaveGroup(ctx context.Context, group domain.Group) error {
	op := runtime.FuncForPC(reflect.ValueOf(client.SaveGroup).Pointer()).Name()

	log := client.log.With(
		slog.String("op", op),
		slog.Any("groupID", group.ID),
	)

	_, span := otel.Tracer(client.cfg.Tracing.GroupTracerName).Start(ctx, "redisClient#SaveGroup")
	defer span.End()

	jsonGroupData, err := json.Marshal(group)
	if err != nil {
		log.Error("unable to marshall group data to json format")
		return errMarshallingGroupData
	}

	if err := client.client.Set(ctx, "group:"+group.ID.String(), jsonGroupData, client.cfg.GetRedisEntityTTL()).Err(); err != nil {
		return errSavingGroupDataToCache
	}

	log.Info("group successfully saved to the cache")
	return nil
}
