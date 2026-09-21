package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func tx(
	ctx context.Context,
	pool *pgxpool.Pool,
	block func(context.Context, pgx.Tx) error,
) (err error) {
	if block == nil {
		return errors.New("блок транзакции обязателен")
	}

	transactionCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	transaction, err := pool.Begin(transactionCtx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", mapPostgresError(err))
	}
	defer func() {
		rollbackCtx, cancelRollback := context.WithTimeout(context.WithoutCancel(ctx), 3*time.Second)
		defer cancelRollback()
		if rollbackErr := transaction.Rollback(rollbackCtx); rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			err = errors.Join(err, fmt.Errorf("rollback transaction: %w", rollbackErr))
		}
	}()

	if err := block(transactionCtx, transaction); err != nil {
		return fmt.Errorf("transaction block: %w", err)
	}
	if err := transaction.Commit(transactionCtx); err != nil {
		return fmt.Errorf("commit transaction: %w", mapPostgresError(err))
	}
	return nil
}
