package repository

import "testing"

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
