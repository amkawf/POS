package bootstrap

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"pos-backend/internal/config"
	"pos-backend/internal/database"
	httpserver "pos-backend/internal/http"
	menuapplication "pos-backend/internal/menu/application"
	menuhttp "pos-backend/internal/menu/http"
	menurepository "pos-backend/internal/menu/repository"
	menudb "pos-backend/internal/menu/repository/generated"
	"pos-backend/internal/order/application"
	orderhttp "pos-backend/internal/order/http"
	"pos-backend/internal/order/repository"
	orderdb "pos-backend/internal/order/repository/generated"
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

	// Order repository.
	orderRepository := repository.NewPostgresOrderRepository(
		queries,
		txManager,
	)

	// Order application use case.
	createOrderUseCase := application.NewCreateOrderUseCase(
		orderRepository,
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

	return &App{
		Router:   router,
		Config:   cfg,
		Database: db,
	}, nil
}
