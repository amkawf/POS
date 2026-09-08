package kitchenhttp

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pos-backend/internal/kitchen/application"
	"pos-backend/internal/kitchen/domain"
	"pos-backend/internal/pkg/httputil"
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

// RegisterRoutes registers the kitchen module routes to the provided router group.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/kitchen/tickets", h.ListTickets)
	rg.PATCH("/kitchen/tickets/:id/status", h.UpdateTicketStatus)
}

func (h *Handler) ListTickets(c *gin.Context) {
	companyID, ok := httputil.ParseUUIDQuery(c, "company_id")
	if !ok {
		return
	}

	storeID, ok := httputil.ParseUUIDQuery(c, "store_id")
	if !ok {
		return
	}

	status := httputil.ParseOptStringQuery(c, "status")

	tickets, err := h.listTicketsUseCase.Execute(c.Request.Context(), application.ListTicketsInput{
		CompanyID: companyID,
		StoreID:   storeID,
		Status:    status,
	})
	if err != nil {
		httputil.InternalError(c, "LIST_TICKETS_FAILED", err.Error())
		return
	}

	response := ListKitchenTicketsResponse{
		Tickets: make([]KitchenTicketResponse, 0, len(tickets)),
	}

	for _, t := range tickets {
		response.Tickets = append(response.Tickets, mapTicketToResponse(t))
	}

	httputil.JSON(c, http.StatusOK, response)
}

func (h *Handler) UpdateTicketStatus(c *gin.Context) {
	ticketID, ok := httputil.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	var req UpdateTicketStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "INVALID_REQUEST", err.Error())
		return
	}

	companyID, err := uuid.Parse(req.CompanyID)
	if err != nil {
		httputil.BadRequest(c, "INVALID_REQUEST", "company_id must be a valid UUID")
		return
	}

	storeID, err := uuid.Parse(req.StoreID)
	if err != nil {
		httputil.BadRequest(c, "INVALID_REQUEST", "store_id must be a valid UUID")
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
			httputil.NotFound(c, "TICKET_NOT_FOUND", "Kitchen ticket not found")
			return
		}
		if errors.Is(err, domain.ErrInvalidStatusTransition) || errors.Is(err, domain.ErrTicketAlreadyFinalized) {
			httputil.BadRequest(c, "INVALID_STATUS_TRANSITION", err.Error())
			return
		}

		httputil.InternalError(c, "UPDATE_TICKET_STATUS_FAILED", err.Error())
		return
	}

	httputil.JSON(c, http.StatusOK, mapTicketToResponse(*updated))
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
