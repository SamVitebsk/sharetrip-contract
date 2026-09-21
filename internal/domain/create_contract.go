package domain

import (
	"slices"
	"time"
)

type CreateContractRequest struct {
	ID         ContractID
	ClientID   ClientID
	ValidFrom  time.Time
	ValidUntil *time.Time
	Services   []ContractService
	Now        time.Time
}

type CreateContractResponse struct {
	Contract Contract
}

func CreateContract(request CreateContractRequest) (CreateContractResponse, error) {
	if err := request.ID.Validate(); err != nil {
		return CreateContractResponse{}, err
	}
	if err := request.ClientID.Validate(); err != nil {
		return CreateContractResponse{}, err
	}
	if request.Now.IsZero() {
		return CreateContractResponse{}, ErrInvalidContractTime
	}
	if request.ValidFrom.IsZero() {
		return CreateContractResponse{}, ErrInvalidContractStart
	}
	if request.ValidUntil != nil && !request.ValidUntil.After(request.ValidFrom) {
		return CreateContractResponse{}, ErrInvalidContractPeriod
	}
	if len(request.Services) == 0 {
		return CreateContractResponse{}, ErrEmptyContractServices
	}

	usedServiceCodes := make(map[ContractServiceCode]struct{}, len(request.Services))
	for _, service := range request.Services {
		if err := service.ServiceCode.Validate(); err != nil {
			return CreateContractResponse{}, err
		}
		if _, exists := usedServiceCodes[service.ServiceCode]; exists {
			return CreateContractResponse{}, ErrDuplicateContractService
		}
		usedServiceCodes[service.ServiceCode] = struct{}{}
	}

	var validUntil *time.Time
	if request.ValidUntil != nil {
		validUntilCopy := *request.ValidUntil
		validUntil = &validUntilCopy
	}

	return CreateContractResponse{Contract: Contract{
		ID:         request.ID,
		ClientID:   request.ClientID,
		Status:     ContractStatusDraft,
		ValidFrom:  request.ValidFrom,
		ValidUntil: validUntil,
		Services:   slices.Clone(request.Services),
		CreatedAt:  request.Now,
		UpdatedAt:  request.Now,
	}}, nil
}
