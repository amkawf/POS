package ingredienthttp

import (
	"github.com/google/uuid"
	"pos-backend/internal/ingredient/domain"
)

// CreateIngredientRequest mengatur payload JSON saat admin mendaftarkan bahan baku baru
type CreateIngredientRequest struct {
	CompanyID     uuid.UUID `json:"company_id" binding:"required"` // ID perusahaan pemilik bahan
	Code          *string   `json:"code"`                          // Kode SKU bahan (opsional, misal: ING-AYAM-01)
	Name          string    `json:"name" binding:"required"`       // Nama bahan (misal: "Daging Ayam Fillet")
	Unit          string    `json:"unit" binding:"required"`       // Satuan bebas (misal: "kg", "gram", "liter")
	MinStockAlert float64   `json:"min_stock_alert"`               // Batas minimum stok sebelum memunculkan peringatan
}

// RestockIngredientRequest mengatur payload JSON saat manajer mencatat belanja bahan baku masuk
type RestockIngredientRequest struct {
	CompanyID    uuid.UUID `json:"company_id" binding:"required"`    // ID perusahaan
	StoreID      uuid.UUID `json:"store_id" binding:"required"`      // Cabang toko tempat bahan disimpan
	IngredientID uuid.UUID `json:"ingredient_id" binding:"required"` // ID bahan yang dibeli
	Quantity     float64   `json:"quantity" binding:"required"`      // Jumlah kuantitas belanja (misal: 10.5 kg atau 500 gram)
	Notes        string    `json:"notes"`                            // Catatan bon pembelian/supplier
}

// ListIngredientsResponse membungkus daftar bahan baku yang dikirim kembali ke web kasir
type ListIngredientsResponse struct {
	Ingredients []domain.Ingredient `json:"ingredients"`
}