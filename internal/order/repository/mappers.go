package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func uuidToPgtype(value uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: value,
		Valid: true,
	}
}

func nullableUUIDToPgtype(value *uuid.UUID) pgtype.UUID {
	if value == nil {
		return pgtype.UUID{
			Valid: false,
		}
	}

	return uuidToPgtype(*value)
}

func textToPgtype(value *string) pgtype.Text {
	if value == nil {
		return pgtype.Text{
			Valid: false,
		}
	}

	return pgtype.Text{
		String: *value,
		Valid:  true,
	}
}

func timestamptzToPgtype(value time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  value,
		Valid: true,
	}
}

func nullableTimestamptzToPgtype(
	value *time.Time,
) pgtype.Timestamptz {
	if value == nil {
		return pgtype.Timestamptz{
			Valid: false,
		}
	}

	return timestamptzToPgtype(*value)
}
