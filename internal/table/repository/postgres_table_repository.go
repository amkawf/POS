package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"pos-backend/internal/table/domain"
	tabledb "pos-backend/internal/table/repository/generated"
)

type PostgresTableRepository struct {
	queries *tabledb.Queries
}

func NewPostgresTableRepository(
	queries *tabledb.Queries,
) *PostgresTableRepository {
	return &PostgresTableRepository{
		queries: queries,
	}
}

func (r *PostgresTableRepository) ListActiveByStore(
	ctx context.Context,
	companyID, storeID uuid.UUID,
) ([]domain.Table, error) {
	rows, err := r.queries.ListActiveTablesByStore(
		ctx,
		tabledb.ListActiveTablesByStoreParams{
			CompanyID: pgtype.UUID{Bytes: companyID, Valid: true},
			StoreID:   pgtype.UUID{Bytes: storeID, Valid: true},
		},
	)
	if err != nil {
		return nil, err
	}

	tables := make([]domain.Table, 0, len(rows))
	for _, row := range rows {
		tables = append(tables, toDomainTable(row))
	}

	return tables, nil
}

func (r *PostgresTableRepository) GetByID(
	ctx context.Context,
	companyID, storeID, id uuid.UUID,
) (*domain.Table, error) {
	row, err := r.queries.GetTableByID(
		ctx,
		tabledb.GetTableByIDParams{
			CompanyID: pgtype.UUID{Bytes: companyID, Valid: true},
			StoreID:   pgtype.UUID{Bytes: storeID, Valid: true},
			ID:        pgtype.UUID{Bytes: id, Valid: true},
		},
	)
	if err != nil {
		return nil, err
	}

	table := toDomainTable(row)
	return &table, nil
}

func (r *PostgresTableRepository) UpdateStatus(
	ctx context.Context,
	companyID, storeID, id uuid.UUID,
	status domain.TableStatus,
) error {
	return r.queries.UpdateTableStatus(
		ctx,
		tabledb.UpdateTableStatusParams{
			CompanyID: pgtype.UUID{Bytes: companyID, Valid: true},
			StoreID:   pgtype.UUID{Bytes: storeID, Valid: true},
			ID:        pgtype.UUID{Bytes: id, Valid: true},
			Status:    string(status),
		},
	)
}

func toDomainTable(row tabledb.Table) domain.Table {
	return domain.Table{
		ID:          uuid.UUID(row.ID.Bytes),
		CompanyID:   uuid.UUID(row.CompanyID.Bytes),
		StoreID:     uuid.UUID(row.StoreID.Bytes),
		TableNumber: row.TableNumber,
		Capacity:    row.Capacity,
		Status:      domain.TableStatus(row.Status),
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
}
