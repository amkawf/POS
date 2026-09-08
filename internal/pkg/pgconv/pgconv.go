package pgconv

import (
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// UUID converts a uuid.UUID to a valid pgtype.UUID.
func UUID(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: u, Valid: true}
}

// OptUUID converts a nullable *uuid.UUID to pgtype.UUID.
func OptUUID(u *uuid.UUID) pgtype.UUID {
	if u == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *u, Valid: true}
}

// ToUUID converts pgtype.UUID to uuid.UUID.
func ToUUID(p pgtype.UUID) uuid.UUID {
	return uuid.UUID(p.Bytes)
}

// ToOptUUID converts pgtype.UUID to nullable *uuid.UUID.
func ToOptUUID(p pgtype.UUID) *uuid.UUID {
	if !p.Valid {
		return nil
	}
	u := uuid.UUID(p.Bytes)
	return &u
}

// Text converts a string to pgtype.Text.
func Text(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: true}
}

// OptText converts a *string to pgtype.Text.
func OptText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

// ToText converts pgtype.Text to string.
func ToText(p pgtype.Text) string {
	if !p.Valid {
		return ""
	}
	return p.String
}

// ToOptText converts pgtype.Text to *string.
func ToOptText(p pgtype.Text) *string {
	if !p.Valid {
		return nil
	}
	return &p.String
}

// Timestamptz converts time.Time to pgtype.Timestamptz.
func Timestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

// OptTimestamptz converts *time.Time to pgtype.Timestamptz.
func OptTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

// ToTime converts pgtype.Timestamptz to time.Time.
func ToTime(p pgtype.Timestamptz) time.Time {
	return p.Time
}

// ToOptTime converts pgtype.Timestamptz to *time.Time.
func ToOptTime(p pgtype.Timestamptz) *time.Time {
	if !p.Valid {
		return nil
	}
	return &p.Time
}

// Int64ToNumeric converts int64 to pgtype.Numeric.
func Int64ToNumeric(value int64) pgtype.Numeric {
	return pgtype.Numeric{
		Int:              big.NewInt(value),
		Exp:              0,
		NaN:              false,
		InfinityModifier: pgtype.Finite,
		Valid:            true,
	}
}

// NumericToInt64 converts pgtype.Numeric to int64, handling scale and exponential representation.
func NumericToInt64(value pgtype.Numeric) (int64, error) {
	if !value.Valid {
		return 0, nil
	}
	if value.NaN || value.InfinityModifier != pgtype.Finite {
		return 0, fmt.Errorf("numeric value is not finite")
	}
	if value.Int == nil {
		return 0, nil
	}

	result := new(big.Int).Set(value.Int)
	switch {
	case value.Exp > 0:
		scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(value.Exp)), nil)
		result.Mul(result, scale)
	case value.Exp < 0:
		scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(-value.Exp)), nil)
		quotient, remainder := new(big.Int).QuoRem(result, scale, new(big.Int))
		if remainder.Sign() != 0 {
			return 0, fmt.Errorf("numeric value has a fractional component: %s", value.Int.String())
		}
		result = quotient
	}

	if !result.IsInt64() {
		return 0, fmt.Errorf("numeric value overflows int64")
	}

	return result.Int64(), nil
}

// NumericToInt64Safe returns 0 on error, useful for non-critical reads like ticket quantities.
func NumericToInt64Safe(value pgtype.Numeric) int64 {
	val, err := NumericToInt64(value)
	if err != nil {
		return 0
	}
	return val
}
