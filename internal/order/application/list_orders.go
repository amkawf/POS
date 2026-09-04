package application

import (
	"context"

	"github.com/google/uuid"

	"pos-backend/internal/order/domain"
	"pos-backend/internal/order/repository"
)

type ListOrdersInput struct {
	CompanyID uuid.UUID
	StoreID   uuid.UUID
	Status    *string
	Limit     int32
}

type ListOrdersUseCase struct {
	orderRepo repository.OrderRepository
}

func NewListOrdersUseCase(orderRepo repository.OrderRepository) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		orderRepo: orderRepo,
	}
}

func (uc *ListOrdersUseCase) Execute(
	ctx context.Context,
	input ListOrdersInput,
) ([]*domain.Order, error) {
	if input.Status != nil && *input.Status != "" {
		return uc.orderRepo.ListByStatus(
			ctx,
			input.CompanyID,
			input.StoreID,
			*input.Status,
			input.Limit,
		)
	}

	return uc.orderRepo.ListByStore(
		ctx,
		input.CompanyID,
		input.StoreID,
		input.Limit,
	)
}
