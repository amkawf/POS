package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"pos-backend/internal/auth/domain"
)

var ErrNoActiveShift = errors.New("tidak ada sesi kasir yang sedang aktif")

type ShiftRepository interface {
	GetActiveShift(ctx context.Context, storeID, userID uuid.UUID) (*domain.CashierShift, error)
	OpenShift(ctx context.Context, shift *domain.CashierShift) error
	CloseShift(ctx context.Context, shiftID uuid.UUID, actualEndingCash, expectedEndingCash float64, notes *string) error
}

type PostgresShiftRepository struct {
	db *pgxpool.Pool
}

func NewPostgresShiftRepository(db *pgxpool.Pool) *PostgresShiftRepository {
	return &PostgresShiftRepository{db: db}
}

// GetActiveShift mengambil shift kasir yang sedang terbuka di toko ini
func (r *PostgresShiftRepository) GetActiveShift(ctx context.Context, storeID, userID uuid.UUID) (*domain.CashierShift, error) {
	query := `
		SELECT s.id, s.company_id, s.store_id, s.user_id, u.name, s.opened_at, s.starting_cash, s.status
		FROM cashier_shifts s
		JOIN users u ON u.id = s.user_id
		WHERE s.store_id = $1 AND s.user_id = $2 AND s.status = 'OPEN'
		ORDER BY s.opened_at DESC
		LIMIT 1
	`
	var s domain.CashierShift
	err := r.db.QueryRow(ctx, query, storeID, userID).Scan(
		&s.ID,
		&s.CompanyID,
		&s.StoreID,
		&s.UserID,
		&s.UserName,
		&s.OpenedAt,
		&s.StartingCash,
		&s.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoActiveShift
		}
		return nil, err
	}
	return &s, nil
}

// OpenShift mencatat pembukaan laci kasir baru dengan modal uang receh awal
func (r *PostgresShiftRepository) OpenShift(ctx context.Context, shift *domain.CashierShift) error {
	query := `
		INSERT INTO cashier_shifts (company_id, store_id, user_id, starting_cash, status)
		VALUES ($1, $2, $3, $4, 'OPEN')
		RETURNING id, opened_at
	`
	return r.db.QueryRow(ctx, query,
		shift.CompanyID,
		shift.StoreID,
		shift.UserID,
		shift.StartingCash,
	).Scan(&shift.ID, &shift.OpenedAt)
}

// CloseShift mencatat penutupan kasir dan menghitung selisih uang
func (r *PostgresShiftRepository) CloseShift(
	ctx context.Context,
	shiftID uuid.UUID,
	actualEndingCash, expectedEndingCash float64,
	notes *string,
) error {
	query := `
		UPDATE cashier_shifts
		SET status = 'CLOSED',
			closed_at = NOW(),
			actual_ending_cash = $2,
			expected_ending_cash = $3,
			notes = $4
		WHERE id = $1 AND status = 'OPEN'
	`
	tag, err := r.db.Exec(ctx, query, shiftID, actualEndingCash, expectedEndingCash, notes)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("sesi kasir tidak ditemukan atau sudah ditutup")
	}
	return nil
}
