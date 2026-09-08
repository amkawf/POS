package repository

import (
	"github.com/jackc/pgx/v5/pgtype"

	"pos-backend/internal/pkg/pgconv"
)

func int64ToNumeric(value int64) pgtype.Numeric {
	return pgconv.Int64ToNumeric(value)
}

func numericToInt64(value pgtype.Numeric) (int64, error) {
	return pgconv.NumericToInt64(value)
}
