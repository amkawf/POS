package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"pos-backend/internal/order/domain"
	"pos-backend/internal/order/repository"
)

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrInsufficientAmount = errors.New("amount paid is less than total amount")
)

type PayOrderInput struct {
	CompanyID     uuid.UUID
	StoreID       uuid.UUID
	OrderID       uuid.UUID
	PaymentMethod string
	AmountPaid    int64
	Notes         *string
}

type PayOrderUseCase struct {
	orderRepo repository.OrderRepository
}

func NewPayOrderUseCase(orderRepo repository.OrderRepository) *PayOrderUseCase {
	return &PayOrderUseCase{
		orderRepo: orderRepo,
	}
}

func (uc *PayOrderUseCase) Execute(
	ctx context.Context,
	input PayOrderInput,
) (*domain.Order, error) {
	order, err := uc.orderRepo.FindByID(ctx, input.CompanyID, input.StoreID, input.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	if input.AmountPaid < order.TotalAmount {
		return nil, ErrInsufficientAmount
	}

	change := input.AmountPaid - order.TotalAmount
	paymentNote := fmt.Sprintf("Payment: %s | Tendered: %d | Change: %d", input.PaymentMethod, input.AmountPaid, change)
	if input.Notes != nil && *input.Notes != "" {
		paymentNote = fmt.Sprintf("%s | %s", paymentNote, *input.Notes)
	}

	if err := order.Complete(&paymentNote); err != nil {
		return nil, err
	}

	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}
