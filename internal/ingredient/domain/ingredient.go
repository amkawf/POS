package domain

import (
	"time"

	"github.com/google/uuid"
)

type Ingredient struct {
	ID             uuid.UUID `json:"id"`
	CompanyID      uuid.UUID `json:"company_id"`
	Code           *string   `json:"code,omitempty"`
	Name           string    `json:"name"`
	Unit           string    `json:"unit"` // kg, gram, liter, butir, dll
	MinStockAlert  float64   `json:"min_stock_alert"`
	Stock          *float64  `json:"stock,omitempty"` // Nilai stok di toko aktif
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type IngredientMovement struct {
	ID           uuid.UUID  `json:"id"`
	CompanyID    uuid.UUID  `json:"company_id"`
	StoreID      uuid.UUID  `json:"store_id"`
	IngredientID uuid.UUID  `json:"ingredient_id"`
	Quantity     float64    `json:"quantity"`
	Type         string     `json:"type"` // PURCHASE, COOKING_USAGE, WASTE, ADJUSTMENT
	ReferenceID  *uuid.UUID `json:"reference_id,omitempty"`
	Notes        *string    `json:"notes,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}