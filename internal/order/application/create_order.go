package application

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"pos-backend/internal/order/domain"
	"pos-backend/internal/order/repository"
)

type CreateOrderInput struct {
	CompanyID uuid.UUID
	StoreID   uuid.UUID
	TableID   *uuid.UUID

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
	Notes *string
}

type KitchenTicketCreator interface {
	CreateTicket(
		ctx context.Context,
		companyID, storeID, orderID uuid.UUID,
		orderNumber, orderType string,
		tableID *uuid.UUID,
		items []CreateOrderItemInput,
	) error
}

type TableStatusUpdater interface {
	UpdateTableStatus(ctx context.Context, companyID, storeID, tableID uuid.UUID, status string) error
}

type CreateOrderUseCase struct {
	orderRepository repository.OrderRepository
	kitchenCreator  KitchenTicketCreator
	tableUpdater    TableStatusUpdater
}

func NewCreateOrderUseCase(
	orderRepository repository.OrderRepository,
	kitchenCreator KitchenTicketCreator,
	tableUpdater TableStatusUpdater,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderRepository: orderRepository,
		kitchenCreator:  kitchenCreator,
		tableUpdater:    tableUpdater,
	}
}

func (uc *CreateOrderUseCase) Execute(
	ctx context.Context,
	input CreateOrderInput,
) (*domain.Order, error) {
	order, err := domain.NewOrder(
		input.CompanyID,
		input.StoreID,
		generateOrderNumber(),
		input.OrderType,
		input.OrderSource,
		input.CreatedBy,
	)
	if err != nil {
		return nil, err
	}

	order.CustomerName = input.CustomerName
	order.Notes = input.Notes
	order.TableID = input.TableID

	for _, item := range input.Items {
		if err := order.AddItem(
			item.MenuItemID,
			item.ItemName,
			item.SKU,
			item.Quantity,
			item.UnitPrice,
			item.Notes,
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

	// If order is bound to a table, mark the table as OCCUPIED
	if uc.tableUpdater != nil && order.TableID != nil {
		_ = uc.tableUpdater.UpdateTableStatus(ctx, order.CompanyID, order.StoreID, *order.TableID, "OCCUPIED")
	}

	if uc.kitchenCreator != nil {
		_ = uc.kitchenCreator.CreateTicket(
			ctx,
			order.CompanyID,
			order.StoreID,
			order.ID,
			order.OrderNumber,
			string(order.OrderType),
			order.TableID,
			input.Items,
		)
	}

	return order, nil
}

// generateOrderNumber produces a human-scannable, collision-resistant order
// number without a per-store sequence/lock. A true sequential daily counter
// (e.g. ORD-20260903-001) needs a locking strategy that is not yet decided
// (see CLAUDE.md deferred: detailed transaction isolation/locking).
func generateOrderNumber() string {
	suffix := strings.ToUpper(strings.ReplaceAll(uuid.New().String(), "-", "")[:8])
	return fmt.Sprintf("ORD-%s-%s", time.Now().UTC().Format("20060102"), suffix)
}
