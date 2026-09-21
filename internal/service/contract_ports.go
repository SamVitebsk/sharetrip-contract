package service

import (
	"context"

	"github.com/SamVitebsk/sharetrip-contract/internal/domain"
)

type RepositoryTx interface {
	Create(ctx context.Context, contract domain.Contract) error

	GetForUpdate(ctx context.Context, id domain.ContractID) (domain.Contract, error)

	SaveSigned(ctx context.Context, contract domain.Contract) (domain.Contract, error)
}

type TxRunner func(ctx context.Context, fn func(context.Context, RepositoryTx) error) error
