package httputil

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse defines the standard error envelope matching the API contract.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// JSON sends a JSON response with status code.
func JSON(c *gin.Context, statusCode int, data any) {
	c.JSON(statusCode, data)
}

// Error sends a standardized JSON error response.
func Error(c *gin.Context, statusCode int, code, message string) {
	_ = c.Error(errors.New(message))
	c.JSON(statusCode, ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}

// BadRequest sends a 400 Bad Request error.
func BadRequest(c *gin.Context, code, message string) {
	Error(c, http.StatusBadRequest, code, message)
}

// NotFound sends a 404 Not Found error.
func NotFound(c *gin.Context, code, message string) {
	Error(c, http.StatusNotFound, code, message)
}

// Conflict sends a 409 Conflict error.
func Conflict(c *gin.Context, code, message string) {
	Error(c, http.StatusConflict, code, message)
}

// InternalError sends a 500 Internal Server Error.
func InternalError(c *gin.Context, code, message string) {
	Error(c, http.StatusInternalServerError, code, message)
}
