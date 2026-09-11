package application

import (
	"context"

	"github.com/google/uuid"

	"pos-backend/internal/inventory/repository"
)

type AdjustStockUseCase struct {
	inventoryRepo repository.InventoryRepository
}

func NewAdjustStockUseCase(inventoryRepo repository.InventoryRepository) *AdjustStockUseCase {
	return &AdjustStockUseCase{
		inventoryRepo: inventoryRepo,
	}
}

func (uc *AdjustStockUseCase) Execute(
	ctx context.Context,
	companyID, storeID, menuItemID uuid.UUID,
	quantity int64,
	notes string,
) error {
	return uc.inventoryRepo.AdjustStock(ctx, companyID, storeID, menuItemID, quantity, notes)
}