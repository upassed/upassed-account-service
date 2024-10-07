package config

import (
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	ErrorConfigFlagEmpty   error = errors.New("config flag is not passed")
	ErrorConfigFileInvalid error = errors.New("config file has invalid format")
)

type EnvType string

const (
	EnvLocal EnvType = "local"
	EnvDev   EnvType = "dev"
)

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

	var configFilename string

	flag.StringVar(&configFilename, "config", "", "path to config file")
	flag.Parse()

	if configFilename == "" {
		return nil, fmt.Errorf("%s -> %w", op, ErrorConfigFlagEmpty)
	}

	return loadByPath(configFilename)
}

func loadByPath(pathToConfig string) (*Config, error) {
	const op = "config.LoadByPath()"

	var config Config

	if err := cleanenv.ReadConfig(pathToConfig, &config); err != nil {
		return nil, fmt.Errorf("%s -> %w; %w", op, ErrorConfigFileInvalid, err)
	}

	return &config, nil
}
