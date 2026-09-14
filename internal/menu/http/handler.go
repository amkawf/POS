package menuhttp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pos-backend/internal/menu/application"
	"pos-backend/internal/pkg/httputil"
)

type Handler struct {
	listMenuItemsUseCase  *application.ListMenuItemsUseCase
	createMenuItemUseCase *application.CreateMenuItemUseCase
	listCategoriesUseCase *application.ListCategoriesUseCase
}

func NewHandler(
	listMenuItemsUseCase *application.ListMenuItemsUseCase,
	createMenuItemUseCase *application.CreateMenuItemUseCase,
	listCategoriesUseCase *application.ListCategoriesUseCase,
) *Handler {
	return &Handler{
		listMenuItemsUseCase:  listMenuItemsUseCase,
		createMenuItemUseCase: createMenuItemUseCase,
		listCategoriesUseCase: listCategoriesUseCase,
	}
}

// RegisterRoutes registers the menu module routes to the provided router group.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/menu-items", h.ListMenuItems)
	rg.POST("/menu-items", h.CreateMenuItem)
	rg.GET("/menu-categories", h.ListCategories)
}

func (h *Handler) ListMenuItems(c *gin.Context) {
	companyID, ok := httputil.ParseUUIDQuery(c, "company_id")
	if !ok {
		return
	}

	// Baca query param store_id jika kasir mengirimkannya
	var storeID *uuid.UUID
	if storeIDStr := c.Query("store_id"); storeIDStr != "" {
		if parsed, err := uuid.Parse(storeIDStr); err == nil {
			storeID = &parsed
		}
	}

	// Panggil use case dengan storeID opsional
	items, err := h.listMenuItemsUseCase.Execute(c.Request.Context(), companyID, storeID)
	if err != nil {
		httputil.InternalError(c, "LIST_MENU_ITEMS_FAILED", err.Error())
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
			ID:              item.ID,
			SKU:             item.SKU,
			Name:            item.Name,
			Description:     item.Description,
			BasePrice:       item.BasePrice,
			CategoryIDs:     categoryIDs,
			Stock:           item.Stock,
			FulfillmentType: item.FulfillmentType,
		})
	}

	httputil.JSON(c, http.StatusOK, response)
}

func (h *Handler) ListCategories(c *gin.Context) {
	companyID, ok := httputil.ParseUUIDQuery(c, "company_id")
	if !ok {
		return
	}

	categories, err := h.listCategoriesUseCase.Execute(c.Request.Context(), companyID)
	if err != nil {
		httputil.InternalError(c, "LIST_CATEGORIES_FAILED", err.Error())
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

	httputil.JSON(c, http.StatusOK, response)
}

// Tambahkan fungsi CreateMenuItem di bagian bawah handler.go:
func (h *Handler) CreateMenuItem(c *gin.Context) {
	var req CreateMenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "INVALID_PAYLOAD", err.Error())
		return
	}
	item, err := h.createMenuItemUseCase.Execute(
		c.Request.Context(),
		req.CompanyID,
		req.SKU,
		req.Name,
		req.Description,
		req.BasePrice,
		req.CategoryIDs,
		req.FulfillmentType,
	)
	if err != nil {
		httputil.InternalError(c, "CREATE_MENU_ITEM_FAILED", err.Error())
		return
	}
	httputil.JSON(c, http.StatusCreated, item)
}
