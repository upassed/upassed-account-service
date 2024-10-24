package student

import (
	"github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/logging"
	"github.com/upassed/upassed-account-service/internal/service/student"
	"github.com/wagslane/go-rabbitmq"
	"log/slog"
	"reflect"
	"runtime"
)

type rabbitClient struct {
	service          student.Service
	rabbitConnection *rabbitmq.Conn
	cfg              *config.Config
	log              *slog.Logger
}

func Initialize(service student.Service, rabbitConnection *rabbitmq.Conn, cfg *config.Config, log *slog.Logger) {
	op := runtime.FuncForPC(reflect.ValueOf(Initialize).Pointer()).Name()

	log = log.With(
		slog.String("op", op),
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
