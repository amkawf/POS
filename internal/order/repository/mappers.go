package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"pos-backend/internal/order/domain"
	orderdb "pos-backend/internal/order/repository/generated"
	"pos-backend/internal/pkg/pgconv"
)

func uuidToPgtype(value uuid.UUID) pgtype.UUID               { return pgconv.UUID(value) }
func nullableUUIDToPgtype(value *uuid.UUID) pgtype.UUID      { return pgconv.OptUUID(value) }
func textToPgtype(value *string) pgtype.Text                 { return pgconv.OptText(value) }
func timestamptzToPgtype(value time.Time) pgtype.Timestamptz { return pgconv.Timestamptz(value) }
func nullableTimestamptzToPgtype(value *time.Time) pgtype.Timestamptz {
	return pgconv.OptTimestamptz(value)
}
func pgtypeToUUID(value pgtype.UUID) uuid.UUID                 { return pgconv.ToUUID(value) }
func pgtypeToNullableUUID(value pgtype.UUID) *uuid.UUID        { return pgconv.ToOptUUID(value) }
func pgtypeToNullableText(value pgtype.Text) *string           { return pgconv.ToOptText(value) }
func pgtypeToTime(value pgtype.Timestamptz) time.Time          { return pgconv.ToTime(value) }
func pgtypeToNullableTime(value pgtype.Timestamptz) *time.Time { return pgconv.ToOptTime(value) }

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
