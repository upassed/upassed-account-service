package app

import (
	"log/slog"

	config "github.com/upassed/upassed-account-service/internal/config/app"
	repository "github.com/upassed/upassed-account-service/internal/repository/postgres"
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

	teacherRepository, err := repository.NewTeacherRepository(config, log)
	if err != nil {
		return nil, err
	}

	server := server.New(server.AppServerCreateParams{
		Config:         config,
		Log:            log,
		TeacherService: service.NewTeacherService(log, teacherRepository),
	})

	log.Info("app successfully created")
	return &App{
		Server: server,
	}, nil
}
