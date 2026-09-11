package inventoryhttp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pos-backend/internal/inventory/application"
	"pos-backend/internal/pkg/httputil"
)

type Handler struct {
	adjustStockUseCase *application.AdjustStockUseCase
}

func NewHandler(adjustStockUseCase *application.AdjustStockUseCase) *Handler {
	return &Handler{
		adjustStockUseCase: adjustStockUseCase,
	}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/inventory/adjust", h.AdjustStock)
}

func (h *Handler) AdjustStock(c *gin.Context) {
	var req AdjustStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "INVALID_PAYLOAD", err.Error())
		return
	}

	if req.Quantity == 0 {
		httputil.BadRequest(c, "INVALID_QUANTITY", "Jumlah penyesuaian stok tidak boleh 0")
		return
	}

	if req.Notes == "" {
		req.Notes = "Penyesuaian stok manual"
	}

	err := h.adjustStockUseCase.Execute(
		c.Request.Context(),
		req.CompanyID,
		req.StoreID,
		req.MenuItemID,
		req.Quantity,
		req.Notes,
	)
	if err != nil {
		httputil.InternalError(c, "ADJUST_STOCK_FAILED", err.Error())
		return
	}

	httputil.JSON(c, http.StatusOK, AdjustStockResponse{
		Message: "Stok berhasil disesuaikan",
	})
}