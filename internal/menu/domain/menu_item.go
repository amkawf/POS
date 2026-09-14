package domain

import "github.com/google/uuid"

type MenuItem struct {
	ID          uuid.UUID
	CompanyID   uuid.UUID
	SKU         string
	Name        string
	Description *string
	BasePrice   int64
	Status      string
	CategoryIDs []uuid.UUID
	Stock		*int64
	FulfillmentType string      `json:"fulfillment_type"` // "BATCH_COOKING" atau "MADE_TO_ORDER"
}
