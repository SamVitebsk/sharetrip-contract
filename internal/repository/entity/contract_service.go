package entity

import (
	"time"

	"github.com/google/uuid"
)

type ContractService struct {
	ContractID  uuid.UUID
	ServiceCode string
	Allowed     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
