package repository

import (
	"context"

	"github.com/google/uuid"

	"pos-backend/internal/table/domain"
)

type TableRepository interface {
	ListActiveByStore(ctx context.Context, companyID, storeID uuid.UUID) ([]domain.Table, error)
	GetByID(ctx context.Context, companyID, storeID, id uuid.UUID) (*domain.Table, error)
	UpdateStatus(ctx context.Context, companyID, storeID, id uuid.UUID, status domain.TableStatus) error
}
