package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"pos-backend/internal/table/application"
	"pos-backend/internal/table/domain"
)

type mockTableRepo struct {
	tables []domain.Table
	err    error
}

func (m *mockTableRepo) ListActiveByStore(ctx context.Context, companyID, storeID uuid.UUID) ([]domain.Table, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.tables, nil
}

func (m *mockTableRepo) GetByID(ctx context.Context, companyID, storeID, id uuid.UUID) (*domain.Table, error) {
	return nil, nil
}

func (m *mockTableRepo) UpdateStatus(ctx context.Context, companyID, storeID, id uuid.UUID, status domain.TableStatus) error {
	return nil
}

func TestListTablesUseCase(t *testing.T) {
	companyID := uuid.New()
	storeID := uuid.New()

	expected := []domain.Table{
		{
			ID:          uuid.New(),
			CompanyID:   companyID,
			StoreID:     storeID,
			TableNumber: "Meja 01",
			Capacity:    2,
			Status:      domain.TableStatusAvailable,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			CompanyID:   companyID,
			StoreID:     storeID,
			TableNumber: "Meja 02",
			Capacity:    4,
			Status:      domain.TableStatusOccupied,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	repo := &mockTableRepo{tables: expected}
	uc := application.NewListTablesUseCase(repo)

	result, err := uc.Execute(context.Background(), companyID, storeID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != len(expected) {
		t.Fatalf("expected %d tables, got %d", len(expected), len(result))
	}

	for i, tbl := range result {
		if tbl.TableNumber != expected[i].TableNumber {
			t.Errorf("table %d: expected number %q, got %q", i, expected[i].TableNumber, tbl.TableNumber)
		}
		if tbl.Status != expected[i].Status {
			t.Errorf("table %d: expected status %q, got %q", i, expected[i].Status, tbl.Status)
		}
	}
}
