package app

import (
	"log/slog"

	config "github.com/upassed/upassed-account-service/internal/config/app"
	"github.com/upassed/upassed-account-service/internal/repository"
	"github.com/upassed/upassed-account-service/internal/server"
	"github.com/upassed/upassed-account-service/internal/service"
)

type App struct {
	Server *server.AppServer
}

func New(config *config.Config, log *slog.Logger) (*App, error) {
	const op = "app.New()"
	log = log.With(
		slog.String("op", op),
	)

	server := server.New(server.AppServerCreateParams{
		Config:         config,
		Log:            log,
		TeacherService: service.NewTeacherService(log, repository.NewTeacherRepository()),
	})

	log.Info("app successfully created")
	return &App{
		Server: server,
	}, nil
}
