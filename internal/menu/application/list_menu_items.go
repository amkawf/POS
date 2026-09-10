package application

import (
	"context"

	"github.com/google/uuid"

	"pos-backend/internal/menu/domain"
	"pos-backend/internal/menu/repository"
)

type StoreStockReaer interface {
	GetStockByStore(ctx context.Context, storeID uuid.UUID) (map[uuid.UUID]int64, error)
}

type ListMenuItemsUseCase struct {
	menuItemRepository repository.MenuItemRepository
	stockReader StoreStockReaer
}

func NewListMenuItemsUseCase(
	menuItemRepository repository.MenuItemRepository,
	stockReader StoreStockReaer,
) *ListMenuItemsUseCase {
	return &ListMenuItemsUseCase{
		menuItemRepository: menuItemRepository,
		stockReader: stockReader,
	}
}

func (uc *ListMenuItemsUseCase) Execute(
	ctx context.Context,
	companyID uuid.UUID,
	storeID *uuid.UUID, // <-- Sekarang menerima storeID opsional
) ([]domain.MenuItem, error) {
	items, err := uc.menuItemRepository.ListActiveByCompany(ctx, companyID)
	if err != nil {
		return nil, err
	}
	// Jika kasir menyertakan storeID, tempelkan sisa stok ke masing-masing menu
	if storeID != nil && uc.stockReader != nil {
		stockMap, err := uc.stockReader.GetStockByStore(ctx, *storeID)
		if err != nil {
			return nil, err
		}
		for i := range items {
			val := stockMap[items[i].ID] // jika menu belum terdaftar di stok, default otomatis 0
			items[i].Stock = &val
		}
	}
	return items, nil
}
