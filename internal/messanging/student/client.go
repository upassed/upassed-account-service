package student

import (
	"github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/service/student"
	"github.com/wagslane/go-rabbitmq"
	"log/slog"
)

type rabbitClient struct {
	service          student.Service
	rabbitConnection *rabbitmq.Conn
	cfg              *config.Config
	log              *slog.Logger
}

func Initialize(service student.Service, rabbitConnection *rabbitmq.Conn, cfg *config.Config, log *slog.Logger) {
	log = logging.Wrap(log,
		logging.WithOp(Initialize),
	)

	client := &rabbitClient{
		service:          service,
		rabbitConnection: rabbitConnection,
		cfg:              cfg,
		log:              log,
	}

	go func() {
		if err := InitializeCreateQueueConsumer(client); err != nil {
			log.Error("error while initializing student queue consumer", logging.Error(err))
			return
		}
	}()
}