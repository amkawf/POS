package bootstrap

import (
	"context"
	"fmt"

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
	tabledomain "pos-backend/internal/table/domain"
	tablehttp "pos-backend/internal/table/http"
	tablerepository "pos-backend/internal/table/repository"
	tabledb "pos-backend/internal/table/repository/generated"
	inventoryrepository "pos-backend/internal/inventory/repository"
	inventoryapplication "pos-backend/internal/inventory/application"
	inventoryhttp "pos-backend/internal/inventory/http"
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
			URL:      cfg.Database.URL,
			Host:     cfg.Database.Host,
			Port:     cfg.Database.Port,
			Name:     cfg.Database.Name,
			User:     cfg.Database.User,
			Password: cfg.Database.Password,
			SSLMode:  cfg.Database.SSLMode,
		},
	)

	if err != nil {
		return nil, err
	}

	if err := database.Ping(ctx, db); err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	// Database infrastructure.
	queries := orderdb.New(db)
	txManager := database.NewTransactionManager(db)

	// Kitchen tickets repository and use cases.
	kitchenQueries := kitchendb.New(db)
	kitchenRepo := kitchenrepository.NewPostgresKitchenRepository(kitchenQueries, txManager)
	createTicketUseCase := kitchenapplication.NewCreateTicketUseCase(kitchenRepo)
	listTicketsUseCase := kitchenapplication.NewListTicketsUseCase(kitchenRepo)
	updateTicketStatusUseCase := kitchenapplication.NewUpdateTicketStatusUseCase(kitchenRepo)
	kitchenHandler := kitchenhttp.NewHandler(listTicketsUseCase, updateTicketStatusUseCase)

	// Adapter to trigger kitchen tickets when orders are created
	kitchenAdapter := &orderKitchenAdapter{createTicketUseCase: createTicketUseCase}

	// Table repository, use case and HTTP handler.
	tableQueries := tabledb.New(db)
	tableRepository := tablerepository.NewPostgresTableRepository(tableQueries)
	listTablesUseCase := tableapplication.NewListTablesUseCase(tableRepository)
	updateTableStatusUseCase := tableapplication.NewUpdateTableStatusUseCase(tableRepository)
	tableHandler := tablehttp.NewHandler(listTablesUseCase, updateTableStatusUseCase)
	tableAdapter := &orderTableAdapter{tableRepo: tableRepository}

	// Order repository.
	orderRepository := repository.NewPostgresOrderRepository(
		queries,
		txManager,
	)
		// Inventory repository & adapter
	inventoryRepository := inventoryrepository.NewPostgresInventoryRepository(db, txManager)
	inventoryAdapter := &orderInventoryAdapter{inventoryRepo: inventoryRepository}

	// Order application use cases.
	createOrderUseCase := application.NewCreateOrderUseCase(
		orderRepository,
		kitchenAdapter,
		tableAdapter,
		inventoryAdapter,
	)

	listOrdersUseCase := application.NewListOrdersUseCase(
		orderRepository,
	)
	getOrderUseCase := application.NewGetOrderUseCase(
		orderRepository,
	)
	payOrderUseCase := application.NewPayOrderUseCase(
		orderRepository,
		tableAdapter,
	)
	deleteOrderUseCase := application.NewDeleteOrderUseCase(
		orderRepository,
		tableAdapter,
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
	listMenuItemsUseCase := menuapplication.NewListMenuItemsUseCase(menuItemRepository, inventoryRepository)
	listCategoriesUseCase := menuapplication.NewListCategoriesUseCase(categoryRepository)
	menuHandler := menuhttp.NewHandler(listMenuItemsUseCase, listCategoriesUseCase)
	// Inventory use case and HTTP handler.
	adjustStockUseCase := inventoryapplication.NewAdjustStockUseCase(inventoryRepository)
	inventoryHandler := inventoryhttp.NewHandler(adjustStockUseCase)

	// HTTP router and modular route registration.
	router := httpserver.NewRouter(httpserver.RouterConfig{
		Database: db,
		Version:  "1.0.0",
	})

	apiV1 := router.Group("/api/v1")
	orderHandler.RegisterRoutes(apiV1)
	menuHandler.RegisterRoutes(apiV1)
	tableHandler.RegisterRoutes(apiV1)
	kitchenHandler.RegisterRoutes(apiV1)
	inventoryHandler.RegisterRoutes(apiV1)

	return &App{
		Router:   router,
		Config:   cfg,
		Database: db,
	}, nil
}

type orderTableAdapter struct {
	tableRepo tablerepository.TableRepository
}

func (a *orderTableAdapter) UpdateTableStatus(
	ctx context.Context,
	companyID, storeID, tableID uuid.UUID,
	status string,
) error {
	return a.tableRepo.UpdateStatus(ctx, companyID, storeID, tableID, tabledomain.TableStatus(status))
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

type orderInventoryAdapter struct {
	inventoryRepo inventoryrepository.InventoryRepository
}

func (a *orderInventoryAdapter) DeductStock(
	ctx context.Context,
	companyID, storeID, orderID uuid.UUID,
	items []application.CreateOrderItemInput,
) error {
	deductions := make([]inventoryrepository.ItemDeduction, 0, len(items))
	for _, it := range items {
		deductions = append(deductions, inventoryrepository.ItemDeduction{
			MenuItemID: it.MenuItemID,
			ItemName:   it.ItemName,
			Quantity:   it.Quantity,
		})
	}
	return a.inventoryRepo.DeductStockForOrder(ctx, companyID, storeID, orderID, deductions)
}
