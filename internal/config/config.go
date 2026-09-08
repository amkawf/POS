package config

import (
	"bufio"
	"os"
	"strings"
)

type Config struct {
	AppEnv  string
	AppPort string

	Database DatabaseConfig
}

type DatabaseConfig struct {
	URL      string
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// loadDotEnv looks for a .env file and sets environment variables if not already set.
func loadDotEnv() {
	paths := []string{".env", "../.env"}
	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				val = strings.Trim(val, `"'`)
				if os.Getenv(key) == "" {
					_ = os.Setenv(key, val)
				}
			}
		}
		if err := scanner.Err(); err != nil {
			// ignore read error on dot env
		}
		break
	}
}

const DefaultDatabaseURL = "postgresql://postgres.gkmkfwafueezngtszqhz:P0s3.142857phi@aws-0-ap-southeast-1.pooler.supabase.com:5432/postgres"

func Load() (Config, error) {
	loadDotEnv()

	cfg := Config{
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "8080"),

		Database: DatabaseConfig{
			URL:      getEnv("DATABASE_URL", DefaultDatabaseURL),
			Host:     getEnv("DATABASE_HOST", "localhost"),
			Port:     getEnv("DATABASE_PORT", "5433"),
			Name:     getEnv("DATABASE_NAME", "pos"),
			User:     getEnv("DATABASE_USER", "pos"),
			Password: getEnv("DATABASE_PASSWORD", "pos_dev_password"),
			SSLMode:  getEnv("DATABASE_SSLMODE", "disable"),
		},
	}

	return cfg, nil
}
