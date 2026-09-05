package menuhttp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pos-backend/internal/menu/application"
)

type Handler struct {
	listMenuItemsUseCase  *application.ListMenuItemsUseCase
	listCategoriesUseCase *application.ListCategoriesUseCase
}

func NewHandler(
	listMenuItemsUseCase *application.ListMenuItemsUseCase,
	listCategoriesUseCase *application.ListCategoriesUseCase,
) *Handler {
	return &Handler{
		listMenuItemsUseCase:  listMenuItemsUseCase,
		listCategoriesUseCase: listCategoriesUseCase,
	}
}

func (h *Handler) ListMenuItems(c *gin.Context) {
	companyID, err := uuid.Parse(c.Query("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "company_id must be a valid UUID",
			},
		})
		return
	}

	items, err := h.listMenuItemsUseCase.Execute(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "LIST_MENU_ITEMS_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	response := ListMenuItemsResponse{
		Items: make([]MenuItemResponse, 0, len(items)),
	}

	for _, item := range items {
		categoryIDs := item.CategoryIDs
		if categoryIDs == nil {
			categoryIDs = []uuid.UUID{}
		}

		response.Items = append(response.Items, MenuItemResponse{
			ID:          item.ID,
			SKU:         item.SKU,
			Name:        item.Name,
			Description: item.Description,
			BasePrice:   item.BasePrice,
			CategoryIDs: categoryIDs,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) ListCategories(c *gin.Context) {
	companyID, err := uuid.Parse(c.Query("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "company_id must be a valid UUID",
			},
		})
		return
	}

	categories, err := h.listCategoriesUseCase.Execute(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "LIST_CATEGORIES_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	response := ListCategoriesResponse{
		Categories: make([]CategoryResponse, 0, len(categories)),
	}

	for _, cat := range categories {
		response.Categories = append(response.Categories, CategoryResponse{
			ID:        cat.ID,
			MenuID:    cat.MenuID,
			Name:      cat.Name,
			SortOrder: cat.SortOrder,
		})
	}

	c.JSON(http.StatusOK, response)
}
