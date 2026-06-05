package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dwinanda09/forte-commerce/internal/config"
	"github.com/dwinanda09/forte-commerce/internal/handler"
	"github.com/dwinanda09/forte-commerce/internal/promotion"
	"github.com/dwinanda09/forte-commerce/internal/queue"
	"github.com/dwinanda09/forte-commerce/internal/resource"
	"github.com/dwinanda09/forte-commerce/internal/router"
	"github.com/dwinanda09/forte-commerce/internal/usecase"
	"github.com/dwinanda09/forte-commerce/util"
	"github.com/labstack/echo/v4"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup logging
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	// Connect to database
	db, err := resource.NewDB(cfg)
	if err != nil {
		slog.Error("Failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("Connected to database")

	// Connect to RabbitMQ
	mq, err := queue.NewRabbitMQ(cfg.RabbitURL)
	if err != nil {
		slog.Error("Failed to connect to RabbitMQ", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer mq.Close()

	slog.Info("Connected to RabbitMQ")

	// Initialize logger
	logger := util.NewLogger()

	// Initialize repositories
	productRepo := resource.NewProductResource(db, logger)
	checkoutRepo := resource.NewCheckoutResource(db, logger)
	orderRepo := resource.NewOrderResource(db, logger)
	userRepo := resource.NewUserResource(db, logger)
	campaignRepo := resource.NewCampaignResource(db, logger)

	// Initialize promotion engine
	promoEngine := promotion.NewDynamicEngine(campaignRepo, logger)

	// Initialize usecases
	txManager := resource.NewDBTransactor(db)
	checkoutUC := usecase.NewCheckoutUsecase(
		productRepo,
		checkoutRepo,
		orderRepo,
		userRepo,
		promoEngine,
		mq,
		txManager,
		logger,
	)

	authUC := usecase.NewAuthUsecase(userRepo, cfg.JWTSecret, logger)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUC)
	productHandler := handler.NewProductHandler(productRepo)
	checkoutHandler := handler.NewCheckoutHandler(checkoutUC)
	campaignHandler := handler.NewCampaignHandler(campaignRepo)

	// Setup Echo
	e := echo.New()

	// Setup routes
	router.Setup(e, authHandler, productHandler, checkoutHandler, campaignHandler, cfg.JWTSecret)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start workers
	usecase.StartWorker(ctx, checkoutUC, mq, logger)
	usecase.StartExpiryWorker(ctx, checkoutUC, logger)

	slog.Info("Checkout worker started")
	slog.Info("Expiry worker started")

	// Start server in goroutine
	go func() {
		slog.Info("Starting HTTP server", slog.String("port", cfg.Port))
		if err := e.Start(":" + cfg.Port); err != nil {
			slog.Error("HTTP server error", slog.String("error", err.Error()))
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	slog.Info("Shutting down server")

	// Cancel context to stop workers
	cancel()

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		slog.Error("Error during shutdown", slog.String("error", err.Error()))
	}

	slog.Info("Server stopped")
}
