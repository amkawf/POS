package repository

import (
	"math/big"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestInt64ToNumeric(t *testing.T) {
	tests := []struct {
		name  string
		input int64
	}{
		{
			name:  "zero",
			input: 0,
		},
		{
			name:  "positive",
			input: 25000,
		},
		{
			name:  "negative",
			input: -100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			numeric := int64ToNumeric(tt.input)

			if !numeric.Valid {
				t.Fatal("expected numeric to be valid")
			}

			if numeric.Int == nil {
				t.Fatal("expected numeric Int to be set")
			}

			if numeric.Int.Int64() != tt.input {
				t.Fatalf(
					"expected %d, got %d",
					tt.input,
					numeric.Int.Int64(),
				)
			}
		})
	}
}

func TestNumericToInt64(t *testing.T) {
	tests := []struct {
		name  string
		input int64
	}{
		{
			name:  "zero",
			input: 0,
		},
		{
			name:  "positive",
			input: 25000,
		},
		{
			name:  "negative",
			input: -100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			numeric := int64ToNumeric(tt.input)

			result, err := numericToInt64(numeric)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if result != tt.input {
				t.Fatalf(
					"expected %d, got %d",
					tt.input,
					result,
				)
			}
		})
	}
}

func TestNumericToInt64ScaledExponent(t *testing.T) {
	// Reflects how Postgres actually sends NUMERIC(15,2) over the wire,
	// e.g. 28000.00 decoded as Int=2800000, Exp=-2 — not Exp=0 like
	// int64ToNumeric produces locally.
	numeric := pgtype.Numeric{
		Int:   big.NewInt(2800000),
		Exp:   -2,
		Valid: true,
	}

	result, err := numericToInt64(numeric)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != 28000 {
		t.Fatalf("expected 28000, got %d", result)
	}
}

func TestNumericToInt64FractionalRejected(t *testing.T) {
	numeric := pgtype.Numeric{
		Int:   big.NewInt(280005),
		Exp:   -2,
		Valid: true,
	}

	if _, err := numericToInt64(numeric); err == nil {
		t.Fatal("expected error for a value with a fractional remainder")
	}
}

func TestNumericToInt64Invalid(t *testing.T) {
	numeric := int64ToNumeric(25000)
	numeric.Valid = false

	result, err := numericToInt64(numeric)

	if err != nil {
		t.Fatalf("expected no error for invalid nullable numeric, got %v", err)
	}

	if result != 0 {
		t.Fatalf("expected 0, got %d", result)
	}
}
