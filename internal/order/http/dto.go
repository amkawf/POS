package orderhttp

import "github.com/google/uuid"

type CreateOrderRequest struct {
	CompanyID    uuid.UUID              `json:"company_id" binding:"required"`
	StoreID      uuid.UUID              `json:"store_id" binding:"required"`
	OrderNumber  string                 `json:"order_number" binding:"required"`
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