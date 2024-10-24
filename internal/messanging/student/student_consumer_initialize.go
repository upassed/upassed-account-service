package student

import (
	"errors"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/wagslane/go-rabbitmq"
	"log/slog"
	"reflect"
	"runtime"
)

var (
	errCreatingStudentCreateQueueConsumer = errors.New("unable to create student queue consumer")
	errRunningStudentCreateQueueConsumer  = errors.New("unable to run student queue consumer")
)

func InitializeCreateQueueConsumer(client *rabbitClient) error {
	op := runtime.FuncForPC(reflect.ValueOf(InitializeCreateQueueConsumer).Pointer()).Name()

	log := client.log.With(
		slog.String("op", op),
	)

	log.Info("started crating student create queue consumer")
	studentCreateGroupConsumer, err := rabbitmq.NewConsumer(
		client.rabbitConnection,
		client.cfg.Rabbit.Queues.StudentCreate.Name,
		rabbitmq.WithConsumerOptionsRoutingKey(client.cfg.Rabbit.Queues.StudentCreate.RoutingKey),
		rabbitmq.WithConsumerOptionsExchangeName(client.cfg.Rabbit.Exchange.Name),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)

	if err != nil {
		log.Error("unable to create student queue consumer", logging.Error(err))
		return errCreatingStudentCreateQueueConsumer
	}

	defer studentCreateGroupConsumer.Close()
	if err := studentCreateGroupConsumer.Run(CreateQueueConsumer(log)); err != nil {
		log.Error("unable to run student queue consumer")
		return errRunningStudentCreateQueueConsumer
	}

	return nil
}
