package domain

import "time"

type Contract struct {
	ID         ContractID
	ClientID   ClientID
	Status     ContractStatus
	ValidFrom  time.Time
	ValidUntil *time.Time
	Services   []ContractService
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type ContractService struct {
	ServiceCode ContractServiceCode
	Allowed     bool
}
