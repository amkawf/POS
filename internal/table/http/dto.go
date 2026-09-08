package tablehttp

import (
	"time"

	"github.com/google/uuid"
)

type TableResponse struct {
	ID          uuid.UUID `json:"id"`
	CompanyID   uuid.UUID `json:"company_id"`
	StoreID     uuid.UUID `json:"store_id"`
	TableNumber string    `json:"table_number"`
	Capacity    int32     `json:"capacity"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ListTablesResponse struct {
	Tables []TableResponse `json:"tables"`
}

type UpdateTableStatusRequest struct {
	CompanyID uuid.UUID `json:"company_id" binding:"required"`
	StoreID   uuid.UUID `json:"store_id" binding:"required"`
	Status    string    `json:"status" binding:"required"`
}
