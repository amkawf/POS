package httputil

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ParseUUIDParam parses a URL path parameter as UUID. If invalid, responds with 400 and returns false.
func ParseUUIDParam(c *gin.Context, paramName string) (uuid.UUID, bool) {
	val := c.Param(paramName)
	id, err := uuid.Parse(val)
	if err != nil {
		BadRequest(c, "INVALID_REQUEST", fmt.Sprintf("%s must be a valid UUID", paramName))
		return uuid.Nil, false
	}
	return id, true
}

// ParseUUIDQuery parses a required query string parameter as UUID. If invalid or missing, responds with 400 and returns false.
func ParseUUIDQuery(c *gin.Context, queryName string) (uuid.UUID, bool) {
	val := strings.TrimSpace(c.Query(queryName))
	if val == "" {
		BadRequest(c, "INVALID_REQUEST", fmt.Sprintf("%s query parameter is required", queryName))
		return uuid.Nil, false
	}

	id, err := uuid.Parse(val)
	if err != nil {
		BadRequest(c, "INVALID_REQUEST", fmt.Sprintf("%s must be a valid UUID", queryName))
		return uuid.Nil, false
	}
	return id, true
}

// ParseOptUUIDQuery parses an optional query string parameter as UUID. If missing, returns nil, true. If malformed, responds with 400 and returns nil, false.
func ParseOptUUIDQuery(c *gin.Context, queryName string) (*uuid.UUID, bool) {
	val := strings.TrimSpace(c.Query(queryName))
	if val == "" {
		return nil, true
	}

	id, err := uuid.Parse(val)
	if err != nil {
		BadRequest(c, "INVALID_REQUEST", fmt.Sprintf("%s must be a valid UUID", queryName))
		return nil, false
	}
	return &id, true
}

// ParseOptStringQuery parses an optional query parameter. Returns nil if missing or empty.
func ParseOptStringQuery(c *gin.Context, queryName string) *string {
	val := strings.TrimSpace(c.Query(queryName))
	if val == "" {
		return nil
	}
	return &val
}

// ParseIntQuery parses an optional integer query parameter, falling back to defaultVal.
func ParseIntQuery(c *gin.Context, queryName string, defaultVal int32) int32 {
	val := strings.TrimSpace(c.Query(queryName))
	if val == "" {
		return defaultVal
	}
	parsed, err := strconv.Atoi(val)
	if err != nil || parsed <= 0 {
		return defaultVal
	}
	return int32(parsed)
}
