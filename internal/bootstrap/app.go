package bootstrap

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"pos-backend/internal/config"
	"pos-backend/internal/database"
	httpserver "pos-backend/internal/http"
	"pos-backend/internal/order/application"
	orderdb "pos-backend/internal/order/repository/generated"
	"pos-backend/internal/order/repository"
	
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

	// HTTP router.
	router := httpserver.NewRouter()

	_ = createOrderUseCase

	return &App{
		Router:   router,
		Config:   cfg,
		Database: db,
	}, nil
}