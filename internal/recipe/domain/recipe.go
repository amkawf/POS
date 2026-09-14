package domain

import "github.com/google/uuid"

// RecipeItem memodelkan satu baris takaran bahan baku untuk sebuah menu makanan
type RecipeItem struct {
	ID                 uuid.UUID `json:"id"`
	MenuItemID         uuid.UUID `json:"menu_item_id"`
	IngredientID       uuid.UUID `json:"ingredient_id"`
	IngredientName     string    `json:"ingredient_name,omitempty"`     // Nama bahan baku (misal: "Daging Ayam")
	IngredientUnit     string    `json:"ingredient_unit,omitempty"`     // Satuan bahan baku (misal: "kg" atau "gram")
	QuantityPerPortion float64   `json:"quantity_per_portion"`          // Takaran untuk 1 porsi (misal: 0.25 atau 5)
}