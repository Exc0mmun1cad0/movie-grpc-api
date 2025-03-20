package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	MovieService MovieService `yaml:"movie_service"`
	Postgres     Postgres     `yaml:"postgres"`
}

type MovieService struct {
	Env string `yaml:"env" env:"ENV" env-default:"prod"`
	HTTPPort uint16 `yaml:"http_port" env:"HTTP_PORT" env-default:"8088"`
	GRPCPort uint16 `yaml:"grpc_port" env:"GRPC_PORT" env-default:"50051"`
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
	Port     uint16 `env:"POSTGRES_PORT" env-default:"5432"`
	Username string `env:"POSTGRES_USER" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	Database string `env:"POSTGRES_DB" env-default:"postgres"`
}

func Load() (*Config, error) {
	const op = "config.MustLoad"

	var cfg Config

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("%s: failed to load .env config: %w", op, err)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("%s: failed to read config from env vars: %w", op, err)
	}

	return &cfg, nil
}
