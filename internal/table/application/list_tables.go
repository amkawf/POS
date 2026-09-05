package application

import (
	"context"

	"github.com/google/uuid"

	"pos-backend/internal/table/domain"
	"pos-backend/internal/table/repository"
)

type ListTablesUseCase struct {
	tableRepository repository.TableRepository
}

func NewListTablesUseCase(
	tableRepository repository.TableRepository,
) *ListTablesUseCase {
	return &ListTablesUseCase{
		tableRepository: tableRepository,
	}
}

func (uc *ListTablesUseCase) Execute(
	ctx context.Context,
	companyID, storeID uuid.UUID,
) ([]domain.Table, error) {
	return uc.tableRepository.ListActiveByStore(ctx, companyID, storeID)
}
