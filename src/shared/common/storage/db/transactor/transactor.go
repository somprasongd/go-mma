// Ref: https://github.com/Thiht/transactor/blob/main/sqlx/transactor.go
package transactor

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Transactor interface {
	WithinTransaction(ctx context.Context, txFunc func(context.Context) error) error
}

func New(db *sqlx.DB, nestedTransactionStrategy nestedTransactionsStrategy) (Transactor, DBContext) {
	sqlDBGetter := func(ctx context.Context) sqlxDB {
		if tx := txFromContext(ctx); tx != nil {
			return tx
		}

		return db
	}

	dbGetter := func(ctx context.Context) DBTX {
		if tx := txFromContext(ctx); tx != nil {
			return tx
		}

		return db
	}

	return &transactor{
		sqlDBGetter,
		nestedTransactionStrategy,
	}, dbGetter
}

type (
	sqlxDBGetter               func(context.Context) sqlxDB
	nestedTransactionsStrategy func(sqlxDB, *sqlx.Tx) (sqlxDB, sqlxTx)
)

type transactor struct {
	sqlxDBGetter
	nestedTransactionsStrategy
}

func (t *transactor) WithinTransaction(ctx context.Context, txFunc func(context.Context) error) error {
	currentDB := t.sqlxDBGetter(ctx)

	tx, err := currentDB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	newDB, currentTX := t.nestedTransactionsStrategy(currentDB, tx)
	defer func() {
		_ = currentTX.Rollback() // If rollback fails, there's nothing to do, the transaction will expire by itself
	}()
	txCtx := txToContext(ctx, newDB)

	if err := txFunc(txCtx); err != nil {
		return err
	}

	if err := currentTX.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func IsWithinTransaction(ctx context.Context) bool {
	return ctx.Value(transactorKey{}) != nil
}
