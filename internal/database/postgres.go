package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	URL      string
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

func NewPostgresPool(
	ctx context.Context,
	cfg Config,
) (*pgxpool.Pool, error) {
	var dsn string
	if cfg.URL != "" {
		dsn = cfg.URL
	} else {
		ssl := cfg.SSLMode
		if ssl == "" {
			ssl = "disable"
		}
		dsn = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Name,
			ssl,
		)
	}

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse postgres config: %w", err)
	}

	poolConfig.MaxConns = 5
	poolConfig.MinConns = 0
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	return pool, nil
}

// Ping checks if the database is reachable with a short timeout.
func Ping(ctx context.Context, pool *pgxpool.Pool) error {
	if pool == nil {
		return fmt.Errorf("database pool is not initialized")
	}
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return pool.Ping(pingCtx)
}
