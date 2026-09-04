package repository

import (
	"context"

	"github.com/google/uuid"

	"pos-backend/internal/order/domain"
)

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	FindByID(ctx context.Context, companyID, storeID, orderID uuid.UUID) (*domain.Order, error)
	ListByStore(ctx context.Context, companyID, storeID uuid.UUID, limit int32) ([]*domain.Order, error)
	ListByStatus(ctx context.Context, companyID, storeID uuid.UUID, status string, limit int32) ([]*domain.Order, error)
	Update(ctx context.Context, order *domain.Order) error
	Delete(ctx context.Context, companyID, storeID, orderID uuid.UUID) error
}
