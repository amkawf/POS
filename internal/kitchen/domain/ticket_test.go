package domain_test

import (
	"testing"

	"github.com/google/uuid"

	"pos-backend/internal/kitchen/domain"
)

func TestKitchenTicket_Lifecycle(t *testing.T) {
	companyID := uuid.New()
	storeID := uuid.New()
	orderID := uuid.New()

	ticket := domain.NewKitchenTicket(
		companyID,
		storeID,
		orderID,
		"ORD-12345",
		"DINE_IN",
		nil,
		domain.TicketPriorityNormal,
		nil,
	)

	menuItemID := uuid.New()
	ticket.AddItem(nil, menuItemID, "Nasi Goreng", "FD-001", 2, nil)

	if ticket.Status != domain.TicketStatusPending {
		t.Fatalf("expected PENDING, got %s", ticket.Status)
	}

	// PENDING -> PREPARING
	if err := ticket.TransitionTo(domain.TicketStatusPreparing); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ticket.Status != domain.TicketStatusPreparing {
		t.Fatalf("expected PREPARING, got %s", ticket.Status)
	}
	if ticket.StartedAt == nil {
		t.Fatalf("expected StartedAt to be set")
	}

	// PREPARING -> READY
	if err := ticket.TransitionTo(domain.TicketStatusReady); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ticket.Status != domain.TicketStatusReady {
		t.Fatalf("expected READY, got %s", ticket.Status)
	}
	if ticket.ReadyAt == nil {
		t.Fatalf("expected ReadyAt to be set")
	}

	// READY -> SERVED
	if err := ticket.TransitionTo(domain.TicketStatusServed); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ticket.Status != domain.TicketStatusServed {
		t.Fatalf("expected SERVED, got %s", ticket.Status)
	}
	if ticket.ServedAt == nil {
		t.Fatalf("expected ServedAt to be set")
	}

	// SERVED is final: cannot transition again
	if err := ticket.TransitionTo(domain.TicketStatusPreparing); err != domain.ErrTicketAlreadyFinalized {
		t.Fatalf("expected ErrTicketAlreadyFinalized, got %v", err)
	}
}
