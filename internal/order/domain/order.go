package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusDraft     OrderStatus = "DRAFT"
	OrderStatusOpen      OrderStatus = "OPEN"
	OrderStatusCompleted OrderStatus = "COMPLETED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
	OrderStatusVoid      OrderStatus = "VOID"
	OrderStatusRefunded  OrderStatus = "REFUNDED"
)

type OrderType string

const (
	OrderTypeDineIn   OrderType = "DINE_IN"
	OrderTypeTakeaway OrderType = "TAKEAWAY"
	OrderTypeDelivery OrderType = "DELIVERY"
)

type OrderSource string

const (
	OrderSourcePOS    OrderSource = "POS"
	OrderSourceWaiter OrderSource = "WAITER"
	OrderSourceQR     OrderSource = "QR"
)

type Order struct {
	ID                uuid.UUID
	CompanyID         uuid.UUID
	StoreID           uuid.UUID
	TableID           *uuid.UUID
	CustomerSessionID *uuid.UUID
	OrderNumber       string
	OrderType         OrderType
	OrderSource       OrderSource
	Status            OrderStatus
	CustomerName      *string
	Items             []OrderItem
	Subtotal          int64
	DiscountAmount    int64
	TaxAmount         int64
	ServiceAmount     int64
	TotalAmount       int64
	Notes             *string
	OpenedAt          time.Time
	CompletedAt       *time.Time
	CancelledAt       *time.Time
	CreatedBy         *uuid.UUID
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type OrderItem struct {
	ID             uuid.UUID
	OrderID        uuid.UUID
	MenuItemID     uuid.UUID
	ItemName       string
	SKU            string
	Quantity       int64
	UnitPrice      int64
	ModifierAmount int64
	DiscountAmount int64
	TaxAmount      int64
	TotalAmount    int64
	Notes          *string
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

var (
	ErrInvalidOrderType       = errors.New("invalid order type")
	ErrInvalidOrderSource     = errors.New("invalid order source")
	ErrInvalidOrderStatus     = errors.New("invalid order status")
	ErrEmptyOrder             = errors.New("order must contain at least one item")
	ErrInvalidItemQuantity    = errors.New("invalid item quantity")
	ErrInvalidItemPrice       = errors.New("invalid item price")
	ErrOrderNotEditable       = errors.New("order is not editable")
	ErrInvalidOrderTransition = errors.New("invalid order status transition")
)

func NewOrder(
	companyID uuid.UUID,
	storeID uuid.UUID,
	orderNumber string,
	orderType OrderType,
	orderSource OrderSource,
	createdBy *uuid.UUID,
) (*Order, error) {
	order := &Order{
		ID:          uuid.New(),
		CompanyID:   companyID,
		StoreID:     storeID,
		OrderNumber: orderNumber,
		OrderType:   orderType,
		OrderSource: orderSource,
		Status:      OrderStatusDraft,
		CreatedBy:   createdBy,
		Items:       make([]OrderItem, 0),
		OpenedAt:    time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := order.Validate(); err != nil {
		return nil, err
	}

	return order, nil
}

func (o Order) Validate() error {
	switch o.OrderType {
	case OrderTypeDineIn, OrderTypeTakeaway, OrderTypeDelivery:
	default:
		return ErrInvalidOrderType
	}

	switch o.OrderSource {
	case OrderSourcePOS, OrderSourceWaiter, OrderSourceQR:
	default:
		return ErrInvalidOrderSource
	}

	switch o.Status {
	case OrderStatusDraft,
		OrderStatusOpen,
		OrderStatusCompleted,
		OrderStatusCancelled,
		OrderStatusVoid,
		OrderStatusRefunded:
	default:
		return ErrInvalidOrderStatus
	}

	return nil
}

func (o *Order) AddItem(
	menuItemID uuid.UUID,
	itemName string,
	sku string,
	quantity int64,
	unitPrice int64,
) error {
	if o.Status != OrderStatusDraft &&
		o.Status != OrderStatusOpen {
		return ErrOrderNotEditable
	}

	if quantity <= 0 {
		return ErrInvalidItemQuantity
	}

	if unitPrice < 0 {
		return ErrInvalidItemPrice
	}

	item := OrderItem{
		ID:         uuid.New(),
		OrderID:    o.ID,
		MenuItemID: menuItemID,
		ItemName:   itemName,
		SKU:        sku,
		Quantity:   quantity,
		UnitPrice:  unitPrice,
		Status:     "ACTIVE",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	item.TotalAmount = item.UnitPrice * item.Quantity

	o.Items = append(o.Items, item)
	o.recalculateTotals()

	return nil
}

func (o *Order) Open() error {
	if o.Status != OrderStatusDraft {
		return ErrInvalidOrderTransition
	}

	if len(o.Items) == 0 {
		return ErrEmptyOrder
	}

	o.Status = OrderStatusOpen
	o.UpdatedAt = time.Now()

	return nil
}

func (o *Order) recalculateTotals() {
	var subtotal int64

	for _, item := range o.Items {
		if item.Status == "ACTIVE" {
			subtotal += item.TotalAmount
		}
	}

	o.Subtotal = subtotal
	o.TotalAmount =
		o.Subtotal -
			o.DiscountAmount +
			o.TaxAmount +
			o.ServiceAmount

	o.UpdatedAt = time.Now()
}
