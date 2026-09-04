package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"pos-backend/internal/order/domain"
	orderdb "pos-backend/internal/order/repository/generated"
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

func pgtypeToUUID(value pgtype.UUID) uuid.UUID {
	return uuid.UUID(value.Bytes)
}

func pgtypeToNullableUUID(value pgtype.UUID) *uuid.UUID {
	if !value.Valid {
		return nil
	}
	id := uuid.UUID(value.Bytes)
	return &id
}

func pgtypeToNullableText(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}
	str := value.String
	return &str
}

func pgtypeToTime(value pgtype.Timestamptz) time.Time {
	return value.Time
}

func pgtypeToNullableTime(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}
	t := value.Time
	return &t
}

func toDomainOrderItem(row orderdb.OrderItem) (domain.OrderItem, error) {
	qty, err := numericToInt64(row.Quantity)
	if err != nil {
		return domain.OrderItem{}, err
	}
	unitPrice, err := numericToInt64(row.UnitPrice)
	if err != nil {
		return domain.OrderItem{}, err
	}
	modifierAmt, err := numericToInt64(row.ModifierAmount)
	if err != nil {
		return domain.OrderItem{}, err
	}
	discountAmt, err := numericToInt64(row.DiscountAmount)
	if err != nil {
		return domain.OrderItem{}, err
	}
	taxAmt, err := numericToInt64(row.TaxAmount)
	if err != nil {
		return domain.OrderItem{}, err
	}
	totalAmt, err := numericToInt64(row.TotalAmount)
	if err != nil {
		return domain.OrderItem{}, err
	}

	return domain.OrderItem{
		ID:             pgtypeToUUID(row.ID),
		OrderID:        pgtypeToUUID(row.OrderID),
		MenuItemID:     pgtypeToUUID(row.MenuItemID),
		ItemName:       row.ItemName,
		SKU:            row.Sku,
		Quantity:       qty,
		UnitPrice:      unitPrice,
		ModifierAmount: modifierAmt,
		DiscountAmount: discountAmt,
		TaxAmount:      taxAmt,
		TotalAmount:    totalAmt,
		Notes:          pgtypeToNullableText(row.Notes),
		Status:         row.Status,
		CreatedAt:      pgtypeToTime(row.CreatedAt),
		UpdatedAt:      pgtypeToTime(row.UpdatedAt),
	}, nil
}

func toDomainOrder(row orderdb.Order, items []domain.OrderItem) (*domain.Order, error) {
	subtotal, err := numericToInt64(row.Subtotal)
	if err != nil {
		return nil, err
	}
	discountAmt, err := numericToInt64(row.DiscountAmount)
	if err != nil {
		return nil, err
	}
	taxAmt, err := numericToInt64(row.TaxAmount)
	if err != nil {
		return nil, err
	}
	serviceAmt, err := numericToInt64(row.ServiceAmount)
	if err != nil {
		return nil, err
	}
	totalAmt, err := numericToInt64(row.TotalAmount)
	if err != nil {
		return nil, err
	}

	return &domain.Order{
		ID:                pgtypeToUUID(row.ID),
		CompanyID:         pgtypeToUUID(row.CompanyID),
		StoreID:           pgtypeToUUID(row.StoreID),
		TableID:           pgtypeToNullableUUID(row.TableID),
		CustomerSessionID: pgtypeToNullableUUID(row.CustomerSessionID),
		OrderNumber:       row.OrderNumber,
		OrderType:         domain.OrderType(row.OrderType),
		OrderSource:       domain.OrderSource(row.OrderSource),
		Status:            domain.OrderStatus(row.Status),
		CustomerName:      pgtypeToNullableText(row.CustomerName),
		Items:             items,
		Subtotal:          subtotal,
		DiscountAmount:    discountAmt,
		TaxAmount:         taxAmt,
		ServiceAmount:     serviceAmt,
		TotalAmount:       totalAmt,
		Notes:             pgtypeToNullableText(row.Notes),
		OpenedAt:          pgtypeToTime(row.OpenedAt),
		CompletedAt:       pgtypeToNullableTime(row.CompletedAt),
		CancelledAt:       pgtypeToNullableTime(row.CancelledAt),
		CreatedBy:         pgtypeToNullableUUID(row.CreatedBy),
		CreatedAt:         pgtypeToTime(row.CreatedAt),
		UpdatedAt:         pgtypeToTime(row.UpdatedAt),
	}, nil
}
