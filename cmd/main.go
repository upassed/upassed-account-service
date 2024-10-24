package main

import (
	"github.com/upassed/upassed-account-service/internal/tracing"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"runtime"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/upassed/upassed-account-service/internal/app"
	"github.com/upassed/upassed-account-service/internal/config"
	"github.com/upassed/upassed-account-service/internal/logging"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	if err := os.Setenv(config.EnvConfigPath, filepath.Join("config", "local.yml")); err != nil {
		log.Fatal(err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	logger := logging.New(cfg.Env).With(
		slog.String("op", runtime.FuncForPC(reflect.ValueOf(main).Pointer()).Name()),
	)

	logger.Info("logger successfully initialized", slog.Any("env", cfg.Env))

	traceProviderShutdownFunc, err := tracing.InitTracer(cfg, logger)
	if err != nil {
		logger.Error("unable to initialize traceProvider", logging.Error(err))
		os.Exit(1)
	}

	defer traceProviderShutdownFunc()
	logger.Info("trace provider successfully initialized")

	application, err := app.New(cfg, logger)
	if err != nil {
		logger.Error("error occurred while creating an app", logging.Error(err))
		os.Exit(1)
	}

	go func(app *app.App) {
		if err := app.Server.Run(); err != nil {
			logger.Error("error occurred while running a gRPC server", logging.Error(err))
			os.Exit(1)
		}
	}(application)

	stopSignalChannel := make(chan os.Signal, 1)
	signal.Notify(stopSignalChannel, syscall.SIGTERM, syscall.SIGINT)
	<-stopSignalChannel

	application.Server.GracefulStop()
	_ = application.RabbitConn.Close()
	logger.Info("server gracefully stopped")
}
