package repository

import (
	"context"
	"fmt"
	"slices"

	"github.com/SamVitebsk/sharetrip-contract/internal/domain"
	"github.com/SamVitebsk/sharetrip-contract/internal/repository/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *ContractTx) GetForUpdate(ctx context.Context, id domain.ContractID) (domain.Contract, error) {
	contractEntity, err := scanContract(r.tx.QueryRow(ctx, `
		SELECT id, client_id, status, valid_from, valid_until, created_at, updated_at
		FROM contracts
		WHERE id = $1
		FOR UPDATE`, id.Value()))
	if err != nil {
		return domain.Contract{}, fmt.Errorf("get contract for update: %w", mapPostgresError(err))
	}
	serviceEntities, err := r.loadContractServices(ctx, contractEntity.ID)
	if err != nil {
		return domain.Contract{}, err
	}
	return toDomainContract(contractEntity, serviceEntities)
}

func (r *ContractTx) SaveSigned(ctx context.Context, contract domain.Contract) (domain.Contract, error) {
	err := r.tx.QueryRow(ctx, `
		UPDATE contracts
		SET status = $2, updated_at = $3
		WHERE id = $1
		RETURNING status, updated_at`,
		contract.ID.Value(), string(contract.Status), contract.UpdatedAt,
	).Scan(&contract.Status, &contract.UpdatedAt)
	if err != nil {
		return domain.Contract{}, fmt.Errorf("save signed contract: %w", mapPostgresError(err))
	}
	contract.Services = slices.Clone(contract.Services)
	contract.ValidUntil = copyContractEndTime(contract.ValidUntil)
	return contract, nil
}

func scanContract(row pgx.Row) (entity.Contract, error) {
	var contractEntity entity.Contract
	if err := row.Scan(
		&contractEntity.ID,
		&contractEntity.ClientID,
		&contractEntity.Status,
		&contractEntity.ValidFrom,
		&contractEntity.ValidUntil,
		&contractEntity.CreatedAt,
		&contractEntity.UpdatedAt,
	); err != nil {
		return entity.Contract{}, err
	}
	return contractEntity, nil
}

func (r *ContractTx) loadContractServices(ctx context.Context, contractID uuid.UUID) ([]entity.ContractService, error) {
	rows, err := r.tx.Query(ctx, `
		SELECT contract_id, service_code, allowed, created_at, updated_at
		FROM contract_services
		WHERE contract_id = $1
		ORDER BY service_code`, contractID)
	if err != nil {
		return nil, fmt.Errorf("query contract services: %w", mapPostgresError(err))
	}
	defer rows.Close()

	var serviceEntities []entity.ContractService
	for rows.Next() {
		var serviceEntity entity.ContractService
		if err := rows.Scan(
			&serviceEntity.ContractID,
			&serviceEntity.ServiceCode,
			&serviceEntity.Allowed,
			&serviceEntity.CreatedAt,
			&serviceEntity.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan contract service: %w", mapPostgresError(err))
		}
		serviceEntities = append(serviceEntities, serviceEntity)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read contract services: %w", mapPostgresError(err))
	}
	return serviceEntities, nil
}
