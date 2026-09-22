package repository

import (
	"context"
	"fmt"

	"github.com/SamVitebsk/sharetrip-contract/internal/domain"
	"github.com/SamVitebsk/sharetrip-contract/internal/repository/entity"
)

func (r *ContractTx) GetClientActiveContract(ctx context.Context, clientID domain.ClientID) (domain.Contract, error) {
	const query = `
		SELECT 
			id, client_id, status, valid_from, valid_until, created_at, updated_at
		FROM contracts
		WHERE client_id = $1
		ORDER BY 
			CASE status WHEN 'active' THEN 1 ELSE 2 END,
			created_at DESC
		LIMIT 1
	`

	var ent entity.Contract
	err := r.tx.QueryRow(ctx, query, clientID.Value()).Scan(
		&ent.ID,
		&ent.ClientID,
		&ent.Status,
		&ent.ValidFrom,
		&ent.ValidUntil,
		&ent.CreatedAt,
		&ent.UpdatedAt,
	)
	if err != nil {
		return domain.Contract{}, fmt.Errorf("repository: failed to get active contract: %w", mapPostgresError(err))
	}

	const servicesQuery = `
		SELECT 
			service_code, allowed 
		FROM contract_services 
		WHERE contract_id = $1
	`
	rows, err := r.tx.Query(ctx, servicesQuery, ent.ID)
	if err != nil {
		return domain.Contract{}, fmt.Errorf("repository: failed to get contract services: %w", mapPostgresError(err))
	}
	defer rows.Close()

	var services []entity.ContractService
	for rows.Next() {
		var contractSrv entity.ContractService
		if err := rows.Scan(&contractSrv.ServiceCode, &contractSrv.Allowed); err != nil {
			return domain.Contract{}, fmt.Errorf("repository: failed to scan contract service: %w", err)
		}
		contractSrv.ContractID = ent.ID
		services = append(services, contractSrv)
	}

	if err := rows.Err(); err != nil {
		return domain.Contract{}, fmt.Errorf("repository: rows iteration error: %w", mapPostgresError(err))
	}

	contract, err := toDomainContract(ent, services)
	if err != nil {
		return domain.Contract{}, fmt.Errorf("repository: failed to map contract: %w", err)
	}

	return contract, nil
}
