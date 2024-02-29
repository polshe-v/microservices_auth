package transaction

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	"github.com/polshe-v/microservices_auth/internal/client/db"
	"github.com/polshe-v/microservices_auth/internal/client/db/pg"
)

type manager struct {
	db db.Transactor
}

// NewTransactionManager creates transaction manager which implements db.TxManager interface.
func NewTransactionManager(db db.Transactor) db.TxManager {
	return &manager{
		db: db,
	}
}

func (m *manager) ReadCommitted(ctx context.Context, f db.Handler) error {
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.transaction(ctx, txOpts, f)
}

func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, fn db.Handler) (err error) {
	// Check for nested transactions: if there is nested transaction, then there is no need to create a new transaction.
	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if ok {
		return fn(ctx)
	}

	// Start new transaction.
	tx, err = m.db.BeginTx(ctx, opts)
	if err != nil {
		return errors.Wrap(err, "can't begin transaction")
	}

	// Save transaction in context.
	ctx = pg.MakeContextTx(ctx, tx)

	defer func() {
		// Recover from panic.
		if r := recover(); r != nil {
			err = errors.Errorf("panic recovered: %v", r)
		}

		// If there is an error, rollback transaction.
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = errors.Wrapf(err, "errRollback: %v", errRollback)
			}
			return
		}

		// If there is no error, commit transaction.
		if err == nil {
			err = tx.Commit(ctx)
			if err != nil {
				err = errors.Wrapf(err, "errCommit: %v", err)
			}
		}
	}()

	if err = fn(ctx); err != nil {
		err = errors.Wrap(err, "failed to execute transaction")
	}
	return err
}
