package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"pos-backend/internal/database"
	"pos-backend/internal/ingredient/domain"
)

type IngredientRepository interface {
	Create(ctx context.Context, ing *domain.Ingredient) error
	ListByStore(ctx context.Context, companyID, storeID uuid.UUID) ([]domain.Ingredient, error)
	Restock(ctx context.Context, companyID, storeID, ingredientID uuid.UUID, quantity float64, notes string) error
	Update(ctx context.Context, ing *domain.Ingredient) error 
}

type PostgresIngredientRepository struct {
	db        *pgxpool.Pool
	txManager *database.TransactionManager
}

func NewPostgresIngredientRepository(
	db *pgxpool.Pool,
	txManager *database.TransactionManager,
) *PostgresIngredientRepository {
	return &PostgresIngredientRepository{
		db:        db,
		txManager: txManager,
	}
}

func (r *PostgresIngredientRepository) Create(ctx context.Context, ing *domain.Ingredient) error {
	query := `
		INSERT INTO ingredients (company_id, code, name, unit, min_stock_alert)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		ing.CompanyID,
		ing.Code,
		ing.Name,
		ing.Unit,
		ing.MinStockAlert,
	).Scan(&ing.ID, &ing.CreatedAt, &ing.UpdatedAt)
}

func (r *PostgresIngredientRepository) ListByStore(
	ctx context.Context,
	companyID, storeID uuid.UUID,
) ([]domain.Ingredient, error) {
	query := `
		SELECT 
			i.id, i.company_id, i.code, i.name, i.unit, i.min_stock_alert, i.created_at, i.updated_at,
			COALESCE(sii.stock, 0) AS stock
		FROM ingredients i
		LEFT JOIN store_ingredient_inventory sii 
			ON sii.ingredient_id = i.id AND sii.store_id = $2
		WHERE i.company_id = $1
		ORDER BY i.name ASC
	`
	rows, err := r.db.Query(ctx, query, companyID, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.Ingredient
	for rows.Next() {
		var ing domain.Ingredient
		var stock float64
		if err := rows.Scan(
			&ing.ID,
			&ing.CompanyID,
			&ing.Code,
			&ing.Name,
			&ing.Unit,
			&ing.MinStockAlert,
			&ing.CreatedAt,
			&ing.UpdatedAt,
			&stock,
		); err != nil {
			return nil, err
		}
		ing.Stock = &stock
		items = append(items, ing)
	}

	return items, rows.Err()
}

func (r *PostgresIngredientRepository) Restock(
	ctx context.Context,
	companyID, storeID, ingredientID uuid.UUID,
	quantity float64,
	notes string,
) error {
	return r.txManager.WithinTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		// 1. Eksekusi UPSERT penambahan stok bahan baku
		_, err := tx.Exec(ctx, `
			INSERT INTO store_ingredient_inventory (company_id, store_id, ingredient_id, stock, updated_at)
			VALUES ($1, $2, $3, $4, NOW())
			ON CONFLICT (store_id, ingredient_id)
			DO UPDATE SET 
				stock = store_ingredient_inventory.stock + EXCLUDED.stock,
				updated_at = NOW()
		`, companyID, storeID, ingredientID, quantity)
		if err != nil {
			return err
		}

		// 2. Catat riwayat pembelian ke buku besar
		_, err = tx.Exec(ctx, `
			INSERT INTO ingredient_movements (
				company_id, store_id, ingredient_id, quantity, type, notes
			) VALUES (
				$1, $2, $3, $4, 'PURCHASE', $5
			)
		`, companyID, storeID, ingredientID, quantity, notes)
		return err
	})
}

// Update memperbarui atribut master bahan baku (Nama, Satuan, Alert, Kode)
func (r *PostgresIngredientRepository) Update(ctx context.Context, ing *domain.Ingredient) error {
	query := `
		UPDATE ingredients
		SET name = $1, unit = $2, min_stock_alert = $3, code = $4, updated_at = NOW()
		WHERE id = $5 AND company_id = $6
	`
	tag, err := r.db.Exec(ctx, query,
		ing.Name,
		ing.Unit,
		ing.MinStockAlert,
		ing.Code,
		ing.ID,
		ing.CompanyID,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("bahan baku tidak ditemukan atau tidak memiliki akses")
	}

	return nil
}