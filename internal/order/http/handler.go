package orderhttp

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"pos-backend/internal/order/application"
	"pos-backend/internal/order/domain"
	"pos-backend/internal/pkg/httputil"
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

// RegisterRoutes registers the order module routes to the provided router group.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/orders", h.CreateOrder)
	rg.GET("/orders", h.ListOrders)
	rg.GET("/orders/:id", h.GetOrder)
	rg.POST("/orders/:id/pay", h.PayOrder)
	rg.DELETE("/orders/:id", h.DeleteOrder)
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "INVALID_REQUEST", err.Error())
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
			Notes:		item.Notes,
		})
	}

	order, err := h.createOrderUseCase.Execute(
		c.Request.Context(),
		application.CreateOrderInput{
			CompanyID:    req.CompanyID,
			StoreID:      req.StoreID,
			TableID:      req.TableID,
			OrderType:    domain.OrderType(req.OrderType),
			OrderSource:  domain.OrderSource(req.OrderSource),
			CustomerName: req.CustomerName,
			Notes:        req.Notes,
			CreatedBy:    req.CreatedBy,
			Items:        items,
		},
	)
	if err != nil {
		httputil.BadRequest(c, "CREATE_ORDER_FAILED", err.Error())
		return
	}

	httputil.JSON(c, http.StatusCreated, CreateOrderResponse{
		ID:          order.ID,
		CompanyID:   order.CompanyID,
		StoreID:     order.StoreID,
		TableID:     order.TableID,
		OrderNumber: order.OrderNumber,
		Status:      string(order.Status),
		Subtotal:    order.Subtotal,
		TotalAmount: order.TotalAmount,
	})
}

func (h *Handler) ListOrders(c *gin.Context) {
	companyID, ok := httputil.ParseUUIDQuery(c, "company_id")
	if !ok {
		return
	}

	storeID, ok := httputil.ParseUUIDQuery(c, "store_id")
	if !ok {
		return
	}

	status := httputil.ParseOptStringQuery(c, "status")
	limit := httputil.ParseIntQuery(c, "limit", 50)

	orders, err := h.listOrdersUseCase.Execute(c.Request.Context(), application.ListOrdersInput{
		CompanyID: companyID,
		StoreID:   storeID,
		Status:    status,
		Limit:     limit,
	})
	if err != nil {
		httputil.InternalError(c, "LIST_ORDERS_FAILED", err.Error())
		return
	}

	res := make([]OrderResponse, 0, len(orders))
	for _, o := range orders {
		res = append(res, toOrderResponse(o))
	}

	httputil.JSON(c, http.StatusOK, gin.H{
		"orders": res,
	})
}

func (h *Handler) GetOrder(c *gin.Context) {
	orderID, ok := httputil.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	companyID, ok := httputil.ParseUUIDQuery(c, "company_id")
	if !ok {
		return
	}

	storeID, ok := httputil.ParseUUIDQuery(c, "store_id")
	if !ok {
		return
	}

	order, err := h.getOrderUseCase.Execute(c.Request.Context(), application.GetOrderInput{
		CompanyID: companyID,
		StoreID:   storeID,
		OrderID:   orderID,
	})
	if err != nil {
		httputil.InternalError(c, "GET_ORDER_FAILED", err.Error())
		return
	}

	if order == nil {
		httputil.NotFound(c, "ORDER_NOT_FOUND", "order not found")
		return
	}

	httputil.JSON(c, http.StatusOK, toOrderResponse(order))
}

func (h *Handler) PayOrder(c *gin.Context) {
	orderID, ok := httputil.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	var req PayOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "INVALID_REQUEST", err.Error())
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
			httputil.NotFound(c, "ORDER_NOT_FOUND", err.Error())
			return
		}

		httputil.BadRequest(c, "PAY_ORDER_FAILED", err.Error())
		return
	}

	httputil.JSON(c, http.StatusOK, toOrderResponse(order))
}

func (h *Handler) DeleteOrder(c *gin.Context) {
	orderID, ok := httputil.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	companyID, ok := httputil.ParseUUIDQuery(c, "company_id")
	if !ok {
		return
	}

	storeID, ok := httputil.ParseUUIDQuery(c, "store_id")
	if !ok {
		return
	}

	err := h.deleteOrderUseCase.Execute(c.Request.Context(), application.DeleteOrderInput{
		CompanyID: companyID,
		StoreID:   storeID,
		OrderID:   orderID,
	})
	if err != nil {
		if err == application.ErrOrderNotFound {
			httputil.NotFound(c, "ORDER_NOT_FOUND", err.Error())
			return
		}

		httputil.BadRequest(c, "DELETE_ORDER_FAILED", err.Error())
		return
	}

	httputil.JSON(c, http.StatusOK, gin.H{
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
		TableID:        order.TableID,
		OrderNumber:    order.OrderNumber,
		OrderType:      string(order.OrderType),
		OrderSource:    string(order.OrderSource),
		Status:         string(order.Status),
		CustomerName:   order.CustomerName,
		CreatedBy:		order.CreatedBy,
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
