package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"pos-backend/internal/auth/domain"
)

var ErrUserNotFound = errors.New("pengguna tidak ditemukan atau PIN tidak valid")

type UserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	FindByEmail(ctx context.Context, companyID uuid.UUID, email string) (*domain.User, error)
	ListActiveByStore(ctx context.Context, storeID uuid.UUID) ([]domain.User, error)
	FindByStoreAndPIN(ctx context.Context, storeID uuid.UUID, plainPIN string) (*domain.User, error)
}

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

// FindByID mencari user berdasarkan ID unik
func (r *PostgresUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, company_id, name, email, password_hash, pin_hash, role, status, created_at, updated_at
		FROM users
		WHERE id = $1 AND status = 'ACTIVE'
	`
	var u domain.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.CompanyID,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.PinHash,
		&u.Role,
		&u.Status,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

// FindByEmail mencari user untuk login master manager/owner
func (r *PostgresUserRepository) FindByEmail(ctx context.Context, companyID uuid.UUID, email string) (*domain.User, error) {
	query := `
		SELECT id, company_id, name, email, password_hash, pin_hash, role, status, created_at, updated_at
		FROM users
		WHERE company_id = $1 AND email = $2 AND status = 'ACTIVE'
	`
	var u domain.User
	err := r.db.QueryRow(ctx, query, companyID, email).Scan(
		&u.ID,
		&u.CompanyID,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.PinHash,
		&u.Role,
		&u.Status,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

// ListActiveByStore mengambil seluruh staf aktif yang ditugaskan di 1 cabang toko
func (r *PostgresUserRepository) ListActiveByStore(ctx context.Context, storeID uuid.UUID) ([]domain.User, error) {
	query := `
		SELECT u.id, u.company_id, u.name, u.email, u.password_hash, u.pin_hash, u.role, u.status, u.created_at, u.updated_at
		FROM users u
		JOIN user_stores us ON us.user_id = u.id
		WHERE us.store_id = $1 AND u.status = 'ACTIVE'
		ORDER BY u.name ASC
	`
	rows, err := r.db.Query(ctx, query, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(
			&u.ID,
			&u.CompanyID,
			&u.Name,
			&u.Email,
			&u.PasswordHash,
			&u.PinHash,
			&u.Role,
			&u.Status,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

// FindByStoreAndPIN mencari staf di cabang toko tertentu yang PIN-nya cocok
func (r *PostgresUserRepository) FindByStoreAndPIN(ctx context.Context, storeID uuid.UUID, plainPIN string) (*domain.User, error) {
	staffList, err := r.ListActiveByStore(ctx, storeID)
	if err != nil {
		return nil, err
	}

	for _, staff := range staffList {
		if staff.CheckPIN(plainPIN) {
			return &staff, nil
		}
	}

	return nil, ErrUserNotFound
}
