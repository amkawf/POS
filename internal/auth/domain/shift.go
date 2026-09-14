package domain

import (
	"time"

	"github.com/google/uuid"
)

// ShiftStatus merepresentasikan status laci kasir (buka atau sudah tutup buku)
type ShiftStatus string

const (
	ShiftStatusOpen   ShiftStatus = "OPEN"
	ShiftStatusClosed ShiftStatus = "CLOSED"
)

// CashierShift merepresentasikan sesi buka kasir harian
type CashierShift struct {
	ID                 uuid.UUID   `json:"id"`
	CompanyID          uuid.UUID   `json:"company_id"`
	StoreID            uuid.UUID   `json:"store_id"`
	UserID             uuid.UUID   `json:"user_id"`
	UserName           string      `json:"user_name,omitempty"`
	OpenedAt           time.Time   `json:"opened_at"`
	StartingCash       float64     `json:"starting_cash"`                  // Modal uang receh di laci saat buka kasir
	ClosedAt           *time.Time  `json:"closed_at,omitempty"`            // Waktu tutup kasir
	ActualEndingCash   *float64    `json:"actual_ending_cash,omitempty"`   // Hitungan uang fisik nyata saat tutup
	ExpectedEndingCash *float64    `json:"expected_ending_cash,omitempty"` // Uang yang seharusnya ada menurut transaksi POS
	Notes              *string     `json:"notes,omitempty"`
	Status             ShiftStatus `json:"status"`
}

// IsOpen memeriksa apakah sesi kasir sedang aktif
func (s *CashierShift) IsOpen() bool {
	return s.Status == ShiftStatusOpen
}
