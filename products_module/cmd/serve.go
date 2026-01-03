package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/sparque/products_module/config"
	"github.com/sparque/products_module/internal/database"
	"github.com/sparque/products_module/internal/handlers"
	"github.com/sparque/products_module/internal/middleware"
	"github.com/sparque/products_module/internal/services"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP API server",
	Long:  `Starts the Echo HTTP server with products API endpoints`,
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func runServer() {
	// Initialize database connection
	db, err := database.Connect(cfg.MongoDB.URI, cfg.MongoDB.Database, cfg.Auth.DomainsDB)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.Disconnect()

	log.Println("✅ Connected to MongoDB")
	log.Printf("   Products DB: %s", cfg.MongoDB.Database)
	log.Printf("   Auth DB: %s (read-only access to domains)", cfg.Auth.DomainsDB)

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: cfg.CORS.AllowedOrigins,
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
	}))

	// Routes
	setupRoutes(e, db, cfg)

	// Start server with graceful shutdown
	go func() {
		addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
		log.Printf("🚀 Starting server on %s", addr)
		if err := e.Start(addr); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Server stopped")
}

func setupRoutes(e *echo.Echo, db *database.DB, cfg *config.Config) {
	// Initialize services
	productService := services.NewProductService(db)
	reviewService := services.NewReviewService(db.Database)
	contactService := services.NewContactService(db.Database)

	// Initialize handlers
	adminHandler := handlers.NewAdminHandler(productService)
	publicHandler := handlers.NewPublicHandler(productService, reviewService, contactService)

	// Initialize middleware
	apiKeyAuth := middleware.APIKeyMiddleware(cfg)
	requireWrite := middleware.RequirePermission("products.write")
	requireRead := middleware.RequirePermission("products.read")

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "healthy",
			"env":    cfg.Server.Env,
		})
	})

	// API v1
	v1 := e.Group("/api/v1")

	// Public API (no auth required)
	public := v1.Group("/public/:domain")
	public.GET("/products", publicHandler.ListProducts)
	public.GET("/products/:id", publicHandler.GetProduct)
	public.GET("/products/search", publicHandler.SearchProducts)
	public.GET("/promotions", publicHandler.GetPromotions)

	// Reviews endpoints
	public.POST("/reviews", publicHandler.CreateReview)
	public.GET("/reviews", publicHandler.ListReviews)
	public.GET("/reviews/:id", publicHandler.GetReview)

	// Contact form endpoint
	public.POST("/contact", publicHandler.CreateContact)

	// Admin API (API key required)
	admin := v1.Group("/products", apiKeyAuth)

	// Create and list require appropriate permissions
	admin.POST("", adminHandler.CreateProduct, requireWrite)
	admin.GET("", adminHandler.ListProducts, requireRead)

	// Individual product operations
	admin.GET("/:id", adminHandler.GetProduct, requireRead)
	admin.PUT("/:id", adminHandler.UpdateProduct, requireWrite)
	admin.DELETE("/:id", adminHandler.DeleteProduct, requireWrite)

	// Stock management
	admin.PUT("/:id/variants/:variantId/stock", adminHandler.UpdateStock, requireWrite)

	log.Println("✅ Routes configured")
	log.Println("   Public API: /api/v1/public/:domain/products")
	log.Println("   Public API: /api/v1/public/:domain/reviews")
	log.Println("   Public API: /api/v1/public/:domain/contact")
	log.Println("   Admin API: /api/v1/products (requires API key)")
}
