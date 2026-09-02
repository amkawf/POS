package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewOrder(t *testing.T) {
	companyID := uuid.New()
	storeID := uuid.New()

	order, err := NewOrder(
		companyID,
		storeID,
		"ORD-001",
		OrderTypeDineIn,
		OrderSourcePOS,
		nil,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if order.Status != OrderStatusDraft {
		t.Fatalf(
			"expected status %s, got %s",
			OrderStatusDraft,
			order.Status,
		)
	}

	if len(order.Items) != 0 {
		t.Fatalf(
			"expected empty items, got %d",
			len(order.Items),
		)
	}
}

func TestOrderAddItem(t *testing.T) {
	order, err := NewOrder(
		uuid.New(),
		uuid.New(),
		"ORD-001",
		OrderTypeDineIn,
		OrderSourcePOS,
		nil,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = order.AddItem(
		uuid.New(),
		"Nasi Goreng",
		"FOOD-001",
		2,
		25000,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(order.Items) != 1 {
		t.Fatalf(
			"expected 1 item, got %d",
			len(order.Items),
		)
	}

	if order.Subtotal != 50000 {
		t.Fatalf(
			"expected subtotal 50000, got %d",
			order.Subtotal,
		)
	}

	if order.TotalAmount != 50000 {
		t.Fatalf(
			"expected total 50000, got %d",
			order.TotalAmount,
		)
	}
}

func TestOrderOpen(t *testing.T) {
	order, err := NewOrder(
		uuid.New(),
		uuid.New(),
		"ORD-001",
		OrderTypeDineIn,
		OrderSourcePOS,
		nil,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = order.Open()

	if err != ErrEmptyOrder {
		t.Fatalf(
			"expected ErrEmptyOrder, got %v",
			err,
		)
	}

	err = order.AddItem(
		uuid.New(),
		"Nasi Goreng",
		"FOOD-001",
		1,
		25000,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = order.Open()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if order.Status != OrderStatusOpen {
		t.Fatalf(
			"expected status %s, got %s",
			OrderStatusOpen,
			order.Status,
		)
	}
}

func TestOrderCannotOpenTwice(t *testing.T) {
	order, err := NewOrder(
		uuid.New(),
		uuid.New(),
		"ORD-001",
		OrderTypeDineIn,
		OrderSourcePOS,
		nil,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = order.AddItem(
		uuid.New(),
		"Nasi Goreng",
		"FOOD-001",
		1,
		25000,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = order.Open()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = order.Open()

	if err != ErrInvalidOrderTransition {
		t.Fatalf(
			"expected ErrInvalidOrderTransition, got %v",
			err,
		)
	}
}
