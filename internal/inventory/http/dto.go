package inventoryhttp

import "github.com/google/uuid"

type AdjustStockRequest struct {
	CompanyID  uuid.UUID `json:"company_id" binding:"required"`
	StoreID    uuid.UUID `json:"store_id" binding:"required"`
	MenuItemID uuid.UUID `json:"menu_item_id" binding:"required"`
	Quantity   int64     `json:"quantity" binding:"required"`
	Notes      string    `json:"notes"`
}

type AdjustStockResponse struct {
	Message string `json:"message"`
}