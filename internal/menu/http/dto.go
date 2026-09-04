package menuhttp

import "github.com/google/uuid"

type MenuItemResponse struct {
	ID          uuid.UUID `json:"id"`
	SKU         string    `json:"sku"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	BasePrice   int64     `json:"base_price"`
}

type ListMenuItemsResponse struct {
	Items []MenuItemResponse `json:"items"`
}
