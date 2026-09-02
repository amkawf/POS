package repository

import (
	"context"

	"pos-backend/internal/order/domain"
)

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
}
