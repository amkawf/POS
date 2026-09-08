package tablehttp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pos-backend/internal/pkg/httputil"
	"pos-backend/internal/table/application"
	"pos-backend/internal/table/domain"
)

type Handler struct {
	listTablesUseCase        *application.ListTablesUseCase
	updateTableStatusUseCase *application.UpdateTableStatusUseCase
}

func NewHandler(
	listTablesUseCase *application.ListTablesUseCase,
	updateTableStatusUseCase *application.UpdateTableStatusUseCase,
) *Handler {
	return &Handler{
		listTablesUseCase:        listTablesUseCase,
		updateTableStatusUseCase: updateTableStatusUseCase,
	}
}

// RegisterRoutes registers the table module routes to the provided router group.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/tables", h.ListTables)
	rg.PATCH("/tables/:id/status", h.UpdateTableStatus)
}

func (h *Handler) ListTables(c *gin.Context) {
	companyID, ok := httputil.ParseUUIDQuery(c, "company_id")
	if !ok {
		return
	}

	storeID, ok := httputil.ParseUUIDQuery(c, "store_id")
	if !ok {
		return
	}

	tables, err := h.listTablesUseCase.Execute(c.Request.Context(), companyID, storeID)
	if err != nil {
		httputil.InternalError(c, "LIST_TABLES_FAILED", err.Error())
		return
	}

	response := ListTablesResponse{
		Tables: make([]TableResponse, 0, len(tables)),
	}

	for _, t := range tables {
		response.Tables = append(response.Tables, TableResponse{
			ID:          t.ID,
			CompanyID:   t.CompanyID,
			StoreID:     t.StoreID,
			TableNumber: t.TableNumber,
			Capacity:    t.Capacity,
			Status:      string(t.Status),
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		})
	}

	httputil.JSON(c, http.StatusOK, response)
}

func (h *Handler) UpdateTableStatus(c *gin.Context) {
	tableID, ok := httputil.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	var req UpdateTableStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "INVALID_REQUEST", err.Error())
		return
	}

	if err := h.updateTableStatusUseCase.Execute(
		c.Request.Context(),
		req.CompanyID,
		req.StoreID,
		tableID,
		domain.TableStatus(req.Status),
	); err != nil {
		httputil.InternalError(c, "UPDATE_TABLE_STATUS_FAILED", err.Error())
		return
	}

	httputil.JSON(c, http.StatusOK, gin.H{
		"message": "table status updated successfully",
		"id":      tableID,
		"status":  req.Status,
	})
}
