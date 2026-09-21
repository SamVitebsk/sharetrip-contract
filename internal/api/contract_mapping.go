package api

import (
	"time"

	"github.com/SamVitebsk/sharetrip-contract/gen"
	"github.com/SamVitebsk/sharetrip-contract/internal/service"
)

func toCreateContractCommand(request gen.CreateContractRequest) service.CreateContractCommand {
	services := make([]service.ContractService, len(request.Services))
	for i, contractService := range request.Services {
		services[i] = service.ContractService{
			ServiceCode: string(contractService.ServiceCode),
			Allowed:     contractService.Allowed,
		}
	}
	return service.CreateContractCommand{
		ClientID:   request.ClientId,
		ValidFrom:  request.ValidFrom,
		ValidUntil: copyContractEndTime(request.ValidUntil),
		Services:   services,
	}
}

func toSignContractCommand(contractID gen.ContractId) service.SignContractCommand {
	return service.SignContractCommand{ContractID: contractID}
}

func toCreateContractResponse(result service.CreateContractResult) gen.CreateContractResponse {
	return gen.CreateContractResponse{Contract: gen.Contract{
		Id:         result.ID,
		ClientId:   result.ClientID,
		Status:     gen.ContractStatus(result.Status),
		ValidFrom:  result.ValidFrom,
		ValidUntil: copyContractEndTime(result.ValidUntil),
		Services:   toContractServiceResponses(result.Services),
		CreatedAt:  result.CreatedAt,
		UpdatedAt:  result.UpdatedAt,
	}}
}

func toSignContractResponse(result service.SignContractResult) gen.SignContractResponse {
	return gen.SignContractResponse{Contract: gen.Contract{
		Id:         result.ID,
		ClientId:   result.ClientID,
		Status:     gen.ContractStatus(result.Status),
		ValidFrom:  result.ValidFrom,
		ValidUntil: copyContractEndTime(result.ValidUntil),
		Services:   toContractServiceResponses(result.Services),
		CreatedAt:  result.CreatedAt,
		UpdatedAt:  result.UpdatedAt,
	}}
}

func toContractServiceResponses(services []service.ContractService) []gen.ContractService {
	result := make([]gen.ContractService, len(services))
	for index, contractService := range services {
		result[index] = gen.ContractService{
			ServiceCode: gen.ServiceCode(contractService.ServiceCode),
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
