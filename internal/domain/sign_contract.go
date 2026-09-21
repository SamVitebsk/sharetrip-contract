package domain

import (
	"slices"
	"time"
)

type SignContractRequest struct {
	Now time.Time
}

type SignContractResponse struct {
	Contract Contract
}

func (c Contract) SignContract(request SignContractRequest) (SignContractResponse, error) {
	if err := c.ID.Validate(); err != nil {
		return SignContractResponse{}, err
	}
	if err := c.ClientID.Validate(); err != nil {
		return SignContractResponse{}, err
	}
	if err := c.Status.Validate(); err != nil {
		return SignContractResponse{}, err
	}
	if request.Now.IsZero() {
		return SignContractResponse{}, ErrInvalidContractTime
	}
	if c.ValidFrom.IsZero() {
		return SignContractResponse{}, ErrInvalidContractStart
	}
	if c.ValidUntil != nil && !c.ValidUntil.After(c.ValidFrom) {
		return SignContractResponse{}, ErrInvalidContractPeriod
	}
	if c.Status != ContractStatusDraft {
		return SignContractResponse{}, ErrContractNotDraft
	}
	if c.ValidUntil != nil && !c.ValidUntil.After(request.Now) {
		return SignContractResponse{}, ErrContractExpired
	}

	c.Services = slices.Clone(c.Services)
	if c.ValidUntil != nil {
		validUntil := *c.ValidUntil
		c.ValidUntil = &validUntil
	}
	c.Status = ContractStatusActive
	c.UpdatedAt = request.Now

	return SignContractResponse{Contract: c}, nil
}
