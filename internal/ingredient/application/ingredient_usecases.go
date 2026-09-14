package application

import (
	"context"

	"github.com/google/uuid"

	"pos-backend/internal/ingredient/domain"
	"pos-backend/internal/ingredient/repository"
)

type CreateIngredientUseCase struct {
	repo repository.IngredientRepository
}

func NewCreateIngredientUseCase(repo repository.IngredientRepository) *CreateIngredientUseCase {
	return &CreateIngredientUseCase{repo: repo}
}

func (uc *CreateIngredientUseCase) Execute(
	ctx context.Context,
	companyID uuid.UUID,
	code *string,
	name string,
	unit string,
	minAlert float64,
) (*domain.Ingredient, error) {
	ing := &domain.Ingredient{
		CompanyID:     companyID,
		Code:          code,
		Name:          name,
		Unit:          unit,
		MinStockAlert: minAlert,
	}
	if err := uc.repo.Create(ctx, ing); err != nil {
		return nil, err
	}
	return ing, nil
}

type ListIngredientsUseCase struct {
	repo repository.IngredientRepository
}

func NewListIngredientsUseCase(repo repository.IngredientRepository) *ListIngredientsUseCase {
	return &ListIngredientsUseCase{repo: repo}
}

func (uc *ListIngredientsUseCase) Execute(
	ctx context.Context,
	companyID, storeID uuid.UUID,
) ([]domain.Ingredient, error) {
	return uc.repo.ListByStore(ctx, companyID, storeID)
}

type RestockIngredientUseCase struct {
	repo repository.IngredientRepository
}

func NewRestockIngredientUseCase(repo repository.IngredientRepository) *RestockIngredientUseCase {
	return &RestockIngredientUseCase{repo: repo}
}

func (uc *RestockIngredientUseCase) Execute(
	ctx context.Context,
	companyID, storeID, ingredientID uuid.UUID,
	quantity float64,
	notes string,
) error {
	return uc.repo.Restock(ctx, companyID, storeID, ingredientID, quantity, notes)
}

// UpdateIngredientUseCase menangani pembaruan atribut bahan baku
type UpdateIngredientUseCase struct {
	repo repository.IngredientRepository
}

func NewUpdateIngredientUseCase(repo repository.IngredientRepository) *UpdateIngredientUseCase {
	return &UpdateIngredientUseCase{repo: repo}
}

func (uc *UpdateIngredientUseCase) Execute(
	ctx context.Context,
	companyID, ingredientID uuid.UUID,
	code *string,
	name string,
	unit string,
	minAlert float64,
) error {
	ing := &domain.Ingredient{
		ID:            ingredientID,
		CompanyID:     companyID,
		Code:          code,
		Name:          name,
		Unit:          unit,
		MinStockAlert: minAlert,
	}
	return uc.repo.Update(ctx, ing)
}