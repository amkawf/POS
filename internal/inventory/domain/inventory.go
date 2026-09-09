package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

//domain errors
var (
	ErrInsufficientStock = errors.New("stok inventori tidak mencukupi")
	ErrInventoryNotFound = errors.New("data inventori tidak ditemukan")
)

type MovementType string

const (
	MovementTypeSale		MovementType = "SALE"
	MovementTypeRestock		MovementType = "RESTOCK"
	MovementTypeAdjustment	MovementType = "ADJUSTMENT"
	MovementTypeWaste		MovementType = "WASTE"
)

type StoreInventory struct {
	ID			uuid.UUID
	CompanyID	uuid.UUID
	StoreID		uuid.UUID
	MenuItemID	uuid.UUID
	Stock		float64
	CreatedAt	time.Time
	UpdatedAt	time.Time
}

type InventoryMovement struct {
	ID				uuid.UUID
	CompanyID		uuid.UUID
	StoreID			uuid.UUID
	MenuItemID		uuid.UUID
	OrderID			*uuid.UUID
	Quantity		float64
	MovementType	MovementType
	Notes			*string
	CreatedAt		time.Time
}