package recipehttp

import (
	"github.com/google/uuid"
	"pos-backend/internal/recipe/domain"
)

// RecipeItemPayload adalah baris input takaran bahan saat admin menyimpan resep di web
type RecipeItemPayload struct {
	IngredientID       uuid.UUID `json:"ingredient_id" binding:"required"`       // ID bahan baku
	QuantityPerPortion float64   `json:"quantity_per_portion" binding:"required"` // Takaran per 1 porsi (misal: 0.25 atau 5)
}

// SaveRecipeRequest membungkus array takaran bahan untuk satu menu
type SaveRecipeRequest struct {
	Items []RecipeItemPayload `json:"items"`
}

// RecipeResponse mengirimkan detail resep lengkap kembali ke frontend
type RecipeResponse struct {
	MenuItemID uuid.UUID           `json:"menu_item_id"`
	Items      []domain.RecipeItem `json:"items"`
}