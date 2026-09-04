package repository

import (
	"fmt"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
)

func numericToInt64(value pgtype.Numeric) (int64, error) {
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
		scale := new(big.Int).Exp(
			big.NewInt(10),
			big.NewInt(int64(value.Exp)),
			nil,
		)
		result.Mul(result, scale)

	case value.Exp < 0:
		scale := new(big.Int).Exp(
			big.NewInt(10),
			big.NewInt(int64(-value.Exp)),
			nil,
		)

		quotient, remainder := new(big.Int).QuoRem(
			result,
			scale,
			new(big.Int),
		)

		if remainder.Sign() != 0 {
			return 0, fmt.Errorf(
				"numeric value has a fractional component: %s",
				value.Int.String(),
			)
		}

		result = quotient
	}

	if !result.IsInt64() {
		return 0, fmt.Errorf("numeric value overflows int64")
	}

	return result.Int64(), nil
}
