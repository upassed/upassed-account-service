package student

import (
	"errors"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/wagslane/go-rabbitmq"
)

var (
	errCreatingStudentCreateQueueConsumer = errors.New("unable to create student queue consumer")
	errRunningStudentCreateQueueConsumer  = errors.New("unable to run student queue consumer")
)

func InitializeCreateQueueConsumer(client *rabbitClient) error {
	log := logging.Wrap(client.log,
		logging.WithOp(InitializeCreateQueueConsumer),
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
	if err := studentCreateGroupConsumer.Run(client.CreateQueueConsumer()); err != nil {
		log.Error("unable to run student queue consumer")
		return errRunningStudentCreateQueueConsumer
	}

	log.Info("student queue consumer successfully initialized")
	return nil
}
