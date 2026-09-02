package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppEnv  string
	AppPort string

	Database DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func Load() (Config, error) {
	cfg := Config{
		AppEnv:  os.Getenv("APP_ENV"),
		AppPort: os.Getenv("APP_PORT"),

		Database: DatabaseConfig{
			Host:     os.Getenv("DATABASE_HOST"),
			Port:     os.Getenv("DATABASE_PORT"),
			Name:     os.Getenv("DATABASE_NAME"),
			User:     os.Getenv("DATABASE_USER"),
			Password: os.Getenv("DATABASE_PASSWORD"),
		},
	}

	if cfg.AppEnv == "" {
		cfg.AppEnv = "development"
	}

	if cfg.AppPort == "" {
		cfg.AppPort = "8080"
	}

	if cfg.AppPort == "" {
		return Config{}, fmt.Errorf("APP_PORT is required")
	}

	if cfg.Database.Port == "" {
		cfg.Database.Port = "5432"
	}

	return cfg, nil
}
