package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/sparque/auth_module/internal/database"
	"github.com/sparque/auth_module/internal/services"
	"github.com/spf13/cobra"
)

var apiKeyCmd = &cobra.Command{
	Use:   "apikey",
	Short: "Manage API keys",
	Long:  `Create, list, and revoke service-scoped API keys`,
}

var createAPIKeyCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new API key for a service",
	Long:  `Creates a new service-scoped API key with specified permissions and expiration`,
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		service, _ := cmd.Flags().GetString("service")
		description, _ := cmd.Flags().GetString("description")
		permissions, _ := cmd.Flags().GetStringSlice("permissions")
		expiresInDays, _ := cmd.Flags().GetInt("expires-in")
		createdBy, _ := cmd.Flags().GetString("created-by")

		if domain == "" || service == "" {
			log.Fatal("❌ --domain and --service are required")
		}

		// Connect to database
		db, err := database.Connect(cfg.MongoDB.URI, cfg.MongoDB.Database)
		if err != nil {
			log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
		}
		defer db.Disconnect()

		// Create API key service
		apiKeyService := services.NewAPIKeyService(db, cfg)

		// Create API key
		req := services.CreateAPIKeyRequest{
			Domain:      domain,
			Service:     service,
			Description: description,
			Permissions: permissions,
			ExpiresIn:   time.Duration(expiresInDays) * 24 * time.Hour,
		}

		token, apiKey, err := apiKeyService.CreateAPIKey(context.Background(), req, createdBy)
		if err != nil {
			log.Fatalf("❌ Failed to create API key: %v", err)
		}

		fmt.Println("\n✅ API Key created successfully!")
		fmt.Println("\n⚠️  IMPORTANT: Save this API key securely. It cannot be retrieved again.")
		fmt.Println("\n" + strings.Repeat("=", 80))
		fmt.Printf("API Key: %s\n", token)
		fmt.Println(strings.Repeat("=", 80))
		fmt.Printf("\nKey ID: %s\n", apiKey.KeyID)
		fmt.Printf("Domain: %s\n", apiKey.Domain)
		fmt.Printf("Service: %s\n", apiKey.Service)
		fmt.Printf("Description: %s\n", apiKey.Description)
		fmt.Printf("Permissions: %v\n", apiKey.Permissions)
		fmt.Printf("Expires At: %s\n", apiKey.ExpiresAt.Format(time.RFC3339))
		fmt.Printf("Created At: %s\n", apiKey.CreatedAt.Format(time.RFC3339))
		fmt.Printf("Days Until Expiry: %d\n\n", apiKey.DaysUntilExpiry())
	},
}

var listAPIKeysCmd = &cobra.Command{
	Use:   "list",
	Short: "List API keys for a domain",
	Long:  `Lists all API keys for a domain, optionally filtered by service`,
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		service, _ := cmd.Flags().GetString("service")

		if domain == "" {
			log.Fatal("❌ --domain is required")
		}

		// Connect to database
		db, err := database.Connect(cfg.MongoDB.URI, cfg.MongoDB.Database)
		if err != nil {
			log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
		}
		defer db.Disconnect()

		// Create API key service
		apiKeyService := services.NewAPIKeyService(db, cfg)

		// List API keys
		keys, err := apiKeyService.ListAPIKeys(context.Background(), domain, service)
		if err != nil {
			log.Fatalf("❌ Failed to list API keys: %v", err)
		}

		if len(keys) == 0 {
			fmt.Printf("No API keys found for domain: %s\n", domain)
			return
		}

		fmt.Printf("\n📋 API Keys for domain: %s\n", domain)
		if service != "" {
			fmt.Printf("   (filtered by service: %s)\n", service)
		}
		fmt.Println()

		for i, key := range keys {
			status := "✅ Valid"
			if key.Revoked {
				status = "❌ Revoked"
			} else if key.IsExpired() {
				status = "⏰ Expired"
			}

			fmt.Printf("%d. %s\n", i+1, status)
			fmt.Printf("   Key ID: %s\n", key.KeyID)
			fmt.Printf("   Service: %s\n", key.Service)
			fmt.Printf("   Description: %s\n", key.Description)
			fmt.Printf("   Permissions: %v\n", key.Permissions)
			fmt.Printf("   Expires: %s (%d days)\n", key.ExpiresAt.Format("2006-01-02"), key.DaysUntilExpiry())
			fmt.Printf("   Created: %s\n", key.CreatedAt.Format("2006-01-02 15:04"))
			if key.LastUsedAt != nil {
				fmt.Printf("   Last Used: %s\n", key.LastUsedAt.Format("2006-01-02 15:04"))
			}
			fmt.Println()
		}
	},
}

var revokeAPIKeyCmd = &cobra.Command{
	Use:   "revoke",
	Short: "Revoke an API key",
	Long:  `Revokes an API key by its key ID`,
	Run: func(cmd *cobra.Command, args []string) {
		keyID, _ := cmd.Flags().GetString("key-id")

		if keyID == "" {
			log.Fatal("❌ --key-id is required")
		}

		// Connect to database
		db, err := database.Connect(cfg.MongoDB.URI, cfg.MongoDB.Database)
		if err != nil {
			log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
		}
		defer db.Disconnect()

		// Create API key service
		apiKeyService := services.NewAPIKeyService(db, cfg)

		// Revoke API key
		err = apiKeyService.RevokeAPIKey(context.Background(), keyID)
		if err != nil {
			log.Fatalf("❌ Failed to revoke API key: %v", err)
		}

		fmt.Printf("✅ API key revoked successfully: %s\n", keyID)
	},
}

func init() {
	rootCmd.AddCommand(apiKeyCmd)

	// Create command
	apiKeyCmd.AddCommand(createAPIKeyCmd)
	createAPIKeyCmd.Flags().String("domain", "", "Domain (required)")
	createAPIKeyCmd.Flags().String("service", "", "Service name (e.g., products, orders) (required)")
	createAPIKeyCmd.Flags().String("description", "", "Description of the API key")
	createAPIKeyCmd.Flags().StringSlice("permissions", []string{}, "Permissions (comma-separated)")
	createAPIKeyCmd.Flags().Int("expires-in", 365, "Expiration in days (default: 365)")
	createAPIKeyCmd.Flags().String("created-by", "cli", "User ID who created the key (default: cli)")

	// List command
	apiKeyCmd.AddCommand(listAPIKeysCmd)
	listAPIKeysCmd.Flags().String("domain", "", "Domain (required)")
	listAPIKeysCmd.Flags().String("service", "", "Filter by service (optional)")

	// Revoke command
	apiKeyCmd.AddCommand(revokeAPIKeyCmd)
	revokeAPIKeyCmd.Flags().String("key-id", "", "Key ID to revoke (required)")
}
