package main

import (
	"log"
	"log/slog"

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
}
