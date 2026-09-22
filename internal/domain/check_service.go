package domain

import "time"

type ServiceAvailabilityReason string

const (
	ReasonSuccess            ServiceAvailabilityReason = "success"
	ReasonContractNotFound   ServiceAvailabilityReason = "contractNotFound"
	ReasonContractNotActive  ServiceAvailabilityReason = "contractNotActive"
	ReasonContractNotStarted ServiceAvailabilityReason = "contractNotStarted"
	ReasonContractExpired    ServiceAvailabilityReason = "contractExpired"
	ReasonServiceNotAllowed  ServiceAvailabilityReason = "serviceNotAllowed"
)

func (c Contract) CheckService(code ContractServiceCode, now time.Time) (bool, ServiceAvailabilityReason) {
	if c.Status != ContractStatusActive {
		return false, ReasonContractNotActive
	}
	if now.Before(c.ValidFrom) {
		return false, ReasonContractNotStarted
	}
	if c.ValidUntil != nil && now.After(*c.ValidUntil) {
		return false, ReasonContractExpired
	}

	for _, s := range c.Services {
		if s.ServiceCode != code {
			continue
		}
		if s.Allowed {
			return true, ReasonSuccess
		}
		break
	}

	return false, ReasonServiceNotAllowed
}
