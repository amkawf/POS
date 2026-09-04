package orderhttp

import "github.com/google/uuid"

type CreateOrderRequest struct {
	CompanyID    uuid.UUID              `json:"company_id" binding:"required"`
	StoreID      uuid.UUID              `json:"store_id" binding:"required"`
	OrderType    string                 `json:"order_type" binding:"required"`
	OrderSource  string                 `json:"order_source" binding:"required"`
	CustomerName *string                `json:"customer_name"`
	Notes        *string                `json:"notes"`
	CreatedBy    *uuid.UUID             `json:"created_by"`
	Items        []CreateOrderItemInput `json:"items" binding:"required,min=1"`
}

type CreateOrderItemInput struct {
	MenuItemID uuid.UUID `json:"menu_item_id" binding:"required"`
	ItemName   string    `json:"item_name" binding:"required"`
	SKU        string    `json:"sku" binding:"required"`
	Quantity   int64     `json:"quantity" binding:"required,gt=0"`
	UnitPrice  int64     `json:"unit_price" binding:"gte=0"`
}

type CreateOrderResponse struct {
	ID          uuid.UUID `json:"id"`
	CompanyID   uuid.UUID `json:"company_id"`
	StoreID     uuid.UUID `json:"store_id"`
	OrderNumber string    `json:"order_number"`
	Status      string    `json:"status"`
	Subtotal    int64     `json:"subtotal"`
	TotalAmount int64     `json:"total_amount"`
}

type OrderResponse struct {
	ID             uuid.UUID           `json:"id"`
	CompanyID      uuid.UUID           `json:"company_id"`
	StoreID        uuid.UUID           `json:"store_id"`
	OrderNumber    string              `json:"order_number"`
	OrderType      string              `json:"order_type"`
	OrderSource    string              `json:"order_source"`
	Status         string              `json:"status"`
	CustomerName   *string             `json:"customer_name"`
	Subtotal       int64               `json:"subtotal"`
	DiscountAmount int64               `json:"discount_amount"`
	TaxAmount      int64               `json:"tax_amount"`
	ServiceAmount  int64               `json:"service_amount"`
	TotalAmount    int64               `json:"total_amount"`
	Notes          *string             `json:"notes"`
	OpenedAt       string              `json:"opened_at"`
	Items          []OrderItemResponse `json:"items,omitempty"`
}

type OrderItemResponse struct {
	ID             uuid.UUID `json:"id"`
	OrderID        uuid.UUID `json:"order_id"`
	MenuItemID     uuid.UUID `json:"menu_item_id"`
	ItemName       string    `json:"item_name"`
	SKU            string    `json:"sku"`
	Quantity       int64     `json:"quantity"`
	UnitPrice      int64     `json:"unit_price"`
	ModifierAmount int64     `json:"modifier_amount"`
	DiscountAmount int64     `json:"discount_amount"`
	TaxAmount      int64     `json:"tax_amount"`
	TotalAmount    int64     `json:"total_amount"`
	Notes          *string   `json:"notes"`
	Status         string    `json:"status"`
}

type PayOrderRequest struct {
	CompanyID     uuid.UUID `json:"company_id" binding:"required"`
	StoreID       uuid.UUID `json:"store_id" binding:"required"`
	PaymentMethod string    `json:"payment_method" binding:"required"`
	AmountPaid    int64     `json:"amount_paid" binding:"required,gte=0"`
	Notes         *string   `json:"notes"`
}
