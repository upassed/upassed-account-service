package amqp

import (
	"context"
	"github.com/upassed/upassed-account-service/internal/messanging"
)

type Middleware func(ctx context.Context, handler messanging.HandlerWithContext) messanging.HandlerWithContext

func ChainMiddleware(handler messanging.HandlerWithContext, middlewares ...Middleware) messanging.HandlerWithContext {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](context.Background(), handler)
	}

	return handler
}
