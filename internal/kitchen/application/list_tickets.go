package application

import (
	"context"

	"github.com/google/uuid"

	"pos-backend/internal/kitchen/domain"
	"pos-backend/internal/kitchen/repository"
)

type ListTicketsInput struct {
	CompanyID uuid.UUID
	StoreID   uuid.UUID
	Status    *string
}

type ListTicketsUseCase struct {
	repo repository.KitchenRepository
}

func NewListTicketsUseCase(repo repository.KitchenRepository) *ListTicketsUseCase {
	return &ListTicketsUseCase{repo: repo}
}

func (uc *ListTicketsUseCase) Execute(
	ctx context.Context,
	input ListTicketsInput,
) ([]domain.KitchenTicket, error) {
	return uc.repo.ListByStore(ctx, input.StoreID, input.Status)
}
