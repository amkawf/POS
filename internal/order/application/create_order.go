package application

import (
	"context"

	"github.com/google/uuid"

	"pos-backend/internal/order/domain"
	"pos-backend/internal/order/repository"
)

type CreateOrderInput struct {
	CompanyID uuid.UUID
	StoreID   uuid.UUID

	OrderNumber string
	OrderType   domain.OrderType
	OrderSource domain.OrderSource

	CustomerName *string
	Notes        *string

	CreatedBy *uuid.UUID

	Items []CreateOrderItemInput
}

type CreateOrderItemInput struct {
	MenuItemID uuid.UUID
	ItemName   string
	SKU        string

	Quantity  int64
	UnitPrice int64
}

type CreateOrderUseCase struct {
	orderRepository repository.OrderRepository
}

func NewCreateOrderUseCase(
	orderRepository repository.OrderRepository,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderRepository: orderRepository,
	}
}

func (uc *CreateOrderUseCase) Execute(
	ctx context.Context,
	input CreateOrderInput,
) (*domain.Order, error) {
	order, err := domain.NewOrder(
		input.CompanyID,
		input.StoreID,
		input.OrderNumber,
		input.OrderType,
		input.OrderSource,
		input.CreatedBy,
	)
	if err != nil {
		return nil, err
	}

	order.CustomerName = input.CustomerName
	order.Notes = input.Notes

	for _, item := range input.Items {
		if err := order.AddItem(
			item.MenuItemID,
			item.ItemName,
			item.SKU,
			item.Quantity,
			item.UnitPrice,
		); err != nil {
			return nil, err
		}
	}

	if err := order.Open(); err != nil {
		return nil, err
	}

	if err := uc.orderRepository.Create(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}
