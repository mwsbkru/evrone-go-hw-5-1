package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Host string `env:"HOST" env-default:"0.0.0.0"`
	Port string `env:"PORT" env-default:"8080"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("Не удалось прочитать параметры конфига: %w", err)
	}

	return &cfg, nil
}
