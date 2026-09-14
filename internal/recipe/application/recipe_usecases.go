package application

import (
	"context"

	"github.com/google/uuid"

	"pos-backend/internal/recipe/domain"
	"pos-backend/internal/recipe/repository"
)

// GetRecipeUseCase mengambil daftar bahan dan takaran per porsi untuk menu tertentu
type GetRecipeUseCase struct {
	repo repository.RecipeRepository
}

func NewGetRecipeUseCase(repo repository.RecipeRepository) *GetRecipeUseCase {
	return &GetRecipeUseCase{repo: repo}
}

func (uc *GetRecipeUseCase) Execute(ctx context.Context, menuItemID uuid.UUID) ([]domain.RecipeItem, error) {
	return uc.repo.GetByMenuItemID(ctx, menuItemID)
}

// SaveRecipeUseCase menyimpan daftar bahan dan takaran baru untuk suatu menu
type SaveRecipeUseCase struct {
	repo repository.RecipeRepository
}

func NewSaveRecipeUseCase(repo repository.RecipeRepository) *SaveRecipeUseCase {
	return &SaveRecipeUseCase{repo: repo}
}

func (uc *SaveRecipeUseCase) Execute(
	ctx context.Context,
	menuItemID uuid.UUID,
	items []domain.RecipeItem,
) error {
	return uc.repo.SaveRecipe(ctx, menuItemID, items)
}