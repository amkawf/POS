package repository

import (
	"fmt"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
)

func int64ToNumeric(value int64) pgtype.Numeric {
	return pgtype.Numeric{
		Int:              big.NewInt(value),
		Exp:              0,
		NaN:              false,
		InfinityModifier: pgtype.Finite,
		Valid:            true,
	}
}

func numericToInt64(value pgtype.Numeric) (int64, error) {
	if !value.Valid {
		return 0, nil
	}

	if value.NaN || value.InfinityModifier != pgtype.Finite {
		return 0, fmt.Errorf("numeric value is not finite")
	}

	if value.Exp != 0 {
		return 0, fmt.Errorf(
			"numeric value has unsupported exponent: %d",
			value.Exp,
		)
	}

	if value.Int == nil {
		return 0, nil
	}

	if !value.Int.IsInt64() {
		return 0, fmt.Errorf("numeric value overflows int64")
	}

	return value.Int.Int64(), nil
}
