package teacher

import (
	"github.com/wagslane/go-rabbitmq"
	"log/slog"
	"reflect"
	"runtime"
)

func CreateQueueConsumer(log *slog.Logger) func(d rabbitmq.Delivery) rabbitmq.Action {
	op := runtime.FuncForPC(reflect.ValueOf(CreateQueueConsumer).Pointer()).Name()

	log = log.With(
		slog.String("op", op),
	)

	return func(delivery rabbitmq.Delivery) rabbitmq.Action {
		log.Info("consumed teacher create message", slog.String("messageBody", string(delivery.Body)))
		return rabbitmq.Ack
	}
}
