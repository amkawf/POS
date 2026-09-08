package httputil

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"

	kitchendomain "pos-backend/internal/kitchen/domain"
	orderdomain "pos-backend/internal/order/domain"
)

// HandleError inspects domain and repository errors and maps them to appropriate HTTP status codes.
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// 404 Not Found
	if errors.Is(err, kitchendomain.ErrTicketNotFound) {
		NotFound(c, "TICKET_NOT_FOUND", err.Error())
		return
	}
	if strings.Contains(strings.ToLower(err.Error()), "not found") {
		NotFound(c, "RESOURCE_NOT_FOUND", err.Error())
		return
	}

	// 409 Conflict
	if errors.Is(err, orderdomain.ErrOrderAlreadyCompleted) {
		Conflict(c, "ORDER_ALREADY_COMPLETED", err.Error())
		return
	}
	if errors.Is(err, orderdomain.ErrInvalidOrderTransition) ||
		errors.Is(err, kitchendomain.ErrInvalidStatusTransition) {
		Conflict(c, "INVALID_STATUS_TRANSITION", err.Error())
		return
	}

	// 400 Bad Request
	if errors.Is(err, orderdomain.ErrInvalidOrderType) ||
		errors.Is(err, orderdomain.ErrInvalidOrderSource) ||
		errors.Is(err, orderdomain.ErrInvalidOrderStatus) ||
		errors.Is(err, orderdomain.ErrEmptyOrder) ||
		errors.Is(err, orderdomain.ErrInvalidItemQuantity) ||
		errors.Is(err, orderdomain.ErrInvalidItemPrice) ||
		errors.Is(err, orderdomain.ErrOrderNotEditable) {
		BadRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}

	// Default 500 Internal Server Error
	InternalError(c, "INTERNAL_ERROR", err.Error())
}
