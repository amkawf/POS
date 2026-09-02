package repository

import (
	"context"

	"github.com/jackc/pgx/v5"

	"pos-backend/internal/database"
	"pos-backend/internal/order/domain"
	orderdb "pos-backend/internal/order/repository/generated"
	"errors"
	"github.com/google/uuid"
	
)

type PostgresOrderRepository struct {
	queries   *orderdb.Queries
	txManager *database.TransactionManager
}

func NewPostgresOrderRepository(
	queries *orderdb.Queries,
	txManager *database.TransactionManager,
) *PostgresOrderRepository {
	return &PostgresOrderRepository{
		queries:   queries,
		txManager: txManager,
	}
}

func (r *PostgresOrderRepository) Create(
	ctx context.Context,
	order *domain.Order,
) error {
	return r.txManager.WithinTransaction(
		ctx,
		func(ctx context.Context, tx pgx.Tx) error {
			txQueries := r.queries.WithTx(tx)

			err := txQueries.CreateOrder(
				ctx,
				orderdb.CreateOrderParams{
					ID:             uuidToPgtype(order.ID),
					CompanyID:      uuidToPgtype(order.CompanyID),
					StoreID:        uuidToPgtype(order.StoreID),
					OrderNumber:    order.OrderNumber,
					OrderType:      string(order.OrderType),
					OrderSource:    string(order.OrderSource),
					Status:         string(order.Status),
					CustomerName:   textToPgtype(order.CustomerName),
					Subtotal:       int64ToNumeric(order.Subtotal),
					DiscountAmount: int64ToNumeric(order.DiscountAmount),
					TaxAmount:      int64ToNumeric(order.TaxAmount),
					ServiceAmount:  int64ToNumeric(order.ServiceAmount),
					TotalAmount:    int64ToNumeric(order.TotalAmount),
					Notes:          textToPgtype(order.Notes),
					OpenedAt:       timestamptzToPgtype(order.OpenedAt),
					CompletedAt:    nullableTimestamptzToPgtype(order.CompletedAt),
					CancelledAt:    nullableTimestamptzToPgtype(order.CancelledAt),
					CreatedBy:      nullableUUIDToPgtype(order.CreatedBy),
					CreatedAt:      timestamptzToPgtype(order.CreatedAt),
					UpdatedAt:      timestamptzToPgtype(order.UpdatedAt),
				},
			)
			if err != nil {
				return err
			}

			for _, item := range order.Items {
				err := txQueries.CreateOrderItem(
					ctx,
					orderdb.CreateOrderItemParams{
						ID:             uuidToPgtype(item.ID),
						OrderID:        uuidToPgtype(item.OrderID),
						MenuItemID:     uuidToPgtype(item.MenuItemID),
						ItemName:       item.ItemName,
						Sku:            item.SKU,
						Quantity:       int64ToNumeric(item.Quantity),
						UnitPrice:      int64ToNumeric(item.UnitPrice),
						ModifierAmount: int64ToNumeric(item.ModifierAmount),
						DiscountAmount: int64ToNumeric(item.DiscountAmount),
						TaxAmount:      int64ToNumeric(item.TaxAmount),
						TotalAmount:    int64ToNumeric(item.TotalAmount),
						Notes:          textToPgtype(item.Notes),
						Status:         item.Status,
						CreatedAt:      timestamptzToPgtype(item.CreatedAt),
						UpdatedAt:      timestamptzToPgtype(item.UpdatedAt),
					},
				)
				if err != nil {
					return err
				}
			}

			return nil
		},
	)
}

func (r *PostgresOrderRepository) FindByID(
	ctx context.Context,
	companyID uuid.UUID,
	storeID uuid.UUID,
	orderID uuid.UUID,
) (*domain.Order, error) {
	return nil, errors.New("FindByID is not implemented")
}

func (r *PostgresOrderRepository) Update(
	ctx context.Context,
	order *domain.Order,
) error {
	return errors.New("Update is not implemented")
}
