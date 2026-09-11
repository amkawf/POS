package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"pos-backend/internal/database"
	"pos-backend/internal/inventory/domain"
)

type ItemDeduction struct {
	MenuItemID uuid.UUID
	ItemName   string
	Quantity   int64
}

type InventoryRepository interface {
	DeductStockForOrder(
		ctx context.Context,
		companyID, storeID, orderID uuid.UUID,
		items []ItemDeduction,
	) error
	GetStockByStore(
		ctx context.Context,
		storeID uuid.UUID,
	) (map[uuid.UUID]int64, error)
	AdjustStock(
		ctx context.Context,
		companyID, storeID, menuItemID uuid.UUID,
		quantity int64,
		notes string,
	) error
}

type PostgresInventoryRepository struct {
	db        *pgxpool.Pool
	txManager *database.TransactionManager
}

func NewPostgresInventoryRepository(
	db *pgxpool.Pool,
	txManager *database.TransactionManager,
) *PostgresInventoryRepository {
	return &PostgresInventoryRepository{
		db:        db,
		txManager: txManager,
	}
}

func (r *PostgresInventoryRepository) DeductStockForOrder(
	ctx context.Context,
	companyID, storeID, orderID uuid.UUID,
	items []ItemDeduction,
) error {
	//  Seluruh pemotongan dibungkus dalam 1 Transaksi Atomik
	return r.txManager.WithinTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		for _, it := range items {
			// 1. Eksekusi Atomic Conditional Update
			tag, err := tx.Exec(ctx, `
				UPDATE store_inventory
				SET stock = stock - $1, updated_at = NOW()
				WHERE store_id = $2 AND menu_item_id = $3 AND stock >= $1
			`, it.Quantity, storeID, it.MenuItemID)
			if err != nil {
				return err
			}

			// Jika RowsAffected == 0, artinya kondisi 'stock >= $1' GAGAL
			if tag.RowsAffected() == 0 {
				return fmt.Errorf("%w: %s (jumlah diminta: %d)", domain.ErrInsufficientStock, it.ItemName, it.Quantity)
			}

			//  Catat riwayat mutasi ke Buku Besar 
			movementNote := fmt.Sprintf("Penjualan Order #%s", orderID.String()[:8])
			_, err = tx.Exec(ctx, `
				INSERT INTO inventory_movements (
					company_id, store_id, menu_item_id, order_id, quantity, movement_type, notes
				) VALUES (
					$1, $2, $3, $4, $5, 'SALE', $6
				)
			`, companyID, storeID, it.MenuItemID, orderID, -float64(it.Quantity), movementNote)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// Menghasilkan map: menuItemId -> jumlah stok
func (r *PostgresInventoryRepository) GetStockByStore(
	ctx context.Context,
	storeID uuid.UUID,
	) (map[uuid.UUID]int64, error) {
		rows, err := r.db.Query(ctx, `
		SELECT menu_item_id, stock
		FROM store_inventory
		WHERE store_id = $1
		`, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stockMap := make(map[uuid.UUID]int64)
	for rows.Next() {
		var menuItemID uuid.UUID
		var stock int64
		if err := rows.Scan(&menuItemID, &stock); err != nil {
			return nil, err
		}
		stockMap[menuItemID] = stock
	}
	return stockMap, rows.Err()
}

func (r *PostgresInventoryRepository) AdjustStock(
	ctx context.Context,
	companyID, storeID, menuItemID uuid.UUID,
	quantity int64,
	notes string,
) error {
	return r.txManager.WithinTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		// 1. Eksekusi UPSERT ke store_inventory (sertakan company_id)
		_, err := tx.Exec(ctx, `
			INSERT INTO store_inventory (company_id, store_id, menu_item_id, stock, updated_at)
			VALUES ($1, $2, $3, $4, NOW())
			ON CONFLICT (store_id, menu_item_id)
			DO UPDATE SET 
				stock = store_inventory.stock + EXCLUDED.stock,
				updated_at = NOW()
		`, companyID, storeID, menuItemID, quantity) // <-- oper companyID sebagai $1
		if err != nil {
			return err
		}

		// 2. Catat ke Buku Besar (inventory_movements)
		movementType := "RESTOCK"
		if quantity < 0 {
			movementType = "WASTE"
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO inventory_movements (
				company_id, store_id, menu_item_id, quantity, movement_type, notes
			) VALUES (
				$1, $2, $3, $4, $5, $6
			)
		`, companyID, storeID, menuItemID, float64(quantity), movementType, notes)
		return err
	})
}