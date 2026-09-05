package application

import (
	"context"

	"github.com/google/uuid"

	"pos-backend/internal/kitchen/domain"
	"pos-backend/internal/kitchen/repository"
)

type TicketItemInput struct {
	OrderItemID *uuid.UUID
	MenuItemID  uuid.UUID
	ItemName    string
	SKU         string
	Quantity    int64
	Notes       *string
}

type CreateTicketInput struct {
	CompanyID   uuid.UUID
	StoreID     uuid.UUID
	OrderID     uuid.UUID
	OrderNumber string
	OrderType   string
	TableID     *uuid.UUID
	Priority    domain.TicketPriority
	Notes       *string
	Items       []TicketItemInput
}

type CreateTicketUseCase struct {
	repo repository.KitchenRepository
}

func NewCreateTicketUseCase(repo repository.KitchenRepository) *CreateTicketUseCase {
	return &CreateTicketUseCase{repo: repo}
}

func (uc *CreateTicketUseCase) Execute(
	ctx context.Context,
	input CreateTicketInput,
) (*domain.KitchenTicket, error) {
	ticket := domain.NewKitchenTicket(
		input.CompanyID,
		input.StoreID,
		input.OrderID,
		input.OrderNumber,
		input.OrderType,
		input.TableID,
		input.Priority,
		input.Notes,
	)

	for _, item := range input.Items {
		ticket.AddItem(
			item.OrderItemID,
			item.MenuItemID,
			item.ItemName,
			item.SKU,
			item.Quantity,
			item.Notes,
		)
	}

	if err := uc.repo.Create(ctx, ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}
