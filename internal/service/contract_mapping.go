package service

import (
	"time"

	"github.com/SamVitebsk/sharetrip-contract/internal/domain"
)

type ContractService struct {
	ServiceCode string
	Allowed     bool
}

func toContractServices(services []domain.ContractService) []ContractService {
	result := make([]ContractService, len(services))
	for index, contractService := range services {
		result[index] = ContractService{
			ServiceCode: string(contractService.ServiceCode),
			Allowed:     contractService.Allowed,
		}
	}
	return result
}

func copyContractEndTime(validUntil *time.Time) *time.Time {
	if validUntil == nil {
		return nil
	}
	validUntilCopy := *validUntil
	return &validUntilCopy
}

func toCreateContractResult(contract domain.Contract) CreateContractResult {
	return CreateContractResult{
		ID:         contract.ID.Value(),
		ClientID:   contract.ClientID.Value(),
		Status:     string(contract.Status),
		ValidFrom:  contract.ValidFrom,
		ValidUntil: copyContractEndTime(contract.ValidUntil),
		Services:   toContractServices(contract.Services),
		CreatedAt:  contract.CreatedAt,
		UpdatedAt:  contract.UpdatedAt,
	}
}

func toSignContractResult(contract domain.Contract) SignContractResult {
	return SignContractResult{
		ID:         contract.ID.Value(),
		ClientID:   contract.ClientID.Value(),
		Status:     string(contract.Status),
		ValidFrom:  contract.ValidFrom,
		ValidUntil: copyContractEndTime(contract.ValidUntil),
		Services:   toContractServices(contract.Services),
		CreatedAt:  contract.CreatedAt,
		UpdatedAt:  contract.UpdatedAt,
	}
}
