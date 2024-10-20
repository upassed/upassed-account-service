package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	errConfigEnvEmpty    = errors.New("config path env is not set")
	errConfigFileInvalid = errors.New("config file has invalid format")
)

type EnvType string

const (
	EnvLocal   EnvType = "local"
	EnvDev     EnvType = "dev"
	EnvTesting EnvType = "testing"

	EnvConfigPath string = "APP_CONFIG_PATH"
)

type Config struct {
	Env        EnvType         `yaml:"env" env-required:"true"`
	Storage    Storage         `yaml:"storage" env-required:"true"`
	GrpcServer GrpcServer      `yaml:"grpc_server" env-required:"true"`
	Migration  MigrationConfig `yaml:"migrations" env-required:"true"`
	Timeouts   Timeouts        `yaml:"timeouts" env-required:"true"`
}

type Storage struct {
	Host         string `yaml:"host" env:"POSTGRES_HOST" env-required:"true"`
	Port         string `yaml:"port" env:"POSTGRES_PORT" env-required:"true"`
	DatabaseName string `yaml:"database_name" env:"POSTGRES_DATABASE_NAME" env-required:"true"`
	User         string `yaml:"user" env:"POSTGRES_USER" env-required:"true"`
	Password     string `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
}

type GrpcServer struct {
	Port    string `yaml:"port" env:"GRPC_SERVER_PORT" env-required:"true"`
	Timeout string `yaml:"timeout" env:"GRPC_SERVER_TIMEOUT" env-required:"true"`
}

type MigrationConfig struct {
	MigrationsPath      string `yaml:"migrations_path" env:"MIGRATIONS_PATH" env-required:"true"`
	MigrationsTableName string `yaml:"migrations_table_name" env:"MIGRATIONS_TABLE_NAME" env-default:"migrations"`
}

type Timeouts struct {
	EndpointExecutionTimeoutMS string `yaml:"endpoint_execution_timeout_ms" env:"ENDPOINT_EXECUTION_TIMEOUT_MS" env-required:"true"`
}

func Load() (*Config, error) {
	op := runtime.FuncForPC(reflect.ValueOf(Load).Pointer()).Name()

	pathToConfig := os.Getenv(EnvConfigPath)
	if pathToConfig == "" {
		return nil, fmt.Errorf("%s -> %w", op, errConfigEnvEmpty)
	}

	return loadByPath(pathToConfig)
}

func loadByPath(pathToConfig string) (*Config, error) {
	op := runtime.FuncForPC(reflect.ValueOf(loadByPath).Pointer()).Name()

	var config Config
	if err := cleanenv.ReadConfig(pathToConfig, &config); err != nil {
		return nil, fmt.Errorf("%s -> %w; %w", op, errConfigFileInvalid, err)
	}

	return &config, nil
}

func (cfg *Config) GetEndpointExecutionTimeout() time.Duration {
	op := runtime.FuncForPC(reflect.ValueOf(cfg.GetEndpointExecutionTimeout).Pointer()).Name()

	milliseconds, err := strconv.Atoi(cfg.Timeouts.EndpointExecutionTimeoutMS)
	if err != nil {
		log.Fatal(fmt.Sprintf("%s, op=%s, err=%s", "unable to convert endpoint timeout duration", op, err.Error()))
	}

	return time.Duration(milliseconds) * time.Millisecond
}
