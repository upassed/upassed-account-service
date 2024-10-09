package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	ErrorConfigFlagEmpty   error = errors.New("config flag is not passed")
	ErrorConfigEnvEmpty    error = errors.New("config path env is not set")
	ErrorConfigFileInvalid error = errors.New("config file has invalid format")
)

type EnvType string

const (
	EnvLocal   EnvType = "local"
	EnvDev     EnvType = "dev"
	EnvTesting EnvType = "testing"
)

const EnvConfigPath string = "APP_CONFIG_PATH"

type Config struct {
	Env        EnvType    `yaml:"env" env-required:"true"`
	Storage    Storage    `yaml:"storage" env-required:"true"`
	GrpcServer GrpcServer `yaml:"grpc_server" env-required:"true"`
}

type Storage struct {
	Host         string `yaml:"host" env-required:"true"`
	Port         string `yaml:"port" env-required:"true"`
	DatabaseName string `yaml:"database_name" env-required:"true"`
	User         string `yaml:"user" env-required:"true"`
	Password     string `yaml:"password" env-required:"true"`
}

type GrpcServer struct {
	Port    string        `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-required:"true"`
}

func Load() (*Config, error) {
	const op = "config.Load()"

	pathToConfig := os.Getenv(EnvConfigPath)
	if pathToConfig == "" {
		return nil, fmt.Errorf("%s -> %w", op, ErrorConfigEnvEmpty)
	}

	return loadByPath(pathToConfig)
}

func loadByPath(pathToConfig string) (*Config, error) {
	const op = "config.loadByPath()"

	var config Config

	if err := cleanenv.ReadConfig(pathToConfig, &config); err != nil {
		return nil, fmt.Errorf("%s -> %w; %w", op, ErrorConfigFileInvalid, err)
	}

	return &config, nil
}
