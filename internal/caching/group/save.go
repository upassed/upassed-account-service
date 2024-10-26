package group

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
	errMarshallingGroupData   = errors.New("unable to marshall group data to json format")
	errSavingGroupDataToCache = errors.New("unable to save group data to redis cache")
)

func (client *RedisClient) Save(ctx context.Context, group *domain.Group) error {
	_, span := otel.Tracer(client.cfg.Tracing.GroupTracerName).Start(ctx, "redisClient#Save")
	span.SetAttributes(
		attribute.String("specializationCode", group.SpecializationCode),
		attribute.String("groupNumber", group.GroupNumber),
	)
	defer span.End()

	log := logging.Wrap(client.log,
		logging.WithOp(client.Save),
		logging.WithCtx(ctx),
		logging.WithAny("groupID", group.ID),
	)

	log.Info("marshalling group data to json to save to cache")
	jsonGroupData, err := json.Marshal(group)
	if err != nil {
		log.Error("unable to marshall group data to json format")
		return errMarshallingGroupData
	}

	log.Info("saving group data to the cache")
	if err := client.client.Set(ctx, fmt.Sprintf(keyFormat, group.ID.String()), jsonGroupData, client.cfg.GetRedisEntityTTL()).Err(); err != nil {
		log.Error("error while saving group data to the cache", logging.Error(err))
		return errSavingGroupDataToCache
	}

	log.Info("group successfully saved to the cache")
	return nil
}
