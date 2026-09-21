package domain

import (
	"github.com/google/uuid"
)

type ClientID struct {
	value uuid.UUID
}

func NewClientID(value uuid.UUID) (ClientID, error) {
	id := ClientID{value: value}
	if err := id.Validate(); err != nil {
		return ClientID{}, err
	}
	return id, nil
}

func (id ClientID) Value() uuid.UUID {
	return id.value
}

func (id ClientID) Validate() error {
	if id.value == uuid.Nil {
		return ErrInvalidClientID
	}
	return nil
}
