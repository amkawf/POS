package application_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"pos-backend/internal/menu/application"
	"pos-backend/internal/menu/domain"
)

type mockCategoryRepo struct {
	categories []domain.Category
	err        error
}

func (m *mockCategoryRepo) ListActiveByCompany(ctx context.Context, companyID uuid.UUID) ([]domain.Category, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.categories, nil
}

func (m *mockCategoryRepo) ListCategoryIDsByMenuItem(ctx context.Context, companyID uuid.UUID) (map[uuid.UUID][]uuid.UUID, error) {
	return nil, nil
}

func TestListCategoriesUseCase(t *testing.T) {
	companyID := uuid.New()
	menuID := uuid.New()

	expected := []domain.Category{
		{
			ID:        uuid.New(),
			MenuID:    menuID,
			Name:      "Makanan",
			SortOrder: 1,
			Status:    "ACTIVE",
		},
		{
			ID:        uuid.New(),
			MenuID:    menuID,
			Name:      "Minuman",
			SortOrder: 2,
			Status:    "ACTIVE",
		},
	}

	repo := &mockCategoryRepo{categories: expected}
	uc := application.NewListCategoriesUseCase(repo)

	result, err := uc.Execute(context.Background(), companyID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != len(expected) {
		t.Fatalf("expected %d categories, got %d", len(expected), len(result))
	}

	for i, c := range result {
		if c.Name != expected[i].Name {
			t.Errorf("category %d: expected name %q, got %q", i, expected[i].Name, c.Name)
		}
	}
}
