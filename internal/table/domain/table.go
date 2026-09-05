package domain

import (
	"time"

	"github.com/google/uuid"
)

type TableStatus string

const (
	TableStatusAvailable TableStatus = "AVAILABLE"
	TableStatusOccupied  TableStatus = "OCCUPIED"
	TableStatusReserved  TableStatus = "RESERVED"
	TableStatusInactive  TableStatus = "INACTIVE"
)

type Table struct {
	ID          uuid.UUID
	CompanyID   uuid.UUID
	StoreID     uuid.UUID
	TableNumber string
	Capacity    int32
	Status      TableStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
