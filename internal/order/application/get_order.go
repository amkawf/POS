package application

import (
	"context"

	"github.com/google/uuid"

	"pos-backend/internal/order/domain"
	"pos-backend/internal/order/repository"
)

type GetOrderInput struct {
	CompanyID uuid.UUID
	StoreID   uuid.UUID
	OrderID   uuid.UUID
}

type GetOrderUseCase struct {
	orderRepo repository.OrderRepository
}

func NewGetOrderUseCase(orderRepo repository.OrderRepository) *GetOrderUseCase {
	return &GetOrderUseCase{
		orderRepo: orderRepo,
	}
}

func (uc *GetOrderUseCase) Execute(
	ctx context.Context,
	input GetOrderInput,
) (*domain.Order, error) {
	return uc.orderRepo.FindByID(
		ctx,
		input.CompanyID,
		input.StoreID,
		input.OrderID,
	)
}
