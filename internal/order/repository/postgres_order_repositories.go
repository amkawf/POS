package repository

import (
	"context"

	"github.com/jackc/pgx/v5"

	"errors"
	"pos-backend/internal/database"
	"pos-backend/internal/order/domain"
	orderdb "pos-backend/internal/order/repository/generated"

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
					TableID:        nullableUUIDToPgtype(order.TableID),
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
	row, err := r.queries.GetOrderByID(ctx, orderdb.GetOrderByIDParams{
		CompanyID: uuidToPgtype(companyID),
		StoreID:   uuidToPgtype(storeID),
		ID:        uuidToPgtype(orderID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	itemRows, err := r.queries.ListOrderItemsByOrderID(ctx, row.ID)
	if err != nil {
		return nil, err
	}

	domainItems := make([]domain.OrderItem, 0, len(itemRows))
	for _, itemRow := range itemRows {
		item, err := toDomainOrderItem(itemRow)
		if err != nil {
			return nil, err
		}
		domainItems = append(domainItems, item)
	}

	return toDomainOrder(row, domainItems)
}

func (r *PostgresOrderRepository) ListByStore(
	ctx context.Context,
	companyID uuid.UUID,
	storeID uuid.UUID,
	limit int32,
) ([]*domain.Order, error) {
	if limit <= 0 {
		limit = 50
	}

	rows, err := r.queries.ListOrdersByStore(ctx, orderdb.ListOrdersByStoreParams{
		CompanyID: uuidToPgtype(companyID),
		StoreID:   uuidToPgtype(storeID),
		Limit:     limit,
	})
	if err != nil {
		return nil, err
	}

	orders := make([]*domain.Order, 0, len(rows))
	for _, row := range rows {
		order, err := toDomainOrder(row, nil)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *PostgresOrderRepository) ListByStatus(
	ctx context.Context,
	companyID uuid.UUID,
	storeID uuid.UUID,
	status string,
	limit int32,
) ([]*domain.Order, error) {
	if limit <= 0 {
		limit = 50
	}

	rows, err := r.queries.ListOrdersByStatus(ctx, orderdb.ListOrdersByStatusParams{
		CompanyID: uuidToPgtype(companyID),
		StoreID:   uuidToPgtype(storeID),
		Status:    status,
		Limit:     limit,
	})
	if err != nil {
		return nil, err
	}

	orders := make([]*domain.Order, 0, len(rows))
	for _, row := range rows {
		order, err := toDomainOrder(row, nil)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *PostgresOrderRepository) Update(
	ctx context.Context,
	order *domain.Order,
) error {
	return r.queries.UpdateOrder(ctx, orderdb.UpdateOrderParams{
		CompanyID:   uuidToPgtype(order.CompanyID),
		StoreID:     uuidToPgtype(order.StoreID),
		ID:          uuidToPgtype(order.ID),
		Status:      string(order.Status),
		CompletedAt: nullableTimestamptzToPgtype(order.CompletedAt),
		CancelledAt: nullableTimestamptzToPgtype(order.CancelledAt),
		Notes:       textToPgtype(order.Notes),
		UpdatedAt:   timestamptzToPgtype(order.UpdatedAt),
	})
}

func (r *PostgresOrderRepository) Delete(
	ctx context.Context,
	companyID uuid.UUID,
	storeID uuid.UUID,
	orderID uuid.UUID,
) error {
	return r.txManager.WithinTransaction(
		ctx,
		func(ctx context.Context, tx pgx.Tx) error {
			txQueries := r.queries.WithTx(tx)

			if err := txQueries.DeleteOrderItemsByOrderID(ctx, uuidToPgtype(orderID)); err != nil {
				return err
			}

			return txQueries.DeleteOrder(ctx, orderdb.DeleteOrderParams{
				CompanyID: uuidToPgtype(companyID),
				StoreID:   uuidToPgtype(storeID),
				ID:        uuidToPgtype(orderID),
			})
		},
	)
}
