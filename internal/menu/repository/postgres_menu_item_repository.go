package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"pos-backend/internal/menu/domain"
	menudb "pos-backend/internal/menu/repository/generated"
)

type PostgresMenuItemRepository struct {
	db      *pgxpool.Pool
	queries *menudb.Queries
}

func NewPostgresMenuItemRepository(
	db *pgxpool.Pool,
	queries *menudb.Queries,
) *PostgresMenuItemRepository {
	return &PostgresMenuItemRepository{
		db:      db, // <-- Tambahkan inisialisasi ini
		queries: queries,
	}
}

func (r *PostgresMenuItemRepository) ListActiveByCompany(
	ctx context.Context,
	companyID uuid.UUID,
) ([]domain.MenuItem, error) {
	pgCompanyID := pgtype.UUID{Bytes: companyID, Valid: true}

	catRows, err := r.queries.ListMenuItemCategoryIDsByCompany(
		ctx,
		pgCompanyID,
	)
	if err != nil {
		return nil, err
	}

	catMap := make(map[uuid.UUID][]uuid.UUID)
	for _, cr := range catRows {
		itemID := uuid.UUID(cr.MenuItemID.Bytes)
		catID := uuid.UUID(cr.MenuCategoryID.Bytes)
		catMap[itemID] = append(catMap[itemID], catID)
	}

	// Query menu_items secara langsung agar kolom fulfillment_type ikut terambil
	query := `
		SELECT id, company_id, sku, name, description, base_price, status, COALESCE(fulfillment_type, 'BATCH_COOKING')
		FROM menu_items
		WHERE company_id = $1 AND status = 'ACTIVE'
		ORDER BY name
	`
	rows, err := r.db.Query(ctx, query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.MenuItem
	for rows.Next() {
		var (
			id              uuid.UUID
			cID             uuid.UUID
			sku             string
			name            string
			desc            *string
			basePriceNum    pgtype.Numeric
			status          string
			fulfillmentType string
		)

		if err := rows.Scan(
			&id,
			&cID,
			&sku,
			&name,
			&desc,
			&basePriceNum,
			&status,
			&fulfillmentType,
		); err != nil {
			return nil, err
		}

		basePrice, err := numericToInt64(basePriceNum)
		if err != nil {
			return nil, err
		}

		categoryIDs := catMap[id]
		if categoryIDs == nil {
			categoryIDs = []uuid.UUID{}
		}

		items = append(items, domain.MenuItem{
			ID:              id,
			CompanyID:       cID,
			SKU:             sku,
			Name:            name,
			Description:     desc,
			BasePrice:       basePrice,
			Status:          status,
			CategoryIDs:     categoryIDs,
			FulfillmentType: fulfillmentType,
		})
	}

	return items, rows.Err()
}

// Create menyimpan data menu baru dan menghubungkannya ke kategori terkait
func (r *PostgresMenuItemRepository) Create(ctx context.Context, item *domain.MenuItem) error {
	if item.FulfillmentType == "" {
		item.FulfillmentType = "BATCH_COOKING"
	}
	query := `
		INSERT INTO menu_items (company_id, sku, name, description, base_price, status, fulfillment_type)
		VALUES ($1, $2, $3, $4, $5, 'ACTIVE', $6)
		RETURNING id
	`
	err := r.db.QueryRow(ctx, query,
		item.CompanyID,
		item.SKU,
		item.Name,
		item.Description,
		item.BasePrice,
		item.FulfillmentType,
	).Scan(&item.ID)
	if err != nil {
		return err
	}
	// Hubungkan ke kategori makanan/minuman jika dipilih
	for _, catID := range item.CategoryIDs {
		_, err := r.db.Exec(ctx, `
			INSERT INTO menu_category_items (menu_category_id, menu_item_id, sort_order)
			VALUES ($1, $2, 1)
			ON CONFLICT (menu_category_id, menu_item_id) DO NOTHING
		`, catID, item.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
