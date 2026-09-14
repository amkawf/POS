package application

import (
	"context"

	"github.com/google/uuid"

	"pos-backend/internal/menu/domain"
	"pos-backend/internal/menu/repository"
)

// CreateMenuItemUseCase menangani pembuatan menu makanan/minuman baru
type CreateMenuItemUseCase struct {
	menuItemRepo repository.MenuItemRepository
}

func NewCreateMenuItemUseCase(menuItemRepo repository.MenuItemRepository) *CreateMenuItemUseCase {
	return &CreateMenuItemUseCase{menuItemRepo: menuItemRepo}
}

func (uc *CreateMenuItemUseCase) Execute(
	ctx context.Context,
	companyID uuid.UUID,
	sku string,
	name string,
	description *string,
	basePrice int64,
	categoryIDs []uuid.UUID,
	fulfillmentType string,
) (*domain.MenuItem, error) {
	// Jika fulfillmentType tidak diisi, default adalah BATCH_COOKING
	if fulfillmentType == "" {
		fulfillmentType = "BATCH_COOKING"
	}

	item := &domain.MenuItem{
		CompanyID:       companyID,
		SKU:             sku,
		Name:            name,
		Description:     description,
		BasePrice:       basePrice,
		CategoryIDs:     categoryIDs,
		FulfillmentType: fulfillmentType,
	}

	if err := uc.menuItemRepo.Create(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}