package orderhttp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pos-backend/internal/order/application"
	"pos-backend/internal/order/domain"
)

type Handler struct {
	createOrderUseCase *application.CreateOrderUseCase
}

func NewHandler(
	createOrderUseCase *application.CreateOrderUseCase,
) *Handler {
	return &Handler{
		createOrderUseCase: createOrderUseCase,
	}
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": err.Error(),
			},
		})
		return
	}

	items := make([]application.CreateOrderItemInput, 0, len(req.Items))

	for _, item := range req.Items {
		items = append(items, application.CreateOrderItemInput{
			MenuItemID: item.MenuItemID,
			ItemName:   item.ItemName,
			SKU:        item.SKU,
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice,
		})
	}

	order, err := h.createOrderUseCase.Execute(
		c.Request.Context(),
		application.CreateOrderInput{
			CompanyID:    req.CompanyID,
			StoreID:      req.StoreID,
			OrderNumber:  req.OrderNumber,
			OrderType:    domain.OrderType(req.OrderType),
			OrderSource:  domain.OrderSource(req.OrderSource),
			CustomerName: req.CustomerName,
			Notes:        req.Notes,
			CreatedBy:    req.CreatedBy,
			Items:        items,
},
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "CREATE_ORDER_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, CreateOrderResponse{
		ID:          order.ID,
		CompanyID:   order.CompanyID,
		StoreID:     order.StoreID,
		OrderNumber: order.OrderNumber,
		Status:      string(order.Status),
		Subtotal:    order.Subtotal,
		TotalAmount: order.TotalAmount,
	})
}