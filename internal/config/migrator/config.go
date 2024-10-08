package config

import (
	"errors"
	"flag"
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	ErrorConfigFlagEmpty   error = errors.New("config flag is not passed")
	ErrorConfigFileInvalid error = errors.New("config file has invalid format")
)

type Config struct {
	MigrationsPath      string  `yaml:"migrations_path" env-required:"true"`
	MigrationsTableName string  `yaml:"migrations_table_name" env-default:"migrations"`
	Storage             Storage `yaml:"storage" env-required:"true"`
}

type Storage struct {
	Host         string `yaml:"host" env-required:"true"`
	Port         string `yaml:"port" env-required:"true"`
	DatabaseName string `yaml:"database_name" env-required:"true"`
	User         string `yaml:"user" env-required:"true"`
	Password     string `yaml:"password" env-required:"true"`
}

func Load() (*Config, error) {
	var configFilename string

	flag.StringVar(&configFilename, "config", "", "path to config file")
	flag.Parse()

	if configFilename == "" {
		return nil, ErrorConfigFlagEmpty
	}

	return loadByPath(configFilename)
}

func loadByPath(pathToConfig string) (*Config, error) {
	const op = "config.LoadByPath()"

	var config Config

	if err := cleanenv.ReadConfig(pathToConfig, &config); err != nil {
		return nil, fmt.Errorf("%s - %w; %w", op, ErrorConfigFileInvalid, err)
	}

	return &config, nil
}
