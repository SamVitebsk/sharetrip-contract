package service

import (
	"context"
	"fmt"
	"time"

	"github.com/SamVitebsk/sharetrip-contract/internal/domain"
	"github.com/google/uuid"
)

type CreateContractRequest struct {
	ClientID   uuid.UUID
	ValidFrom  time.Time
	ValidUntil *time.Time
	Services   []ContractService
}

type CreateContractResponse struct {
	ID         uuid.UUID
	ClientID   uuid.UUID
	Status     string
	ValidFrom  time.Time
	ValidUntil *time.Time
	Services   []ContractService
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (s *Service) CreateContract(ctx context.Context, request CreateContractRequest) (CreateContractResponse, error) {
	clientID, err := domain.NewClientID(request.ClientID)
	if err != nil {
		return CreateContractResponse{}, mapCreateContractError(err)
	}
	contractID, err := domain.NewContractID(uuid.New())
	if err != nil {
		return CreateContractResponse{}, fmt.Errorf("generate contract ID: %w", err)
	}

	var validUntil *time.Time
	if request.ValidUntil != nil {
		validUntilCopy := request.ValidUntil.UTC().Truncate(time.Microsecond)
		validUntil = &validUntilCopy
	}
	services := make([]domain.ContractService, len(request.Services))
	for index, contractService := range request.Services {
		services[index] = domain.ContractService{
			ServiceCode: domain.ContractServiceCode(contractService.ServiceCode),
			Allowed:     contractService.Allowed,
		}
	}

	response, err := domain.CreateContract(domain.CreateContractRequest{
		ID:         contractID,
		ClientID:   clientID,
		ValidFrom:  request.ValidFrom.UTC().Truncate(time.Microsecond),
		ValidUntil: validUntil,
		Services:   services,
		Now:        time.Now().UTC().Truncate(time.Microsecond),
	})
	if err != nil {
		return CreateContractResponse{}, mapCreateContractError(err)
	}

	err = s.txRunner(ctx, func(transactionCtx context.Context, repoTx RepositoryTx) error {
		if err := repoTx.Create(transactionCtx, response.Contract); err != nil {
			return fmt.Errorf("save new contract: %w", err)
		}
		return nil
	})
	if err != nil {
		return CreateContractResponse{}, fmt.Errorf("create contract: %w", err)
	}

	return toCreateContractResponse(response.Contract), nil
}
