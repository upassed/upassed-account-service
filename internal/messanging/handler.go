package messanging

import (
	"context"
	"github.com/wagslane/go-rabbitmq"
)

type HandlerWithContext func(ctx context.Context, d rabbitmq.Delivery) (action rabbitmq.Action)
