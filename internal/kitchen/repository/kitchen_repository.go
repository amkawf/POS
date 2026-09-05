package repository

import (
	"context"

	"github.com/google/uuid"

	"pos-backend/internal/kitchen/domain"
)

type KitchenRepository interface {
	Create(ctx context.Context, ticket *domain.KitchenTicket) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.KitchenTicket, error)
	GetByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.KitchenTicket, error)
	ListByStore(ctx context.Context, storeID uuid.UUID, status *string) ([]domain.KitchenTicket, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.TicketStatus) (*domain.KitchenTicket, error)
}
