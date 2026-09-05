package kitchenhttp

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pos-backend/internal/kitchen/application"
	"pos-backend/internal/kitchen/domain"
)

type Handler struct {
	listTicketsUseCase        *application.ListTicketsUseCase
	updateTicketStatusUseCase *application.UpdateTicketStatusUseCase
}

func NewHandler(
	listTicketsUseCase *application.ListTicketsUseCase,
	updateTicketStatusUseCase *application.UpdateTicketStatusUseCase,
) *Handler {
	return &Handler{
		listTicketsUseCase:        listTicketsUseCase,
		updateTicketStatusUseCase: updateTicketStatusUseCase,
	}
}

func (h *Handler) ListTickets(c *gin.Context) {
	companyIDStr := c.Query("company_id")
	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "company_id must be a valid UUID",
			},
		})
		return
	}

	storeIDStr := c.Query("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "store_id must be a valid UUID",
			},
		})
		return
	}

	var status *string
	if s := c.Query("status"); s != "" {
		status = &s
	}

	tickets, err := h.listTicketsUseCase.Execute(c.Request.Context(), application.ListTicketsInput{
		CompanyID: companyID,
		StoreID:   storeID,
		Status:    status,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "LIST_TICKETS_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	response := ListKitchenTicketsResponse{
		Tickets: make([]KitchenTicketResponse, 0, len(tickets)),
	}

	for _, t := range tickets {
		response.Tickets = append(response.Tickets, mapTicketToResponse(t))
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) UpdateTicketStatus(c *gin.Context) {
	ticketIDStr := c.Param("id")
	ticketID, err := uuid.Parse(ticketIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "id must be a valid UUID",
			},
		})
		return
	}

	var req UpdateTicketStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": err.Error(),
			},
		})
		return
	}

	companyID, err := uuid.Parse(req.CompanyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "company_id must be a valid UUID",
			},
		})
		return
	}

	storeID, err := uuid.Parse(req.StoreID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "store_id must be a valid UUID",
			},
		})
		return
	}

	targetStatus := domain.TicketStatus(req.Status)

	updated, err := h.updateTicketStatusUseCase.Execute(c.Request.Context(), application.UpdateTicketStatusInput{
		CompanyID: companyID,
		StoreID:   storeID,
		TicketID:  ticketID,
		Status:    targetStatus,
	})
	if err != nil {
		if errors.Is(err, domain.ErrTicketNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "TICKET_NOT_FOUND",
					"message": "Kitchen ticket not found",
				},
			})
			return
		}
		if errors.Is(err, domain.ErrInvalidStatusTransition) || errors.Is(err, domain.ErrTicketAlreadyFinalized) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_STATUS_TRANSITION",
					"message": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "UPDATE_TICKET_STATUS_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, mapTicketToResponse(*updated))
}

func mapTicketToResponse(t domain.KitchenTicket) KitchenTicketResponse {
	items := make([]KitchenTicketItemResponse, 0, len(t.Items))
	for _, it := range t.Items {
		items = append(items, KitchenTicketItemResponse{
			ID:          it.ID,
			TicketID:    it.TicketID,
			OrderItemID: it.OrderItemID,
			MenuItemID:  it.MenuItemID,
			ItemName:    it.ItemName,
			SKU:         it.SKU,
			Quantity:    it.Quantity,
			Notes:       it.Notes,
			Status:      string(it.Status),
			CreatedAt:   it.CreatedAt,
			UpdatedAt:   it.UpdatedAt,
		})
	}

	return KitchenTicketResponse{
		ID:          t.ID,
		CompanyID:   t.CompanyID,
		StoreID:     t.StoreID,
		OrderID:     t.OrderID,
		OrderNumber: t.OrderNumber,
		OrderType:   t.OrderType,
		TableID:     t.TableID,
		TableNumber: t.TableNumber,
		Status:      string(t.Status),
		Priority:    string(t.Priority),
		Notes:       t.Notes,
		Items:       items,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		StartedAt:   t.StartedAt,
		ReadyAt:     t.ReadyAt,
		ServedAt:    t.ServedAt,
	}
}
