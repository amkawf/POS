package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"pos-backend/internal/menu/domain"
	menudb "pos-backend/internal/menu/repository/generated"
)

type PostgresCategoryRepository struct {
	queries *menudb.Queries
}

func NewPostgresCategoryRepository(
	queries *menudb.Queries,
) *PostgresCategoryRepository {
	return &PostgresCategoryRepository{
		queries: queries,
	}
}

func (r *PostgresCategoryRepository) ListActiveByCompany(
	ctx context.Context,
	companyID uuid.UUID,
) ([]domain.Category, error) {
	rows, err := r.queries.ListActiveMenuCategoriesByCompany(
		ctx,
		pgtype.UUID{Bytes: companyID, Valid: true},
	)
	if err != nil {
		return nil, err
	}

	categories := make([]domain.Category, 0, len(rows))
	for _, row := range rows {
		categories = append(categories, domain.Category{
			ID:        uuid.UUID(row.ID.Bytes),
			MenuID:    uuid.UUID(row.MenuID.Bytes),
			Name:      row.Name,
			SortOrder: row.SortOrder,
			Status:    row.Status,
		})
	}

	return categories, nil
}

func (r *PostgresCategoryRepository) ListCategoryIDsByMenuItem(
	ctx context.Context,
	companyID uuid.UUID,
) (map[uuid.UUID][]uuid.UUID, error) {
	rows, err := r.queries.ListMenuItemCategoryIDsByCompany(
		ctx,
		pgtype.UUID{Bytes: companyID, Valid: true},
	)
	if err != nil {
		return nil, err
	}

	mapping := make(map[uuid.UUID][]uuid.UUID)
	for _, row := range rows {
		itemID := uuid.UUID(row.MenuItemID.Bytes)
		catID := uuid.UUID(row.MenuCategoryID.Bytes)
		mapping[itemID] = append(mapping[itemID], catID)
	}

	return mapping, nil
}
