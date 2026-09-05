package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type TicketStatus string

const (
	TicketStatusPending   TicketStatus = "PENDING"
	TicketStatusPreparing TicketStatus = "PREPARING"
	TicketStatusReady     TicketStatus = "READY"
	TicketStatusServed    TicketStatus = "SERVED"
	TicketStatusCancelled TicketStatus = "CANCELLED"
)

type TicketPriority string

const (
	TicketPriorityNormal TicketPriority = "NORMAL"
	TicketPriorityRush   TicketPriority = "RUSH"
	TicketPriorityVIP    TicketPriority = "VIP"
)

var (
	ErrInvalidStatusTransition = errors.New("invalid kitchen ticket status transition")
	ErrTicketAlreadyFinalized  = errors.New("kitchen ticket is already in final state")
	ErrTicketNotFound          = errors.New("kitchen ticket not found")
)

type KitchenTicketItem struct {
	ID          uuid.UUID
	TicketID    uuid.UUID
	OrderItemID *uuid.UUID
	MenuItemID  uuid.UUID
	ItemName    string
	SKU         string
	Quantity    int64
	Notes       *string
	Status      TicketStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type KitchenTicket struct {
	ID          uuid.UUID
	CompanyID   uuid.UUID
	StoreID     uuid.UUID
	OrderID     uuid.UUID
	OrderNumber string
	OrderType   string
	TableID     *uuid.UUID
	TableNumber *string
	Status      TicketStatus
	Priority    TicketPriority
	Notes       *string
	Items       []KitchenTicketItem
	CreatedAt   time.Time
	UpdatedAt   time.Time
	StartedAt   *time.Time
	ReadyAt     *time.Time
	ServedAt    *time.Time
}

func NewKitchenTicket(
	companyID, storeID, orderID uuid.UUID,
	orderNumber, orderType string,
	tableID *uuid.UUID,
	priority TicketPriority,
	notes *string,
) *KitchenTicket {
	now := time.Now().UTC()
	if priority == "" {
		priority = TicketPriorityNormal
	}
	return &KitchenTicket{
		ID:          uuid.New(),
		CompanyID:   companyID,
		StoreID:     storeID,
		OrderID:     orderID,
		OrderNumber: orderNumber,
		OrderType:   orderType,
		TableID:     tableID,
		Status:      TicketStatusPending,
		Priority:    priority,
		Notes:       notes,
		Items:       make([]KitchenTicketItem, 0),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (t *KitchenTicket) AddItem(
	orderItemID *uuid.UUID,
	menuItemID uuid.UUID,
	itemName, sku string,
	quantity int64,
	notes *string,
) {
	now := time.Now().UTC()
	t.Items = append(t.Items, KitchenTicketItem{
		ID:          uuid.New(),
		TicketID:    t.ID,
		OrderItemID: orderItemID,
		MenuItemID:  menuItemID,
		ItemName:    itemName,
		SKU:         sku,
		Quantity:    quantity,
		Notes:       notes,
		Status:      TicketStatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
}

func (t *KitchenTicket) TransitionTo(targetStatus TicketStatus) error {
	if t.Status == TicketStatusServed || t.Status == TicketStatusCancelled {
		return ErrTicketAlreadyFinalized
	}

	now := time.Now().UTC()

	switch targetStatus {
	case TicketStatusPreparing:
		if t.Status != TicketStatusPending {
			return ErrInvalidStatusTransition
		}
		t.Status = TicketStatusPreparing
		t.StartedAt = &now
	case TicketStatusReady:
		if t.Status != TicketStatusPreparing && t.Status != TicketStatusPending {
			return ErrInvalidStatusTransition
		}
		t.Status = TicketStatusReady
		if t.StartedAt == nil {
			t.StartedAt = &now
		}
		t.ReadyAt = &now
	case TicketStatusServed:
		if t.Status != TicketStatusReady && t.Status != TicketStatusPreparing {
			return ErrInvalidStatusTransition
		}
		t.Status = TicketStatusServed
		t.ServedAt = &now
	case TicketStatusCancelled:
		t.Status = TicketStatusCancelled
	default:
		return ErrInvalidStatusTransition
	}

	t.UpdatedAt = now
	for i := range t.Items {
		t.Items[i].Status = targetStatus
		t.Items[i].UpdatedAt = now
	}

	return nil
}
