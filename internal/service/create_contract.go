package service

import (
	"context"
	"fmt"
	"time"

	"github.com/SamVitebsk/sharetrip-contract/internal/domain"
	"github.com/google/uuid"
)

type CreateContractCommand struct {
	ClientID   uuid.UUID
	ValidFrom  time.Time
	ValidUntil *time.Time
	Services   []ContractService
}

type CreateContractResult struct {
	ID         uuid.UUID
	ClientID   uuid.UUID
	Status     string
	ValidFrom  time.Time
	ValidUntil *time.Time
	Services   []ContractService
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (s *Service) CreateContract(ctx context.Context, command CreateContractCommand) (CreateContractResult, error) {
	clientID, err := domain.NewClientID(command.ClientID)
	if err != nil {
		return CreateContractResult{}, mapCreateContractError(err)
	}
	contractID, err := domain.NewContractID(uuid.New())
	if err != nil {
		return CreateContractResult{}, fmt.Errorf("generate contract ID: %w", err)
	}

	var validUntil *time.Time
	if command.ValidUntil != nil {
		validUntilCopy := command.ValidUntil.UTC().Truncate(time.Microsecond)
		validUntil = &validUntilCopy
	}
	services := make([]domain.ContractService, len(command.Services))
	for index, contractService := range command.Services {
		services[index] = domain.ContractService{
			ServiceCode: domain.ContractServiceCode(contractService.ServiceCode),
			Allowed:     contractService.Allowed,
		}
	}

	response, err := domain.CreateContract(domain.CreateContractRequest{
		ID:         contractID,
		ClientID:   clientID,
		ValidFrom:  command.ValidFrom.UTC().Truncate(time.Microsecond),
		ValidUntil: validUntil,
		Services:   services,
		Now:        time.Now().UTC().Truncate(time.Microsecond),
	})
	if err != nil {
		return CreateContractResult{}, mapCreateContractError(err)
	}

	err = s.txRunner(ctx, func(transactionCtx context.Context, repoTx RepositoryTx) error {
		if err := repoTx.Create(transactionCtx, response.Contract); err != nil {
			return fmt.Errorf("save new contract: %w", err)
		}
		return nil
	})
	if err != nil {
		return CreateContractResult{}, fmt.Errorf("create contract: %w", err)
	}

	return toCreateContractResult(response.Contract), nil
}
