package domain

type ContractServiceCode string

const (
	ContractServiceCodeTripCreation     ContractServiceCode = "tripCreation"
	ContractServiceCodeTripParticipants ContractServiceCode = "tripParticipants"
	ContractServiceCodeNotifications    ContractServiceCode = "notifications"
)

func (code ContractServiceCode) Validate() error {
	switch code {
	case ContractServiceCodeTripCreation, ContractServiceCodeTripParticipants, ContractServiceCodeNotifications:
		return nil
	default:
		return ErrInvalidContractServiceCode
	}
}
