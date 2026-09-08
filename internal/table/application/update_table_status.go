package application

import (
	"context"

	"github.com/google/uuid"

	"pos-backend/internal/table/domain"
	"pos-backend/internal/table/repository"
)

type UpdateTableStatusUseCase struct {
	tableRepo repository.TableRepository
}

func NewUpdateTableStatusUseCase(tableRepo repository.TableRepository) *UpdateTableStatusUseCase {
	return &UpdateTableStatusUseCase{
		tableRepo: tableRepo,
	}
}

func (uc *UpdateTableStatusUseCase) Execute(
	ctx context.Context,
	companyID, storeID, tableID uuid.UUID,
	status domain.TableStatus,
) error {
	return uc.tableRepo.UpdateStatus(ctx, companyID, storeID, tableID, status)
}
