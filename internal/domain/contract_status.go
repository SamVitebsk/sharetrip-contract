package domain

type ContractStatus string

const (
	ContractStatusDraft      ContractStatus = "draft"
	ContractStatusActive     ContractStatus = "active"
	ContractStatusSuspended  ContractStatus = "suspended"
	ContractStatusTerminated ContractStatus = "terminated"
)

func (status ContractStatus) Validate() error {
	switch status {
	case ContractStatusDraft, ContractStatusActive, ContractStatusSuspended, ContractStatusTerminated:
		return nil
	default:
		return ErrInvalidContractStatus
	}
}
