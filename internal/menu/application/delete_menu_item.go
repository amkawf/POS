package application

import (
	"context"

	"github.com/google/uuid"

	"pos-backend/internal/menu/repository"
)

type DeleteMenuItemUseCase struct {
	repo repository.MenuItemRepository
}

func NewDeleteMenuItemUseCase(repo repository.MenuItemRepository) *DeleteMenuItemUseCase {
	return &DeleteMenuItemUseCase{repo: repo}
}

func (uc *DeleteMenuItemUseCase) Execute(ctx context.Context, companyID, itemID uuid.UUID) error {
	return uc.repo.Delete(ctx, companyID, itemID)
}