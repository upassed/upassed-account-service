package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/upassed/upassed-account-service/internal/app"
	config "github.com/upassed/upassed-account-service/internal/config/app"
	"github.com/upassed/upassed-account-service/internal/logger"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	log := logger.New(config.Env)
	log.Info("logger successfully initialized", slog.Any("env", config.Env))

	application, err := app.New(config, log)
	if err != nil {
		log.Error("error occured while creating an app", logger.Error(err))
		os.Exit(1)
	}

	go func(app *app.App) {
		if err := app.Server.Run(); err != nil {
			log.Error("error occured while running a gRPC server", logger.Error(err))
			os.Exit(1)
		}
	}(application)

	stopSignalChannel := make(chan os.Signal, 1)
	signal.Notify(stopSignalChannel, syscall.SIGTERM, syscall.SIGINT)
	<-stopSignalChannel

	application.Server.GracefulStop()
	log.Info("server gracefully stopped")
}
