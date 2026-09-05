package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"pos-backend/internal/menu/domain"
	menudb "pos-backend/internal/menu/repository/generated"
)

type PostgresMenuItemRepository struct {
	queries *menudb.Queries
}

func NewPostgresMenuItemRepository(
	queries *menudb.Queries,
) *PostgresMenuItemRepository {
	return &PostgresMenuItemRepository{
		queries: queries,
	}
}

func (r *PostgresMenuItemRepository) ListActiveByCompany(
	ctx context.Context,
	companyID uuid.UUID,
) ([]domain.MenuItem, error) {
	pgCompanyID := pgtype.UUID{Bytes: companyID, Valid: true}

	rows, err := r.queries.ListActiveMenuItemsByCompany(
		ctx,
		pgCompanyID,
	)
	if err != nil {
		return nil, err
	}

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

	items := make([]domain.MenuItem, 0, len(rows))

	for _, row := range rows {
		basePrice, err := numericToInt64(row.BasePrice)
		if err != nil {
			return nil, err
		}

		var description *string
		if row.Description.Valid {
			description = &row.Description.String
		}

		id := uuid.UUID(row.ID.Bytes)
		categoryIDs := catMap[id]
		if categoryIDs == nil {
			categoryIDs = []uuid.UUID{}
		}

		items = append(items, domain.MenuItem{
			ID:          id,
			CompanyID:   uuid.UUID(row.CompanyID.Bytes),
			SKU:         row.Sku,
			Name:        row.Name,
			Description: description,
			BasePrice:   basePrice,
			Status:      row.Status,
			CategoryIDs: categoryIDs,
		})
	}

	return items, nil
}
