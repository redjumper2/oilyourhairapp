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
	"github.com/sparque/auth_module/config"
	"github.com/sparque/auth_module/internal/database"
	"github.com/sparque/auth_module/internal/handlers"
	"github.com/sparque/auth_module/internal/middleware"
	"github.com/sparque/auth_module/internal/services"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP API server",
	Long:  `Starts the Echo HTTP server with all authentication endpoints`,
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func runServer() {
	// Initialize database connection
	db, err := database.Connect(cfg.MongoDB.URI, cfg.MongoDB.Database)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.Disconnect()

	log.Println("✅ Connected to MongoDB")

	// Initialize OAuth providers (Goth)
	if cfg.Google.ClientID != "" && cfg.Google.ClientSecret != "" {
		services.InitializeProviders(cfg)
		log.Println("✅ OAuth providers initialized (Google)")
	} else {
		log.Println("⚠️  Google OAuth not configured (set GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET)")
	}

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

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
	// Import handlers
	authHandler := handlers.NewAuthHandler(db, cfg)
	oauthHandler := handlers.NewOAuthHandler(db, cfg)
	adminHandler := handlers.NewAdminHandler(db, cfg)
	authMiddleware := middleware.AuthMiddleware(cfg)
	requireAdmin := middleware.RequireRole("admin")

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "healthy",
			"env":    cfg.Server.Env,
		})
	})

	// API v1
	v1 := e.Group("/api/v1")

	// Auth routes (public)
	auth := v1.Group("/auth")

	// Magic link auth
	auth.POST("/magic-link/request", authHandler.RequestMagicLink)
	auth.GET("/magic-link/verify", authHandler.VerifyMagicLink)

	// Invitation auth
	auth.GET("/invitation/verify", authHandler.VerifyInvitation)
	auth.POST("/invitation/accept", authHandler.AcceptInvitation)

	// Google OAuth
	auth.GET("/google", oauthHandler.GoogleLogin)
	auth.GET("/google/callback", oauthHandler.GoogleCallback)
	auth.GET("/google/callback/json", oauthHandler.GoogleCallbackJSON)

	// Protected auth routes
	auth.GET("/me", authHandler.GetMe, authMiddleware)

	// Admin routes (protected, admin only)
	admin := v1.Group("/admin", authMiddleware, requireAdmin)

	// Domain settings
	admin.GET("/domain/settings", adminHandler.GetDomainSettings)
	admin.PUT("/domain/settings", adminHandler.UpdateDomainSettings)

	// User management
	admin.GET("/users", adminHandler.ListUsers)
	admin.POST("/users/invite", adminHandler.InviteUser)
	admin.PUT("/users/:id", adminHandler.UpdateUser)
	admin.DELETE("/users/:id", adminHandler.DeleteUser)

	// Permissions
	admin.GET("/permissions", adminHandler.GetPermissions)
	admin.GET("/permissions/roles", adminHandler.GetRolePermissions)

	log.Println("✅ Routes configured")
}
