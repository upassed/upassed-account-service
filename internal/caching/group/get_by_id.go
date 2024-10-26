package group

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
	ErrGroupIsNotPresentInCache     = errors.New("group id is not present as a key in the cache")
	errFetchingGroupFromCache       = errors.New("unable to get group by id from the cache")
	errUnmarshallingGroupDataToJson = errors.New("unable to unmarshall group data from the cache to json format")
)

func (client *RedisClient) GetByID(ctx context.Context, groupID uuid.UUID) (*domain.Group, error) {
	_, span := otel.Tracer(client.cfg.Tracing.GroupTracerName).Start(ctx, "redisClient#GetByID")
	span.SetAttributes(attribute.String("id", groupID.String()))
	defer span.End()

	log := logging.Wrap(client.log,
		logging.WithOp(client.GetByID),
		logging.WithCtx(ctx),
		logging.WithAny("groupID", groupID),
	)

	log.Info("started getting group data by id from cache")
	groupData, err := client.client.Get(ctx, fmt.Sprintf(keyFormat, groupID.String())).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Error("group by id was not found in cache")
			span.SetAttributes(attribute.String("err", err.Error()))
			return nil, ErrGroupIsNotPresentInCache
		}

		log.Error("error while fetching group by id from cache", logging.Error(err))
		span.SetAttributes(attribute.String("err", err.Error()))
		return nil, errFetchingGroupFromCache
	}

	log.Info("group by id was found in cache, unmarshalling from json")
	var group domain.Group
	if err := json.Unmarshal([]byte(groupData), &group); err != nil {
		log.Error("error while unmarshalling group data to json", logging.Error(err))
		span.SetAttributes(attribute.String("err", err.Error()))
		return nil, errUnmarshallingGroupDataToJson
	}

	log.Info("group was successfully found and unmarshalled")
	return &group, nil
}
