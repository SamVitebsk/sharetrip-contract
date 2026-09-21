package repository

import (
	"fmt"
	"time"

	"github.com/SamVitebsk/sharetrip-contract/internal/domain"
	"github.com/SamVitebsk/sharetrip-contract/internal/repository/entity"
)

func toContractEntity(contract domain.Contract) entity.Contract {
	return entity.Contract{
		ID:         contract.ID.Value(),
		ClientID:   contract.ClientID.Value(),
		Status:     string(contract.Status),
		ValidFrom:  contract.ValidFrom,
		ValidUntil: copyContractEndTime(contract.ValidUntil),
		CreatedAt:  contract.CreatedAt,
		UpdatedAt:  contract.UpdatedAt,
	}
}

func toNewContractServiceEntities(contract domain.Contract) ([]entity.ContractService, error) {
	services := make([]entity.ContractService, len(contract.Services))
	for index, service := range contract.Services {
		serviceCode, err := toStoredContractServiceCode(service.ServiceCode)
		if err != nil {
			return nil, err
		}
		services[index] = entity.ContractService{
			ContractID:  contract.ID.Value(),
			ServiceCode: serviceCode,
			Allowed:     service.Allowed,
			CreatedAt:   contract.CreatedAt,
			UpdatedAt:   contract.CreatedAt,
		}
	}
	return services, nil
}

func toDomainContract(contract entity.Contract, services []entity.ContractService) (domain.Contract, error) {
	contractID, err := domain.NewContractID(contract.ID)
	if err != nil {
		return domain.Contract{}, fmt.Errorf("%w: %v", ErrInvalidContractData, err)
	}
	clientID, err := domain.NewClientID(contract.ClientID)
	if err != nil {
		return domain.Contract{}, fmt.Errorf("%w: %v", ErrInvalidContractData, err)
	}
	status := domain.ContractStatus(contract.Status)
	if err := status.Validate(); err != nil {
		return domain.Contract{}, fmt.Errorf("%w: %v", ErrInvalidContractData, err)
	}
	if contract.ValidFrom.IsZero() || contract.CreatedAt.IsZero() || contract.UpdatedAt.IsZero() {
		return domain.Contract{}, fmt.Errorf("%w: обязательные даты не заданы", ErrInvalidContractData)
	}
	if contract.ValidUntil != nil && !contract.ValidUntil.After(contract.ValidFrom) {
		return domain.Contract{}, fmt.Errorf("%w: некорректный период", ErrInvalidContractData)
	}
	if len(services) == 0 {
		return domain.Contract{}, fmt.Errorf("%w: отсутствуют услуги", ErrInvalidContractData)
	}

	domainServices := make([]domain.ContractService, len(services))
	usedServiceCodes := make(map[domain.ContractServiceCode]struct{}, len(services))
	for index, service := range services {
		if service.ContractID != contract.ID {
			return domain.Contract{}, fmt.Errorf("%w: услуга относится к другому договору", ErrInvalidContractData)
		}
		serviceCode, err := toDomainContractServiceCode(service.ServiceCode)
		if err != nil {
			return domain.Contract{}, err
		}
		if _, exists := usedServiceCodes[serviceCode]; exists {
			return domain.Contract{}, fmt.Errorf("%w: повтор услуги", ErrInvalidContractData)
		}
		usedServiceCodes[serviceCode] = struct{}{}
		domainServices[index] = domain.ContractService{ServiceCode: serviceCode, Allowed: service.Allowed}
	}

	return domain.Contract{
		ID:         contractID,
		ClientID:   clientID,
		Status:     status,
		ValidFrom:  contract.ValidFrom,
		ValidUntil: copyContractEndTime(contract.ValidUntil),
		Services:   domainServices,
		CreatedAt:  contract.CreatedAt,
		UpdatedAt:  contract.UpdatedAt,
	}, nil
}

func toStoredContractServiceCode(code domain.ContractServiceCode) (string, error) {
	switch code {
	case domain.ContractServiceCodeTripCreation:
		return "trip_creation", nil
	case domain.ContractServiceCodeTripParticipants:
		return "trip_participants", nil
	case domain.ContractServiceCodeNotifications:
		return "notifications", nil
	default:
		return "", domain.ErrInvalidContractServiceCode
	}
}

func toDomainContractServiceCode(code string) (domain.ContractServiceCode, error) {
	switch code {
	case "trip_creation":
		return domain.ContractServiceCodeTripCreation, nil
	case "trip_participants":
		return domain.ContractServiceCodeTripParticipants, nil
	case "notifications":
		return domain.ContractServiceCodeNotifications, nil
	default:
		return "", fmt.Errorf("%w: неизвестный код услуги %q", ErrInvalidContractData, code)
	}
}

func copyContractEndTime(validUntil *time.Time) *time.Time {
	if validUntil == nil {
		return nil
	}
	validUntilCopy := *validUntil
	return &validUntilCopy
}
