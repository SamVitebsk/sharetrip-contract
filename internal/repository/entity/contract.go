package entity

import (
	"time"

	"github.com/google/uuid"
)

type Contract struct {
	ID         uuid.UUID
	ClientID   uuid.UUID
	Status     string
	ValidFrom  time.Time
	ValidUntil *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
