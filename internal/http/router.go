package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"pos-backend/internal/database"
	"pos-backend/internal/pkg/middleware"
)

// RouterConfig holds dependencies for the HTTP server router.
type RouterConfig struct {
	Logger   *slog.Logger
	Database *pgxpool.Pool
	Version  string
}

// NewRouter constructs a configured gin.Engine with production-grade middlewares and health routes.
func NewRouter(cfg RouterConfig) *gin.Engine {
	// Set gin mode to Release if in production
	router := gin.New()

	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}

	router.Use(middleware.RequestID())
	router.Use(middleware.StructuredLogger(cfg.Logger))
	router.Use(middleware.Recovery(cfg.Logger))
	router.Use(middleware.CORS())

	// Health check endpoint (supports container probes & k8s readiness)
	router.GET("/health", func(c *gin.Context) {
		status := "ok"
		dbStatus := "disabled"

		if cfg.Database != nil {
			ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
			defer cancel()

			if err := database.Ping(ctx, cfg.Database); err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"status":   "degraded",
					"database": "unreachable",
					"error":    err.Error(),
				})
				return
			}
			dbStatus = "connected"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   status,
			"database": dbStatus,
			"version":  cfg.Version,
		})
	})

	return router
}
