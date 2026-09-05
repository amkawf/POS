package tablehttp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pos-backend/internal/table/application"
)

type Handler struct {
	listTablesUseCase *application.ListTablesUseCase
}

func NewHandler(
	listTablesUseCase *application.ListTablesUseCase,
) *Handler {
	return &Handler{
		listTablesUseCase: listTablesUseCase,
	}
}

func (h *Handler) ListTables(c *gin.Context) {
	companyIDStr := c.Query("company_id")
	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "company_id must be a valid UUID",
			},
		})
		return
	}

	storeIDStr := c.Query("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "store_id must be a valid UUID",
			},
		})
		return
	}

	tables, err := h.listTablesUseCase.Execute(c.Request.Context(), companyID, storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "LIST_TABLES_FAILED",
				"message": err.Error(),
			},
		})
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

	c.JSON(http.StatusOK, response)
}
