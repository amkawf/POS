package repository

import (
	"context"

	"github.com/google/uuid"

	"pos-backend/internal/menu/domain"
)

type MenuItemRepository interface {
	ListActiveByCompany(ctx context.Context, companyID uuid.UUID) ([]domain.MenuItem, error)
}
