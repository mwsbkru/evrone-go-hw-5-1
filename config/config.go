package config

import (
	"errors"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Host               string `env:"HOST" env-default:"0.0.0.0"`
	Port               string `env:"PORT" env-default:"8080"`
	DbConnectionString string `env:"DB_CONNECTION_STRING"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("Не удалось прочитать параметры конфига: %w", err)
	}

	if cfg.DbConnectionString == "" {
		return nil, errors.New("connection string for DB is required")
	}

	return &cfg, nil
}
