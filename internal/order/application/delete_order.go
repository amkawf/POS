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
	orderRepo    repository.OrderRepository
	tableUpdater TableStatusUpdater
}

func NewDeleteOrderUseCase(
	orderRepo repository.OrderRepository,
	tableUpdater TableStatusUpdater,
) *DeleteOrderUseCase {
	return &DeleteOrderUseCase{
		orderRepo:    orderRepo,
		tableUpdater: tableUpdater,
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

	if err := uc.orderRepo.Delete(ctx, input.CompanyID, input.StoreID, input.OrderID); err != nil {
		return err
	}

	if uc.tableUpdater != nil && order.TableID != nil {
		_ = uc.tableUpdater.UpdateTableStatus(ctx, order.CompanyID, order.StoreID, *order.TableID, "AVAILABLE")
	}

	return nil
}
