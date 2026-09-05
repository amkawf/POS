package repository

import (
	"context"

	"github.com/google/uuid"

	"pos-backend/internal/menu/domain"
)

type CategoryRepository interface {
	ListActiveByCompany(ctx context.Context, companyID uuid.UUID) ([]domain.Category, error)
	ListCategoryIDsByMenuItem(ctx context.Context, companyID uuid.UUID) (map[uuid.UUID][]uuid.UUID, error)
}
