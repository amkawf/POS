package application

import (
	"context"

	"github.com/google/uuid"

	"pos-backend/internal/kitchen/domain"
	"pos-backend/internal/kitchen/repository"
)

type UpdateTicketStatusInput struct {
	CompanyID uuid.UUID
	StoreID   uuid.UUID
	TicketID  uuid.UUID
	Status    domain.TicketStatus
}

type UpdateTicketStatusUseCase struct {
	repo repository.KitchenRepository
}

func NewUpdateTicketStatusUseCase(repo repository.KitchenRepository) *UpdateTicketStatusUseCase {
	return &UpdateTicketStatusUseCase{repo: repo}
}

func (uc *UpdateTicketStatusUseCase) Execute(
	ctx context.Context,
	input UpdateTicketStatusInput,
) (*domain.KitchenTicket, error) {
	ticket, err := uc.repo.GetByID(ctx, input.TicketID)
	if err != nil {
		return nil, err
	}

	if err := ticket.TransitionTo(input.Status); err != nil {
		return nil, err
	}

	return uc.repo.UpdateStatus(ctx, input.TicketID, input.Status)
}
