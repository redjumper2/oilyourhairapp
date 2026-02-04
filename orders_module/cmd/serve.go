package cmd

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sparque/orders_module/internal/database"
	"github.com/sparque/orders_module/internal/handlers"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Orders Module API server",
	Long:  `Starts the HTTP server for order processing, payments, and stock management.`,
	Run:   runServer,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func runServer(cmd *cobra.Command, args []string) {
	// Load configuration
	port := viper.GetString("server.port")
	if port == "" {
		port = "9092"
	}

	mongoURI := viper.GetString("mongodb.uri")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	mongoDBName := viper.GetString("mongodb.database")
	if mongoDBName == "" {
		mongoDBName = "orders_module"
	}

	jwtSecret := viper.GetString("jwt.secret")
	if jwtSecret == "" {
		log.Fatal("JWT secret is required")
	}

	stripeKey := viper.GetString("stripe.secret_key")
	if stripeKey == "" {
		log.Fatal("Stripe secret key is required")
	}

	stripeWebhookSecret := viper.GetString("stripe.webhook_secret")
	if stripeWebhookSecret == "" {
		log.Println("⚠️  Warning: Stripe webhook secret not set - webhooks will not work")
	}

	// Connect to MongoDB
	db, err := database.Connect(mongoURI, mongoDBName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.Close()

	log.Printf("✅ Connected to MongoDB: %s/%s", mongoURI, mongoDBName)

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// Initialize handlers
	orderHandler := handlers.NewOrderHandler(db, jwtSecret, stripeKey, stripeWebhookSecret)

	// API Routes
	api := e.Group("/api/v1")

	// Order routes
	api.POST("/orders", orderHandler.CreateOrder)           // Create order and payment intent
	api.GET("/orders/:id", orderHandler.GetOrder)           // Get order by ID
	api.GET("/orders", orderHandler.ListOrders)             // List orders for user
	api.PATCH("/orders/:id", orderHandler.UpdateOrderDetails) // Update order details (customer, addresses)
	api.POST("/webhooks/stripe", orderHandler.StripeWebhook) // Stripe payment webhook

	// Admin routes (require admin role)
	admin := api.Group("/admin")
	admin.GET("/orders", orderHandler.ListAllOrders)         // List all orders
	admin.PATCH("/orders/:id/status", orderHandler.UpdateOrderStatus) // Update order status

	// Start server
	address := fmt.Sprintf(":%s", port)
	log.Printf("🚀 Orders Module starting on port %s", port)
	if err := e.Start(address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
