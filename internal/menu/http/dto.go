package menuhttp

import "github.com/google/uuid"

// CreateMenuItemRequest mengatur payload JSON saat kasir/manajer mendaftarkan menu baru
type CreateMenuItemRequest struct {
	CompanyID       uuid.UUID   `json:"company_id" binding:"required"`
	SKU             string      `json:"sku" binding:"required"`
	Name            string      `json:"name" binding:"required"`
	Description     *string     `json:"description"`
	BasePrice       int64       `json:"base_price" binding:"required"`
	CategoryIDs     []uuid.UUID `json:"category_ids"`
	FulfillmentType string      `json:"fulfillment_type"` // BATCH_COOKING atau MADE_TO_ORDER
}

type MenuItemResponse struct {
	ID          uuid.UUID   `json:"id"`
	SKU         string      `json:"sku"`
	Name        string      `json:"name"`
	Description *string     `json:"description,omitempty"`
	BasePrice   int64       `json:"base_price"`
	CategoryIDs []uuid.UUID `json:"category_ids"`
	Stock       *int64      `json:"stock,omitempty"`
	FulfillmentType string      `json:"fulfillment_type"`
}

type ListMenuItemsResponse struct {
	Items []MenuItemResponse `json:"items"`
}

type CategoryResponse struct {
	ID        uuid.UUID `json:"id"`
	MenuID    uuid.UUID `json:"menu_id"`
	Name      string    `json:"name"`
	SortOrder int32     `json:"sort_order"`
}

type ListCategoriesResponse struct {
	Categories []CategoryResponse `json:"categories"`
}
