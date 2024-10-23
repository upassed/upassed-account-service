package caching

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
	ErrGroupIsNotPresentInCache     = errors.New("group id is not present as a key in the cache")
	errFetchingGroupFromCache       = errors.New("unable to get group by id from the cache")
	errUnmarshallingGroupDataToJson = errors.New("unable to unmarshall group data from the cache to json format")
)

func (client *RedisClient) GetGroupByID(ctx context.Context, groupID uuid.UUID) (domain.Group, error) {
	op := runtime.FuncForPC(reflect.ValueOf(client.GetGroupByID).Pointer()).Name()

	log := client.log.With(
		slog.String("op", op),
		slog.Any("groupID", groupID),
	)

	_, span := otel.Tracer(client.cfg.Tracing.GroupTracerName).Start(ctx, "redisClient#GetGroupByID")
	defer span.End()

	groupData, err := client.client.Get(ctx, "group:"+groupID.String()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Error("group by id was not found in cache")
			return domain.Group{}, ErrGroupIsNotPresentInCache
		}

		log.Error("error while fetching group by id from cache", logging.Error(err))
		return domain.Group{}, errFetchingGroupFromCache
	}

	var group domain.Group
	if err := json.Unmarshal([]byte(groupData), &group); err != nil {
		log.Error("error while unmarshalling group data to json", logging.Error(err))
		return domain.Group{}, errUnmarshallingGroupDataToJson
	}

	log.Info("group was successfully found by id")
	return group, nil
}
