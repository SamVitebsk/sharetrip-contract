package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type ContractTx struct {
	tx pgx.Tx
}

func newContractTx(transaction pgx.Tx) *ContractTx {
	return &ContractTx{tx: transaction}
}

func (r *RepoPg) WithinContractTx(ctx context.Context, operation func(context.Context, *ContractTx) error) error {
	if operation == nil {
		return errors.New("операция транзакции договора обязательна")
	}
	return tx(ctx, r.pool, func(transactionCtx context.Context, transaction pgx.Tx) error {
		return operation(transactionCtx, newContractTx(transaction))
	})
}
