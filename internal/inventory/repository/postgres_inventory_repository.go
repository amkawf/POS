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
	// Seluruh pemotongan dibungkus dalam 1 Transaksi Atomik
	return r.txManager.WithinTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		for _, it := range items {
			// Cek tipe pemenuhan menu item: BATCH_COOKING atau MADE_TO_ORDER
			var fulfillmentType string
			err := tx.QueryRow(ctx, `
				SELECT COALESCE(fulfillment_type, 'BATCH_COOKING')
				FROM menu_items
				WHERE id = $1
			`, it.MenuItemID).Scan(&fulfillmentType)
			if err != nil {
				return err
			}

			if fulfillmentType == "MADE_TO_ORDER" {
				// 🍳 Menu Made-to-Order: potong langsung bahan baku di store_ingredient_inventory sesuai resep
				rows, err := tx.Query(ctx, `
					SELECT ingredient_id, quantity_per_portion
					FROM recipe_items
					WHERE menu_item_id = $1
				`, it.MenuItemID)
				if err != nil {
					return err
				}

				type ingUsage struct {
					id       uuid.UUID
					consumed float64
				}
				var usages []ingUsage
				for rows.Next() {
					var ingID uuid.UUID
					var qtyPerPortion float64
					if err := rows.Scan(&ingID, &qtyPerPortion); err != nil {
						rows.Close()
						return err
					}
					usages = append(usages, ingUsage{
						id:       ingID,
						consumed: qtyPerPortion * float64(it.Quantity),
					})
				}
				rows.Close()
				if err := rows.Err(); err != nil {
					return err
				}

				// Kurangi stok bahan baku & catat mutasi buku besar bahan
				for _, u := range usages {
					_, err := tx.Exec(ctx, `
						INSERT INTO store_ingredient_inventory (company_id, store_id, ingredient_id, stock, updated_at)
						VALUES ($1, $2, $3, $4, NOW())
						ON CONFLICT (store_id, ingredient_id)
						DO UPDATE SET 
							stock = store_ingredient_inventory.stock - $5,
							updated_at = NOW()
					`, companyID, storeID, u.id, -u.consumed, u.consumed)
					if err != nil {
						return err
					}

					movementNote := fmt.Sprintf("Penjualan Order #%s (%s x%d)", orderID.String()[:8], it.ItemName, it.Quantity)
					_, err = tx.Exec(ctx, `
						INSERT INTO ingredient_movements (
							company_id, store_id, ingredient_id, quantity, type, reference_id, notes, created_at
						) VALUES (
							$1, $2, $3, $4, 'SALE', $5, $6, NOW()
						)
					`, companyID, storeID, u.id, -u.consumed, orderID, movementNote)
					if err != nil {
						return err
					}
				}
			} else {
				// 📦 Menu Batch Cooking: potong stok porsi jadi di store_inventory
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

				// Catat riwayat mutasi ke Buku Besar makanan jadi
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
		}

		return nil
	})
}

// Menghasilkan map: menuItemId -> jumlah stok
func (r *PostgresInventoryRepository) GetStockByStore(
	ctx context.Context,
	storeID uuid.UUID,
) (map[uuid.UUID]int64, error) {
	// 1. Ambil stok porsi makanan jadi (hanya untuk BATCH_COOKING)
	rows, err := r.db.Query(ctx, `
		SELECT si.menu_item_id, si.stock
		FROM store_inventory si
		JOIN menu_items m ON m.id = si.menu_item_id
		WHERE si.store_id = $1 
		  AND COALESCE(m.fulfillment_type, 'BATCH_COOKING') = 'BATCH_COOKING'
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
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// 2. Hitung stok virtual porsi untuk menu MADE_TO_ORDER secara dinamis dari sisa bahan baku di store_ingredient_inventory
	mtoRows, err := r.db.Query(ctx, `
		SELECT 
			m.id,
			CASE 
				WHEN COUNT(r.ingredient_id) = 0 THEN 0
				ELSE COALESCE(
					GREATEST(0, FLOOR(MIN(COALESCE(sii.stock, 0) / NULLIF(r.quantity_per_portion, 0))))::BIGINT,
					0
				)
			END AS available_stock
		FROM menu_items m
		LEFT JOIN recipe_items r ON r.menu_item_id = m.id
		LEFT JOIN store_ingredient_inventory sii 
			ON sii.ingredient_id = r.ingredient_id 
			AND sii.store_id = $1
		WHERE m.fulfillment_type = 'MADE_TO_ORDER'
		GROUP BY m.id
	`, storeID)
	if err == nil {
		defer mtoRows.Close()
		for mtoRows.Next() {
			var menuItemID uuid.UUID
			var availableStock int64
			if err := mtoRows.Scan(&menuItemID, &availableStock); err == nil {
				stockMap[menuItemID] = availableStock
			}
		}
	}

	return stockMap, nil
}

func (r *PostgresInventoryRepository) AdjustStock(
	ctx context.Context,
	companyID, storeID, menuItemID uuid.UUID,
	quantity int64,
	notes string,
) error {
	return r.txManager.WithinTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		var fulfillmentType string
		err := tx.QueryRow(ctx, `
			SELECT COALESCE(fulfillment_type, 'BATCH_COOKING')
			FROM menu_items
			WHERE id = $1
		`, menuItemID).Scan(&fulfillmentType)
		if err != nil {
			return err
		}
		if fulfillmentType == "MADE_TO_ORDER" {
			return fmt.Errorf("menu bertipe Made-to-Order tidak memiliki stok makanan jadi fisik; stok dihitung langsung dari bahan baku")
		}

		// 1. Eksekusi UPSERT ke store_inventory (sertakan company_id)
		_, err = tx.Exec(ctx, `
			INSERT INTO store_inventory (company_id, store_id, menu_item_id, stock, updated_at)
			VALUES ($1, $2, $3, $4, NOW())
			ON CONFLICT (store_id, menu_item_id)
			DO UPDATE SET 
				stock = store_inventory.stock + EXCLUDED.stock,
				updated_at = NOW()
		`, companyID, storeID, menuItemID, quantity)
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
