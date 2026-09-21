package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamVitebsk/sharetrip-contract/internal/domain"
	"github.com/SamVitebsk/sharetrip-contract/internal/repository/entity"
)

func (r *ContractTx) Create(ctx context.Context, contract domain.Contract) error {
	contractEntity := toContractEntity(contract)
	serviceEntities, err := toNewContractServiceEntities(contract)
	if err != nil {
		return fmt.Errorf("map new contract services: %w", err)
	}
	if len(serviceEntities) == 0 {
		return ErrEmptyContractServices
	}

	_, err = r.tx.Exec(ctx, `
		INSERT INTO contracts (
			id, client_id, status, valid_from, valid_until, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		contractEntity.ID,
		contractEntity.ClientID,
		contractEntity.Status,
		contractEntity.ValidFrom,
		contractEntity.ValidUntil,
		contractEntity.CreatedAt,
		contractEntity.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert contract: %w", mapPostgresError(err))
	}

	serviceValues, serviceArguments := buildContractServiceValues(serviceEntities)
	_, err = r.tx.Exec(ctx, `
		INSERT INTO contract_services (
			contract_id, service_code, allowed, created_at, updated_at
		) VALUES `+serviceValues, serviceArguments...)
	if err != nil {
		return fmt.Errorf("insert contract services: %w", mapPostgresError(err))
	}
	return nil
}

func buildContractServiceValues(serviceEntities []entity.ContractService) (string, []any) {
	const parametersPerService = 5
	serviceValueGroups := make([]string, len(serviceEntities))
	serviceArguments := make([]any, 0, len(serviceEntities)*parametersPerService)
	for index, serviceEntity := range serviceEntities {
		firstParameter := len(serviceArguments) + 1
		serviceValueGroups[index] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)",
			firstParameter, firstParameter+1, firstParameter+2, firstParameter+3, firstParameter+4,
		)
		serviceArguments = append(serviceArguments,
			serviceEntity.ContractID,
			serviceEntity.ServiceCode,
			serviceEntity.Allowed,
			serviceEntity.CreatedAt,
			serviceEntity.UpdatedAt,
		)
	}
	return strings.Join(serviceValueGroups, ", "), serviceArguments
}
