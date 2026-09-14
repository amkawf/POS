package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"pos-backend/internal/database"
	"pos-backend/internal/recipe/domain"
)

type RecipeRepository interface {
	GetByMenuItemID(ctx context.Context, menuItemID uuid.UUID) ([]domain.RecipeItem, error)
	SaveRecipe(ctx context.Context, menuItemID uuid.UUID, items []domain.RecipeItem) error
}

type PostgresRecipeRepository struct {
	db        *pgxpool.Pool
	txManager *database.TransactionManager
}

func NewPostgresRecipeRepository(
	db *pgxpool.Pool,
	txManager *database.TransactionManager,
) *PostgresRecipeRepository {
	return &PostgresRecipeRepository{
		db:        db,
		txManager: txManager,
	}
}

// GetByMenuItemID mengambil seluruh takaran bahan untuk satu menu makanan
func (r *PostgresRecipeRepository) GetByMenuItemID(
	ctx context.Context,
	menuItemID uuid.UUID,
) ([]domain.RecipeItem, error) {
	query := `
		SELECT 
			r.id, r.menu_item_id, r.ingredient_id, i.name, i.unit, r.quantity_per_portion
		FROM recipe_items r
		JOIN ingredients i ON i.id = r.ingredient_id
		WHERE r.menu_item_id = $1
		ORDER BY i.name ASC
	`
	rows, err := r.db.Query(ctx, query, menuItemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.RecipeItem
	for rows.Next() {
		var item domain.RecipeItem
		if err := rows.Scan(
			&item.ID,
			&item.MenuItemID,
			&item.IngredientID,
			&item.IngredientName,
			&item.IngredientUnit,
			&item.QuantityPerPortion,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

// SaveRecipe mengganti/menyimpan seluruh komposisi resep untuk satu menu dalam 1 transaksi atomik
func (r *PostgresRecipeRepository) SaveRecipe(
	ctx context.Context,
	menuItemID uuid.UUID,
	items []domain.RecipeItem,
) error {
	return r.txManager.WithinTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		// 1. Bersihkan resep lama menu ini terlebih dahulu
		_, err := tx.Exec(ctx, `DELETE FROM recipe_items WHERE menu_item_id = $1`, menuItemID)
		if err != nil {
			return err
		}

		// 2. Masukkan komposisi takaran bahan yang baru
		for _, it := range items {
			_, err = tx.Exec(ctx, `
				INSERT INTO recipe_items (menu_item_id, ingredient_id, quantity_per_portion)
				VALUES ($1, $2, $3)
			`, menuItemID, it.IngredientID, it.QuantityPerPortion)
			if err != nil {
				return err
			}
		}

		return nil
	})
}