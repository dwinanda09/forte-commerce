package router

import (
	"github.com/dwinanda09/forte-commerce/internal/handler"
	"github.com/dwinanda09/forte-commerce/internal/middleware"
	"github.com/labstack/echo/v4"
)

func Setup(
	e *echo.Echo,
	authH *handler.AuthHandler,
	productH *handler.ProductHandler,
	checkoutH *handler.CheckoutHandler,
	campaignH *handler.CampaignHandler,
	jwtSecret string,
) {
	// Global middleware
	e.Use(middleware.RequestID)
	e.Use(middleware.LoggerMiddleware)

	// Public routes
	auth := e.Group("/api/v1/auth")
	auth.POST("/login", authH.Login)

	// Protected routes
	api := e.Group("/api/v1")
	api.Use(middleware.JWTAuth(jwtSecret))

	// Products
	api.GET("/products", productH.List)

	// Checkout
	api.POST("/checkout", checkoutH.Submit)
	api.GET("/checkout/:id", checkoutH.GetSession)
	api.POST("/checkout/:id/confirm", checkoutH.Confirm)

	// Orders
	api.GET("/orders", checkoutH.ListOrders)
	api.GET("/orders/:id", checkoutH.GetOrder)
	api.POST("/orders/:id/pay", checkoutH.PayOrder)
	api.POST("/orders/:id/cancel", checkoutH.CancelOrder)

	// Seller-only routes (require role=seller)
	seller := api.Group("/seller")
	seller.Use(middleware.SellerOnly)
	seller.POST("/products", productH.Create)
	seller.PUT("/products/:id", productH.Update)
	seller.DELETE("/products/:id", productH.Delete)

	// Campaigns — public list
	api.GET("/campaigns", campaignH.List)

	// Campaigns — seller CRUD
	seller.GET("/campaigns", campaignH.List)
	seller.GET("/campaigns/:id", campaignH.GetByID)
	seller.POST("/campaigns", campaignH.Create)
	seller.PUT("/campaigns/:id", campaignH.Update)
	seller.DELETE("/campaigns/:id", campaignH.Delete)
	seller.PATCH("/campaigns/:id/toggle", campaignH.Toggle)

	// API docs — static swagger UI at /docs/
	e.Static("/docs", "docs/swagger")
	e.GET("/docs", func(c echo.Context) error {
		return c.Redirect(301, "/docs/index.html")
	})
}
