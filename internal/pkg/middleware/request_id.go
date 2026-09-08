package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const HeaderXRequestID = "X-Request-ID"
const ContextKeyRequestID = "request_id"

// RequestID attaches a unique request ID to each incoming request context and response header.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader(HeaderXRequestID)
		if reqID == "" {
			reqID = uuid.New().String()
		}

		c.Set(ContextKeyRequestID, reqID)
		c.Header(HeaderXRequestID, reqID)
		c.Next()
	}
}

// GetRequestID retrieves the request ID from gin context.
func GetRequestID(c *gin.Context) string {
	if val, exists := c.Get(ContextKeyRequestID); exists {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}
