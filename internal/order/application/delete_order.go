package application

import (
	"context"

	"github.com/google/uuid"

	"pos-backend/internal/order/repository"
)

type DeleteOrderInput struct {
	CompanyID uuid.UUID
	StoreID   uuid.UUID
	OrderID   uuid.UUID
}

type DeleteOrderUseCase struct {
	orderRepo repository.OrderRepository
}

func NewDeleteOrderUseCase(orderRepo repository.OrderRepository) *DeleteOrderUseCase {
	return &DeleteOrderUseCase{
		orderRepo: orderRepo,
	}
}

func (uc *DeleteOrderUseCase) Execute(
	ctx context.Context,
	input DeleteOrderInput,
) error {
	order, err := uc.orderRepo.FindByID(ctx, input.CompanyID, input.StoreID, input.OrderID)
	if err != nil {
		return err
	}
	if order == nil {
		return ErrOrderNotFound
	}

	if err := order.CanDelete(); err != nil {
		return err
	}

	return uc.orderRepo.Delete(ctx, input.CompanyID, input.StoreID, input.OrderID)
}
