package api

import (
	"time"

	"github.com/SamVitebsk/sharetrip-contract/gen"
	"github.com/SamVitebsk/sharetrip-contract/internal/service"
)

func toCreateContractRequest(request gen.CreateContractRequest) service.CreateContractRequest {
	services := make([]service.ContractService, len(request.Services))
	for i, contractService := range request.Services {
		services[i] = service.ContractService{
			ServiceCode: string(contractService.ServiceCode),
			Allowed:     contractService.Allowed,
		}
	}
	return service.CreateContractRequest{
		ClientID:   request.ClientId,
		ValidFrom:  request.ValidFrom,
		ValidUntil: copyContractEndTime(request.ValidUntil),
		Services:   services,
	}
}

func toSignContractRequest(contractID gen.ContractId) service.SignContractRequest {
	return service.SignContractRequest{ContractID: contractID}
}

func toCreateContractResponse(response service.CreateContractResponse) gen.CreateContractResponse {
	return gen.CreateContractResponse{Contract: gen.Contract{
		Id:         response.ID,
		ClientId:   response.ClientID,
		Status:     gen.ContractStatus(response.Status),
		ValidFrom:  response.ValidFrom,
		ValidUntil: copyContractEndTime(response.ValidUntil),
		Services:   toContractServiceResponses(response.Services),
		CreatedAt:  response.CreatedAt,
		UpdatedAt:  response.UpdatedAt,
	}}
}

func toSignContractResponse(response service.SignContractResponse) gen.SignContractResponse {
	return gen.SignContractResponse{Contract: gen.Contract{
		Id:         response.ID,
		ClientId:   response.ClientID,
		Status:     gen.ContractStatus(response.Status),
		ValidFrom:  response.ValidFrom,
		ValidUntil: copyContractEndTime(response.ValidUntil),
		Services:   toContractServiceResponses(response.Services),
		CreatedAt:  response.CreatedAt,
		UpdatedAt:  response.UpdatedAt,
	}}
}

func toContractServiceResponses(services []service.ContractService) []gen.ContractService {
	responses := make([]gen.ContractService, len(services))
	for index, contractService := range services {
		responses[index] = gen.ContractService{
			ServiceCode: gen.ServiceCode(contractService.ServiceCode),
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

func toCheckServiceAvailabilityRequest(request gen.CheckServiceAvailabilityRequest) service.CheckServiceAvailabilityRequest {
	return service.CheckServiceAvailabilityRequest{
		ClientID:    request.ClientId,
		ServiceCode: string(request.ServiceCode),
	}
}

func toCheckServiceAvailabilityResponse(res service.CheckServiceAvailabilityResponse) gen.CheckServiceAvailabilityResponse {
	response := gen.CheckServiceAvailabilityResponse{
		Allowed:    res.Allowed,
		ContractId: res.ContractID,
	}
	if !res.Allowed && res.Reason != "" {
		reason := gen.CheckServiceAvailabilityResponseReason(res.Reason)
		response.Reason = &reason
	}
	return response
}
