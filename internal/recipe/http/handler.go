package recipehttp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pos-backend/internal/pkg/httputil"
	"pos-backend/internal/recipe/application"
	"pos-backend/internal/recipe/domain"
)

type Handler struct {
	getRecipeUseCase  *application.GetRecipeUseCase
	saveRecipeUseCase *application.SaveRecipeUseCase
}

func NewHandler(
	getRecipeUseCase *application.GetRecipeUseCase,
	saveRecipeUseCase *application.SaveRecipeUseCase,
) *Handler {
	return &Handler{
		getRecipeUseCase:  getRecipeUseCase,
		saveRecipeUseCase: saveRecipeUseCase,
	}
}

// RegisterRoutes mendaftarkan endpoint resep ke router utama Gin
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/recipes/:menu_item_id", h.GetRecipe)
	rg.POST("/recipes/:menu_item_id", h.SaveRecipe)
}

// GetRecipe menangani GET /api/v1/recipes/:menu_item_id
func (h *Handler) GetRecipe(c *gin.Context) {
	menuItemIDStr := c.Param("menu_item_id")
	menuItemID, err := uuid.Parse(menuItemIDStr)
	if err != nil {
		httputil.BadRequest(c, "INVALID_MENU_ITEM_ID", "ID menu tidak valid")
		return
	}

	items, err := h.getRecipeUseCase.Execute(c.Request.Context(), menuItemID)
	if err != nil {
		httputil.InternalError(c, "GET_RECIPE_FAILED", err.Error())
		return
	}

	httputil.JSON(c, http.StatusOK, RecipeResponse{
		MenuItemID: menuItemID,
		Items:      items,
	})
}

// SaveRecipe menangani POST /api/v1/recipes/:menu_item_id
func (h *Handler) SaveRecipe(c *gin.Context) {
	menuItemIDStr := c.Param("menu_item_id")
	menuItemID, err := uuid.Parse(menuItemIDStr)
	if err != nil {
		httputil.BadRequest(c, "INVALID_MENU_ITEM_ID", "ID menu tidak valid")
		return
	}

	var req SaveRecipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "INVALID_PAYLOAD", err.Error())
		return
	}

	domainItems := make([]domain.RecipeItem, 0, len(req.Items))
	for _, it := range req.Items {
		if it.QuantityPerPortion <= 0 {
			httputil.BadRequest(c, "INVALID_QUANTITY", "Takaran per porsi harus lebih dari 0")
			return
		}
		domainItems = append(domainItems, domain.RecipeItem{
			MenuItemID:         menuItemID,
			IngredientID:       it.IngredientID,
			QuantityPerPortion: it.QuantityPerPortion,
		})
	}

	if err := h.saveRecipeUseCase.Execute(c.Request.Context(), menuItemID, domainItems); err != nil {
		httputil.InternalError(c, "SAVE_RECIPE_FAILED", err.Error())
		return
	}

	httputil.JSON(c, http.StatusOK, gin.H{
		"message": "Resep berhasil disimpan",
	})
}