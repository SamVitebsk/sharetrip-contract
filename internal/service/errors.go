package service

import (
	"errors"

	"github.com/SamVitebsk/sharetrip-contract/internal/domain"
)

var (
	ErrInvalidInput = errors.New("некорректные входные данные")
	ErrConflict     = errors.New("конфликт операции с договором")
	ErrNotFound     = errors.New("договор не найден")
)

func mapCreateContractError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidClientID),
		errors.Is(err, domain.ErrInvalidContractStart),
		errors.Is(err, domain.ErrEmptyContractServices),
		errors.Is(err, domain.ErrInvalidContractServiceCode):
		return errors.Join(ErrInvalidInput, err)
	case errors.Is(err, domain.ErrInvalidContractPeriod),
		errors.Is(err, domain.ErrDuplicateContractService):
		return errors.Join(ErrConflict, err)
	default:
		return err
	}
}

func mapSignContractError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidContractID):
		return errors.Join(ErrInvalidInput, err)
	case errors.Is(err, domain.ErrContractNotDraft), errors.Is(err, domain.ErrContractExpired):
		return errors.Join(ErrConflict, err)
	default:
		return err
	}
}
