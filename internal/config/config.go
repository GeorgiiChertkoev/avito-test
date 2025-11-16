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
	User           string `env:"POSTGRES_USER,required"`
	Password       string `env:"POSTGRES_PASSWORD,required"`
	DatabaseName   string `env:"POSTGRES_DBNAME" envDefault:"prreviewer"`
	Host           string `env:"POSTGRES_HOST" envDefault:"postgres"`
	Port           int    `env:"POSTGRES_PORT" envDefault:"5432"`
	MaxConnections int    `env:"PG_MAX_CONNECTIONS" envDefault:"2"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
