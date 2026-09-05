package kitchenhttp

import (
	"time"

	"github.com/google/uuid"
)

type KitchenTicketItemResponse struct {
	ID          uuid.UUID  `json:"id"`
	TicketID    uuid.UUID  `json:"ticket_id"`
	OrderItemID *uuid.UUID `json:"order_item_id,omitempty"`
	MenuItemID  uuid.UUID  `json:"menu_item_id"`
	ItemName    string     `json:"item_name"`
	SKU         string     `json:"sku"`
	Quantity    int64      `json:"quantity"`
	Notes       *string    `json:"notes,omitempty"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type KitchenTicketResponse struct {
	ID          uuid.UUID                   `json:"id"`
	CompanyID   uuid.UUID                   `json:"company_id"`
	StoreID     uuid.UUID                   `json:"store_id"`
	OrderID     uuid.UUID                   `json:"order_id"`
	OrderNumber string                      `json:"order_number"`
	OrderType   string                      `json:"order_type"`
	TableID     *uuid.UUID                  `json:"table_id,omitempty"`
	TableNumber *string                     `json:"table_number,omitempty"`
	Status      string                      `json:"status"`
	Priority    string                      `json:"priority"`
	Notes       *string                     `json:"notes,omitempty"`
	Items       []KitchenTicketItemResponse `json:"items"`
	CreatedAt   time.Time                   `json:"created_at"`
	UpdatedAt   time.Time                   `json:"updated_at"`
	StartedAt   *time.Time                  `json:"started_at,omitempty"`
	ReadyAt     *time.Time                  `json:"ready_at,omitempty"`
	ServedAt    *time.Time                  `json:"served_at,omitempty"`
}

type ListKitchenTicketsResponse struct {
	Tickets []KitchenTicketResponse `json:"tickets"`
}

type UpdateTicketStatusRequest struct {
	CompanyID string `json:"company_id" binding:"required"`
	StoreID   string `json:"store_id" binding:"required"`
	Status    string `json:"status" binding:"required"`
}
