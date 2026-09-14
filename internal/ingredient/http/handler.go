package ingredienthttp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pos-backend/internal/ingredient/application"
	"pos-backend/internal/pkg/httputil"

	"github.com/google/uuid" 
	
)

type Handler struct {
	createIngredientUseCase  *application.CreateIngredientUseCase
	listIngredientsUseCase   *application.ListIngredientsUseCase
	restockIngredientUseCase *application.RestockIngredientUseCase
	updateIngredientUseCase  *application.UpdateIngredientUseCase
}

func NewHandler(
	createIngredientUseCase *application.CreateIngredientUseCase,
	listIngredientsUseCase *application.ListIngredientsUseCase,
	restockIngredientUseCase *application.RestockIngredientUseCase,
	updateIngredientUseCase *application.UpdateIngredientUseCase,
) *Handler {
	return &Handler{
		createIngredientUseCase:  createIngredientUseCase,
		listIngredientsUseCase:   listIngredientsUseCase,
		restockIngredientUseCase: restockIngredientUseCase,
		updateIngredientUseCase:  updateIngredientUseCase,
	}
}

// RegisterRoutes mendaftarkan endpoint bahan baku ke router utama Gin
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/ingredients", h.ListIngredients)
	rg.POST("/ingredients", h.CreateIngredient)
	rg.POST("/ingredients/restock", h.RestockIngredient)
	rg.PUT("/ingredients/:id", h.UpdateIngredient)
}

// ListIngredients menangani GET /api/v1/ingredients?company_id=...&store_id=...
func (h *Handler) ListIngredients(c *gin.Context) {
	companyID, ok := httputil.ParseUUIDQuery(c, "company_id")
	if !ok {
		return
	}

	storeID, ok := httputil.ParseUUIDQuery(c, "store_id")
	if !ok {
		return
	}

	items, err := h.listIngredientsUseCase.Execute(c.Request.Context(), companyID, storeID)
	if err != nil {
		httputil.InternalError(c, "LIST_INGREDIENTS_FAILED", err.Error())
		return
	}

	httputil.JSON(c, http.StatusOK, ListIngredientsResponse{
		Ingredients: items,
	})
}

// CreateIngredient menangani POST /api/v1/ingredients (Pendaftaran Bahan Baku Baru)
func (h *Handler) CreateIngredient(c *gin.Context) {
	var req CreateIngredientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "INVALID_PAYLOAD", err.Error())
		return
	}

	ing, err := h.createIngredientUseCase.Execute(
		c.Request.Context(),
		req.CompanyID,
		req.Code,
		req.Name,
		req.Unit,
		req.MinStockAlert,
	)
	if err != nil {
		httputil.InternalError(c, "CREATE_INGREDIENT_FAILED", err.Error())
		return
	}

	httputil.JSON(c, http.StatusCreated, ing)
}

// RestockIngredient menangani POST /api/v1/ingredients/restock (Pencatatan Belanja Bahan Masuk)
func (h *Handler) RestockIngredient(c *gin.Context) {
	var req RestockIngredientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "INVALID_PAYLOAD", err.Error())
		return
	}

	if req.Quantity <= 0 {
		httputil.BadRequest(c, "INVALID_QUANTITY", "Kuantitas belanja harus lebih dari 0")
		return
	}

	if req.Notes == "" {
		req.Notes = "Pembelian bahan baku (Restock)"
	}

	err := h.restockIngredientUseCase.Execute(
		c.Request.Context(),
		req.CompanyID,
		req.StoreID,
		req.IngredientID,
		req.Quantity,
		req.Notes,
	)
	if err != nil {
		httputil.InternalError(c, "RESTOCK_FAILED", err.Error())
		return
	}

	httputil.JSON(c, http.StatusOK, gin.H{
		"message": "Stok bahan baku berhasil ditambahkan",
	})
}

type UpdateIngredientRequest struct {
    CompanyID     uuid.UUID `json:"company_id" binding:"required"`
    Code          *string   `json:"code"`
    Name          string    `json:"name" binding:"required"`
    Unit          string    `json:"unit" binding:"required"`
    MinStockAlert float64   `json:"min_stock_alert"`
}

func (h *Handler) UpdateIngredient(c *gin.Context) {
    idStr := c.Param("id")
    ingredientID, err := uuid.Parse(idStr)
    if err != nil {
        httputil.BadRequest(c, "INVALID_ID", "ID bahan tidak valid")
        return
    }

    var req UpdateIngredientRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        httputil.BadRequest(c, "INVALID_PAYLOAD", err.Error())
        return
    }

    err = h.updateIngredientUseCase.Execute(
        c.Request.Context(),
        req.CompanyID,
        ingredientID,
        req.Code,
        req.Name,
        req.Unit,
        req.MinStockAlert,
    )
    if err != nil {
        httputil.InternalError(c, "UPDATE_INGREDIENT_FAILED", err.Error())
        return
    }

    httputil.JSON(c, http.StatusOK, gin.H{
        "message": "Bahan baku berhasil diperbarui",
    })
}