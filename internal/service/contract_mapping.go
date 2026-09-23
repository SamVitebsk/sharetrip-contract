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
	responses := make([]ContractService, len(services))
	for index, contractService := range services {
		responses[index] = ContractService{
			ServiceCode: string(contractService.ServiceCode),
			Allowed:     contractService.Allowed,
		}
	}
	return responses
}

func copyContractEndTime(validUntil *time.Time) *time.Time {
	if validUntil == nil {
		return nil
	}
	validUntilCopy := *validUntil
	return &validUntilCopy
}

func toCreateContractResponse(contract domain.Contract) CreateContractResponse {
	return CreateContractResponse{
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

func toSignContractResponse(contract domain.Contract) SignContractResponse {
	return SignContractResponse{
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
