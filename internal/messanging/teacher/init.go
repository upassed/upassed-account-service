package teacher

import (
	"errors"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/wagslane/go-rabbitmq"
	"log/slog"
	"reflect"
	"runtime"
)

var (
	errCreatingTeacherCreateQueueConsumer = errors.New("unable to create teacher queue consumer")
	errRunningTeacherCreateQueueConsumer  = errors.New("unable to run teacher queue consumer")
)

func InitializeCreateQueueConsumer(client *rabbitClient) error {
	op := runtime.FuncForPC(reflect.ValueOf(InitializeCreateQueueConsumer).Pointer()).Name()

	log := client.log.With(
		slog.String("op", op),
	)

	log.Info("started crating teacher create queue consumer")
	teacherCreateGroupConsumer, err := rabbitmq.NewConsumer(
		client.rabbitConnection,
		client.cfg.Rabbit.Queues.TeacherCreate.Name,
		rabbitmq.WithConsumerOptionsRoutingKey(client.cfg.Rabbit.Queues.TeacherCreate.RoutingKey),
		rabbitmq.WithConsumerOptionsExchangeName(client.cfg.Rabbit.Exchange.Name),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)

	if err != nil {
		log.Error("unable to create teacher queue consumer", logging.Error(err))
		return errCreatingTeacherCreateQueueConsumer
	}

	defer teacherCreateGroupConsumer.Close()
	if err := teacherCreateGroupConsumer.Run(client.CreateQueueConsumer(log)); err != nil {
		log.Error("unable to run teacher queue consumer")
		return errRunningTeacherCreateQueueConsumer
	}

	return nil
}
