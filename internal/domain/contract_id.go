package domain

import (
	"github.com/google/uuid"
)

type ContractID struct {
	value uuid.UUID
}

func NewContractID(value uuid.UUID) (ContractID, error) {
	id := ContractID{value: value}
	if err := id.Validate(); err != nil {
		return ContractID{}, err
	}
	return id, nil
}

func (id ContractID) Value() uuid.UUID {
	return id.value
}

func (id ContractID) Validate() error {
	if id.value == uuid.Nil {
		return ErrInvalidContractID
	}
	return nil
}
