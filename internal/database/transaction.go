package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type TransactionManager struct {
	db interface {
		BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
	}
}

func NewTransactionManager(db interface {
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
}) *TransactionManager {
	return &TransactionManager{
		db: db,
	}
}

func (tm *TransactionManager) WithinTransaction(
	ctx context.Context,
	fn func(ctx context.Context, tx pgx.Tx) error,
) error {
	tx, err := tm.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err := fn(ctx, tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}