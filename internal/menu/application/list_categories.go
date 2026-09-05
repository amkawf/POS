package application

import (
	"context"

	"github.com/google/uuid"

	"pos-backend/internal/menu/domain"
	"pos-backend/internal/menu/repository"
)

type ListCategoriesUseCase struct {
	categoryRepository repository.CategoryRepository
}

func NewListCategoriesUseCase(
	categoryRepository repository.CategoryRepository,
) *ListCategoriesUseCase {
	return &ListCategoriesUseCase{
		categoryRepository: categoryRepository,
	}
}

func (uc *ListCategoriesUseCase) Execute(
	ctx context.Context,
	companyID uuid.UUID,
) ([]domain.Category, error) {
	return uc.categoryRepository.ListActiveByCompany(ctx, companyID)
}
