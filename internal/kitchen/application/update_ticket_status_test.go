package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"pos-backend/internal/kitchen/application"
	"pos-backend/internal/kitchen/domain"
)

type mockKitchenRepo struct {
	ticket *domain.KitchenTicket
	err    error
}

func (m *mockKitchenRepo) Create(ctx context.Context, ticket *domain.KitchenTicket) error {
	m.ticket = ticket
	return m.err
}

func (m *mockKitchenRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.KitchenTicket, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.ticket, nil
}

func (m *mockKitchenRepo) GetByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.KitchenTicket, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.ticket, nil
}

func (m *mockKitchenRepo) ListByStore(ctx context.Context, storeID uuid.UUID, status *string) ([]domain.KitchenTicket, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.ticket != nil {
		return []domain.KitchenTicket{*m.ticket}, nil
	}
	return []domain.KitchenTicket{}, nil
}

func (m *mockKitchenRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.TicketStatus) (*domain.KitchenTicket, error) {
	if m.err != nil {
		return nil, m.err
	}
	m.ticket.Status = status
	return m.ticket, nil
}

func TestUpdateTicketStatusUseCase_Success(t *testing.T) {
	companyID := uuid.New()
	storeID := uuid.New()
	ticketID := uuid.New()

	ticket := &domain.KitchenTicket{
		ID:          ticketID,
		CompanyID:   companyID,
		StoreID:     storeID,
		OrderID:     uuid.New(),
		OrderNumber: "ORD-001",
		OrderType:   "DINE_IN",
		Status:      domain.TicketStatusPending,
		Priority:    domain.TicketPriorityNormal,
		Items:       []domain.KitchenTicketItem{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	repo := &mockKitchenRepo{ticket: ticket}
	useCase := application.NewUpdateTicketStatusUseCase(repo)

	updated, err := useCase.Execute(context.Background(), application.UpdateTicketStatusInput{
		CompanyID: companyID,
		StoreID:   storeID,
		TicketID:  ticketID,
		Status:    domain.TicketStatusPreparing,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Status != domain.TicketStatusPreparing {
		t.Fatalf("expected status PREPARING, got %s", updated.Status)
	}
}

func TestUpdateTicketStatusUseCase_InvalidTransition(t *testing.T) {
	companyID := uuid.New()
	storeID := uuid.New()
	ticketID := uuid.New()

	ticket := &domain.KitchenTicket{
		ID:          ticketID,
		CompanyID:   companyID,
		StoreID:     storeID,
		OrderID:     uuid.New(),
		OrderNumber: "ORD-002",
		OrderType:   "DINE_IN",
		Status:      domain.TicketStatusPending,
		Priority:    domain.TicketPriorityNormal,
		Items:       []domain.KitchenTicketItem{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	repo := &mockKitchenRepo{ticket: ticket}
	useCase := application.NewUpdateTicketStatusUseCase(repo)

	// Direct transition from PENDING to SERVED is invalid
	_, err := useCase.Execute(context.Background(), application.UpdateTicketStatusInput{
		CompanyID: companyID,
		StoreID:   storeID,
		TicketID:  ticketID,
		Status:    domain.TicketStatusServed,
	})

	if err == nil {
		t.Fatalf("expected error for invalid transition, got nil")
	}
}
