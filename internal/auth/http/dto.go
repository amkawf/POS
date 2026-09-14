package authhttp

import (
	"github.com/google/uuid"

	"pos-backend/internal/auth/domain"
)

type PinLoginRequest struct {
	StoreID uuid.UUID `json:"store_id" binding:"required"`
	PIN     string    `json:"pin" binding:"required"`
}

type StaffProfileResponse struct {
	ID   uuid.UUID   `json:"id"`
	Name string      `json:"name"`
	Role domain.Role `json:"role"`
}

type OpenShiftRequest struct {
	CompanyID    uuid.UUID `json:"company_id" binding:"required"`
	StoreID      uuid.UUID `json:"store_id" binding:"required"`
	UserID       uuid.UUID `json:"user_id" binding:"required"`
	StartingCash float64   `json:"starting_cash"`
}

type CloseShiftRequest struct {
	ShiftID            uuid.UUID `json:"shift_id" binding:"required"`
	ActualEndingCash   float64   `json:"actual_ending_cash"`
	ExpectedEndingCash float64   `json:"expected_ending_cash"`
	Notes              *string   `json:"notes"`
}
