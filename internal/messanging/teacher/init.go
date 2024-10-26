package teacher

import (
	"errors"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/wagslane/go-rabbitmq"
)

var (
	errCreatingTeacherCreateQueueConsumer = errors.New("unable to create teacher queue consumer")
	errRunningTeacherCreateQueueConsumer  = errors.New("unable to run teacher queue consumer")
)

func InitializeCreateQueueConsumer(client *rabbitClient) error {
	log := logging.Wrap(client.log,
		logging.WithOp(InitializeCreateQueueConsumer),
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
	if err := teacherCreateGroupConsumer.Run(client.CreateQueueConsumer()); err != nil {
		log.Error("unable to run teacher queue consumer")
		return errRunningTeacherCreateQueueConsumer
	}

	log.Info("teacher queue consumer successfully created")
	return nil
}
