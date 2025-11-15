package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	HTTPPort int `env:"HTTP_PORT" envDefault:"8080"`
	Postgres Postgres
}

type Postgres struct {
	URL            string `env:"PG_URL,required"`
	MaxConnections int    `env:"PG_MAX_CONNECTIONS,required"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
