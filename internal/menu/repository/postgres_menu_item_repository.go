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
	rows, err := r.queries.ListActiveMenuItemsByCompany(
		ctx,
		pgtype.UUID{Bytes: companyID, Valid: true},
	)
	if err != nil {
		return nil, err
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

		items = append(items, domain.MenuItem{
			ID:          uuid.UUID(row.ID.Bytes),
			CompanyID:   uuid.UUID(row.CompanyID.Bytes),
			SKU:         row.Sku,
			Name:        row.Name,
			Description: description,
			BasePrice:   basePrice,
			Status:      row.Status,
		})
	}

	return items, nil
}
