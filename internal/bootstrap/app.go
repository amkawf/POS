package bootstrap

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"pos-backend/internal/config"
	"pos-backend/internal/database"
	httpserver "pos-backend/internal/http"
	kitchenapplication "pos-backend/internal/kitchen/application"
	kitchenhttp "pos-backend/internal/kitchen/http"
	kitchenrepository "pos-backend/internal/kitchen/repository"
	kitchendb "pos-backend/internal/kitchen/repository/generated"
	menuapplication "pos-backend/internal/menu/application"
	menuhttp "pos-backend/internal/menu/http"
	menurepository "pos-backend/internal/menu/repository"
	menudb "pos-backend/internal/menu/repository/generated"
	"pos-backend/internal/order/application"
	orderhttp "pos-backend/internal/order/http"
	"pos-backend/internal/order/repository"
	orderdb "pos-backend/internal/order/repository/generated"
	tableapplication "pos-backend/internal/table/application"
	tablehttp "pos-backend/internal/table/http"
	tablerepository "pos-backend/internal/table/repository"
	tabledb "pos-backend/internal/table/repository/generated"
)

type App struct {
	Router   *gin.Engine
	Config   config.Config
	Database *pgxpool.Pool
}

func New(
	ctx context.Context,
	cfg config.Config,
) (*App, error) {
	db, err := database.NewPostgresPool(
		ctx,
		database.Config{
			Host:     cfg.Database.Host,
			Port:     cfg.Database.Port,
			Name:     cfg.Database.Name,
			User:     cfg.Database.User,
			Password: cfg.Database.Password,
		},
	)
	if err != nil {
		return nil, err
	}

	// Database infrastructure.
	queries := orderdb.New(db)
	txManager := database.NewTransactionManager(db)

	// Kitchen repository, use cases and HTTP handler.
	kitchenQueries := kitchendb.New(db)
	kitchenRepo := kitchenrepository.NewPostgresKitchenRepository(kitchenQueries, txManager)
	createTicketUseCase := kitchenapplication.NewCreateTicketUseCase(kitchenRepo)
	listTicketsUseCase := kitchenapplication.NewListTicketsUseCase(kitchenRepo)
	updateTicketStatusUseCase := kitchenapplication.NewUpdateTicketStatusUseCase(kitchenRepo)
	kitchenHandler := kitchenhttp.NewHandler(listTicketsUseCase, updateTicketStatusUseCase)

	// Adapter to trigger kitchen tickets when orders are created
	kitchenAdapter := &orderKitchenAdapter{createTicketUseCase: createTicketUseCase}

	// Order repository.
	orderRepository := repository.NewPostgresOrderRepository(
		queries,
		txManager,
	)

	// Order application use cases.
	createOrderUseCase := application.NewCreateOrderUseCase(
		orderRepository,
		kitchenAdapter,
	)
	listOrdersUseCase := application.NewListOrdersUseCase(
		orderRepository,
	)
	getOrderUseCase := application.NewGetOrderUseCase(
		orderRepository,
	)
	payOrderUseCase := application.NewPayOrderUseCase(
		orderRepository,
	)
	deleteOrderUseCase := application.NewDeleteOrderUseCase(
		orderRepository,
	)

	// Order HTTP handler.
	orderHandler := orderhttp.NewHandler(
		createOrderUseCase,
		listOrdersUseCase,
		getOrderUseCase,
		payOrderUseCase,
		deleteOrderUseCase,
	)

	// Menu repository, use case and HTTP handler.
	menuQueries := menudb.New(db)
	menuItemRepository := menurepository.NewPostgresMenuItemRepository(menuQueries)
	categoryRepository := menurepository.NewPostgresCategoryRepository(menuQueries)
	listMenuItemsUseCase := menuapplication.NewListMenuItemsUseCase(menuItemRepository)
	listCategoriesUseCase := menuapplication.NewListCategoriesUseCase(categoryRepository)
	menuHandler := menuhttp.NewHandler(listMenuItemsUseCase, listCategoriesUseCase)

	// Table repository, use case and HTTP handler.
	tableQueries := tabledb.New(db)
	tableRepository := tablerepository.NewPostgresTableRepository(tableQueries)
	listTablesUseCase := tableapplication.NewListTablesUseCase(tableRepository)
	tableHandler := tablehttp.NewHandler(listTablesUseCase)

	// HTTP router.
	router := httpserver.NewRouter()

	apiV1 := router.Group("/api/v1")
	apiV1.POST("/orders", orderHandler.CreateOrder)
	apiV1.GET("/orders", orderHandler.ListOrders)
	apiV1.GET("/orders/:id", orderHandler.GetOrder)
	apiV1.POST("/orders/:id/pay", orderHandler.PayOrder)
	apiV1.DELETE("/orders/:id", orderHandler.DeleteOrder)
	apiV1.GET("/menu-items", menuHandler.ListMenuItems)
	apiV1.GET("/menu-categories", menuHandler.ListCategories)
	apiV1.GET("/tables", tableHandler.ListTables)
	apiV1.GET("/kitchen/tickets", kitchenHandler.ListTickets)
	apiV1.PATCH("/kitchen/tickets/:id/status", kitchenHandler.UpdateTicketStatus)

	return &App{
		Router:   router,
		Config:   cfg,
		Database: db,
	}, nil
}

type orderKitchenAdapter struct {
	createTicketUseCase *kitchenapplication.CreateTicketUseCase
}

func (a *orderKitchenAdapter) CreateTicket(
	ctx context.Context,
	companyID, storeID, orderID uuid.UUID,
	orderNumber, orderType string,
	tableID *uuid.UUID,
	items []application.CreateOrderItemInput,
) error {
	ticketItems := make([]kitchenapplication.TicketItemInput, 0, len(items))
	for _, it := range items {
		ticketItems = append(ticketItems, kitchenapplication.TicketItemInput{
			MenuItemID: it.MenuItemID,
			ItemName:   it.ItemName,
			SKU:        it.SKU,
			Quantity:   it.Quantity,
		})
	}

	_, err := a.createTicketUseCase.Execute(ctx, kitchenapplication.CreateTicketInput{
		CompanyID:   companyID,
		StoreID:     storeID,
		OrderID:     orderID,
		OrderNumber: orderNumber,
		OrderType:   orderType,
		TableID:     tableID,
		Items:       ticketItems,
	})
	return err
}
