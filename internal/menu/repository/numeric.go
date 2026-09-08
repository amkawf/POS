package repository

import (
	"github.com/jackc/pgx/v5/pgtype"

	"pos-backend/internal/pkg/pgconv"
)

func numericToInt64(value pgtype.Numeric) (int64, error) {
	return pgconv.NumericToInt64(value)
}
