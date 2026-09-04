package application

import (
	"context"

	"github.com/google/uuid"

	"pos-backend/internal/menu/domain"
	"pos-backend/internal/menu/repository"
)

type ListMenuItemsUseCase struct {
	menuItemRepository repository.MenuItemRepository
}

func NewListMenuItemsUseCase(
	menuItemRepository repository.MenuItemRepository,
) *ListMenuItemsUseCase {
	return &ListMenuItemsUseCase{
		menuItemRepository: menuItemRepository,
	}
}

func (uc *ListMenuItemsUseCase) Execute(
	ctx context.Context,
	companyID uuid.UUID,
) ([]domain.MenuItem, error) {
	return uc.menuItemRepository.ListActiveByCompany(ctx, companyID)
}
