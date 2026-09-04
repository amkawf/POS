package orderhttp

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pos-backend/internal/order/application"
	"pos-backend/internal/order/domain"
)

type Handler struct {
	createOrderUseCase *application.CreateOrderUseCase
	listOrdersUseCase  *application.ListOrdersUseCase
	getOrderUseCase    *application.GetOrderUseCase
	payOrderUseCase    *application.PayOrderUseCase
	deleteOrderUseCase *application.DeleteOrderUseCase
}

func NewHandler(
	createOrderUseCase *application.CreateOrderUseCase,
	listOrdersUseCase *application.ListOrdersUseCase,
	getOrderUseCase *application.GetOrderUseCase,
	payOrderUseCase *application.PayOrderUseCase,
	deleteOrderUseCase *application.DeleteOrderUseCase,
) *Handler {
	return &Handler{
		createOrderUseCase: createOrderUseCase,
		listOrdersUseCase:  listOrdersUseCase,
		getOrderUseCase:    getOrderUseCase,
		payOrderUseCase:    payOrderUseCase,
		deleteOrderUseCase: deleteOrderUseCase,
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

func (h *Handler) ListOrders(c *gin.Context) {
	companyIDStr := c.Query("company_id")
	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_COMPANY_ID",
				"message": "valid company_id query parameter is required",
			},
		})
		return
	}

	storeIDStr := c.Query("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_STORE_ID",
				"message": "valid store_id query parameter is required",
			},
		})
		return
	}

	var status *string
	if s := c.Query("status"); s != "" {
		status = &s
	}

	limit := int32(50)
	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = int32(parsedLimit)
		}
	}

	orders, err := h.listOrdersUseCase.Execute(c.Request.Context(), application.ListOrdersInput{
		CompanyID: companyID,
		StoreID:   storeID,
		Status:    status,
		Limit:     limit,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "LIST_ORDERS_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	res := make([]OrderResponse, 0, len(orders))
	for _, o := range orders {
		res = append(res, toOrderResponse(o))
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": res,
	})
}

func (h *Handler) GetOrder(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ORDER_ID",
				"message": "valid order id is required",
			},
		})
		return
	}

	companyIDStr := c.Query("company_id")
	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_COMPANY_ID",
				"message": "valid company_id query parameter is required",
			},
		})
		return
	}

	storeIDStr := c.Query("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_STORE_ID",
				"message": "valid store_id query parameter is required",
			},
		})
		return
	}

	order, err := h.getOrderUseCase.Execute(c.Request.Context(), application.GetOrderInput{
		CompanyID: companyID,
		StoreID:   storeID,
		OrderID:   orderID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "GET_ORDER_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "ORDER_NOT_FOUND",
				"message": "order not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, toOrderResponse(order))
}

func (h *Handler) PayOrder(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ORDER_ID",
				"message": "valid order id is required",
			},
		})
		return
	}

	var req PayOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": err.Error(),
			},
		})
		return
	}

	order, err := h.payOrderUseCase.Execute(c.Request.Context(), application.PayOrderInput{
		CompanyID:     req.CompanyID,
		StoreID:       req.StoreID,
		OrderID:       orderID,
		PaymentMethod: req.PaymentMethod,
		AmountPaid:    req.AmountPaid,
		Notes:         req.Notes,
	})
	if err != nil {
		if err == application.ErrOrderNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "ORDER_NOT_FOUND",
					"message": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "PAY_ORDER_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, toOrderResponse(order))
}

func (h *Handler) DeleteOrder(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ORDER_ID",
				"message": "valid order id is required",
			},
		})
		return
	}

	companyIDStr := c.Query("company_id")
	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_COMPANY_ID",
				"message": "valid company_id query parameter is required",
			},
		})
		return
	}

	storeIDStr := c.Query("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_STORE_ID",
				"message": "valid store_id query parameter is required",
			},
		})
		return
	}

	err = h.deleteOrderUseCase.Execute(c.Request.Context(), application.DeleteOrderInput{
		CompanyID: companyID,
		StoreID:   storeID,
		OrderID:   orderID,
	})
	if err != nil {
		if err == application.ErrOrderNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "ORDER_NOT_FOUND",
					"message": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "DELETE_ORDER_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "order deleted successfully",
	})
}

func toOrderResponse(order *domain.Order) OrderResponse {
	items := make([]OrderItemResponse, 0, len(order.Items))
	for _, it := range order.Items {
		items = append(items, OrderItemResponse{
			ID:             it.ID,
			OrderID:        it.OrderID,
			MenuItemID:     it.MenuItemID,
			ItemName:       it.ItemName,
			SKU:            it.SKU,
			Quantity:       it.Quantity,
			UnitPrice:      it.UnitPrice,
			ModifierAmount: it.ModifierAmount,
			DiscountAmount: it.DiscountAmount,
			TaxAmount:      it.TaxAmount,
			TotalAmount:    it.TotalAmount,
			Notes:          it.Notes,
			Status:         it.Status,
		})
	}

	return OrderResponse{
		ID:             order.ID,
		CompanyID:      order.CompanyID,
		StoreID:        order.StoreID,
		OrderNumber:    order.OrderNumber,
		OrderType:      string(order.OrderType),
		OrderSource:    string(order.OrderSource),
		Status:         string(order.Status),
		CustomerName:   order.CustomerName,
		Subtotal:       order.Subtotal,
		DiscountAmount: order.DiscountAmount,
		TaxAmount:      order.TaxAmount,
		ServiceAmount:  order.ServiceAmount,
		TotalAmount:    order.TotalAmount,
		Notes:          order.Notes,
		OpenedAt:       order.OpenedAt.Format(time.RFC3339),
		Items:          items,
	}
}
