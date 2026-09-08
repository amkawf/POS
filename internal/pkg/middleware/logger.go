package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// StructuredLogger logs HTTP requests using standard library log/slog with latency, status, and request ID.
func StructuredLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		reqID := GetRequestID(c)

		fullPath := path
		if raw != "" {
			fullPath = path + "?" + raw
		}

		logAttrs := []any{
			slog.String("request_id", reqID),
			slog.String("method", method),
			slog.String("path", fullPath),
			slog.Int("status", statusCode),
			slog.Duration("latency", latency),
			slog.String("ip", clientIP),
		}

		if len(c.Errors) > 0 {
			logAttrs = append(logAttrs, slog.String("error", c.Errors.String()))
		}

		switch {
		case statusCode >= 500:
			logger.Error("HTTP request internal error", logAttrs...)
		case statusCode >= 400:
			logger.Warn("HTTP request client error", logAttrs...)
		default:
			logger.Info("HTTP request processed", logAttrs...)
		}
	}
}
