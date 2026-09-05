package menuhttp

import "github.com/google/uuid"

type MenuItemResponse struct {
	ID          uuid.UUID   `json:"id"`
	SKU         string      `json:"sku"`
	Name        string      `json:"name"`
	Description *string     `json:"description,omitempty"`
	BasePrice   int64       `json:"base_price"`
	CategoryIDs []uuid.UUID `json:"category_ids"`
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
