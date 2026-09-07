package config

import (
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

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func Load() (Config, error) {
	cfg := Config{
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "8080"),

		Database: DatabaseConfig{
			Host:     getEnv("DATABASE_HOST", "localhost"),
			Port:     getEnv("DATABASE_PORT", "5433"),
			Name:     getEnv("DATABASE_NAME", "pos"),
			User:     getEnv("DATABASE_USER", "pos"),
			Password: getEnv("DATABASE_PASSWORD", "pos_dev_password"),
		},
	}

	return cfg, nil
}
