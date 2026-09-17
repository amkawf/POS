package application

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"pos-backend/internal/database"
)

type BatchProduceUseCase struct {
	db        *pgxpool.Pool
	txManager *database.TransactionManager
}

func NewBatchProduceUseCase(
	db *pgxpool.Pool,
	txManager *database.TransactionManager,
) *BatchProduceUseCase {
	return &BatchProduceUseCase{
		db:        db,
		txManager: txManager,
	}
}

// Execute menjalankan konversi bahan baku menjadi porsi makanan jadi
func (uc *BatchProduceUseCase) Execute(
	ctx context.Context,
	companyID, storeID, menuItemID uuid.UUID,
	portions int64,
	notes string,
) error {
	if portions <= 0 {
		return fmt.Errorf("jumlah porsi yang dimasak harus lebih dari 0")
	}

	if notes == "" {
		notes = fmt.Sprintf("Produksi dapur: %d porsi", portions)
	}

	// Seluruh pemotongan bahan & penambahan porsi jadi dibungkus 1 Transaksi Atomik
	return uc.txManager.WithinTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
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
			return fmt.Errorf("menu bertipe Made-to-Order tidak memerlukan batch masak; porsi otomatis tersedia berdasarkan stok bahan baku")
		}

		// 1. Ambil resep komposisi bahan baku untuk menu ini
		rows, err := tx.Query(ctx, `
			SELECT ingredient_id, quantity_per_portion
			FROM recipe_items
			WHERE menu_item_id = $1
		`, menuItemID)
		if err != nil {
			return err
		}
		defer rows.Close()

		type ingredientUsage struct {
			id       uuid.UUID
			consumed float64
		}
		var usages []ingredientUsage

		for rows.Next() {
			var ingID uuid.UUID
			var qtyPerPortion float64
			if err := rows.Scan(&ingID, &qtyPerPortion); err != nil {
				return err
			}
			usages = append(usages, ingredientUsage{
				id:       ingID,
				consumed: qtyPerPortion * float64(portions),
			})
		}
		if err := rows.Err(); err != nil {
			return err
		}

		// Validasi: Pastikan menu sudah memiliki konfigurasi resep bahan baku
		if len(usages) == 0 {
			return fmt.Errorf("menu ini belum memiliki konfigurasi resep (BOM). Harap atur resep terlebih dahulu")
		}

		// 2. Potong stok bahan baku di gudang toko (Opsi B: toleran nilai minus)
		for _, u := range usages {
			// Kurangi stok bahan mentah via UPSERT
			_, err := tx.Exec(ctx, `
				INSERT INTO store_ingredient_inventory (company_id, store_id, ingredient_id, stock, updated_at)
				VALUES ($1, $2, $3, -1 * $4, NOW())
				ON CONFLICT (store_id, ingredient_id)
				DO UPDATE SET 
					stock = store_ingredient_inventory.stock - $4,
					updated_at = NOW()
			`, companyID, storeID, u.id, u.consumed)
			if err != nil {
				return err
			}

			// Catat mutasi pemakaian masak ke buku besar bahan baku
			_, err = tx.Exec(ctx, `
				INSERT INTO ingredient_movements (
					company_id, store_id, ingredient_id, quantity, type, notes
				) VALUES (
					$1, $2, $3, -1 * $4, 'COOKING_USAGE', $5
				)
			`, companyID, storeID, u.id, u.consumed, notes)
			if err != nil {
				return err
			}
		}

		// 3. Tambahkan porsi makanan jadi yang siap jual di kasir
		_, err = tx.Exec(ctx, `
			INSERT INTO store_inventory (company_id, store_id, menu_item_id, stock, updated_at)
			VALUES ($1, $2, $3, $4, NOW())
			ON CONFLICT (store_id, menu_item_id)
			DO UPDATE SET 
				stock = store_inventory.stock + EXCLUDED.stock,
				updated_at = NOW()
		`, companyID, storeID, menuItemID, portions)
		if err != nil {
			return err
		}

		// 4. Catat riwayat penambahan porsi jadi ke buku besar menu
		_, err = tx.Exec(ctx, `
			INSERT INTO inventory_movements (
				company_id, store_id, menu_item_id, quantity, movement_type, notes
			) VALUES (
				$1, $2, $3, $4, 'RESTOCK', $5
			)
		`, companyID, storeID, menuItemID, portions, notes)
		return err
	})
}
