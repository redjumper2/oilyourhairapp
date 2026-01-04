package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sparque/auth_module/internal/database"
	"github.com/sparque/auth_module/internal/models"
	"github.com/sparque/auth_module/internal/services"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// domainCmd represents the domain command
var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Manage domains",
	Long:  `Create, list, and manage domains (tenants)`,
}

// domainCreateCmd creates a new domain
var domainCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new domain",
	Long:  `Create a new domain and optionally send magic link invitation to the admin`,
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		name, _ := cmd.Flags().GetString("name")
		adminEmail, _ := cmd.Flags().GetString("admin-email")
		noInvite, _ := cmd.Flags().GetBool("no-invite")

		if domain == "" || name == "" {
			log.Fatal("domain and name are required")
		}

		if !noInvite && adminEmail == "" {
			log.Fatal("admin-email is required when creating invitation (use --no-invite to skip)")
		}

		createDomain(domain, name, adminEmail, noInvite)
	},
}

// domainListCmd lists all domains
var domainListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all domains",
	Long:  `List all registered domains`,
	Run: func(cmd *cobra.Command, args []string) {
		listDomains()
	},
}

// domainDeleteCmd deletes a domain
var domainDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a domain",
	Long:  `Delete a domain and all its users (use with caution!)`,
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")

		if domain == "" {
			log.Fatal("domain is required")
		}

		deleteDomain(domain)
	},
}

func init() {
	rootCmd.AddCommand(domainCmd)

	// Add subcommands
	domainCmd.AddCommand(domainCreateCmd)
	domainCmd.AddCommand(domainListCmd)
	domainCmd.AddCommand(domainDeleteCmd)

	// Flags for create command
	domainCreateCmd.Flags().String("domain", "", "Domain name (e.g., oilyourhair.com)")
	domainCreateCmd.Flags().String("name", "", "Company name (e.g., Oil Your Hair)")
	domainCreateCmd.Flags().String("admin-email", "", "Admin email address")
	domainCreateCmd.Flags().Bool("no-invite", false, "Skip sending invitation email")

	// Flags for delete command
	domainDeleteCmd.Flags().String("domain", "", "Domain name to delete")
}

func createDomain(domainName, companyName, adminEmail string, noInvite bool) {
	// Connect to database
	db, err := database.Connect(cfg.MongoDB.URI, cfg.MongoDB.Database)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.Disconnect()

	ctx := context.Background()

	// Check if domain already exists
	var existing models.Domain
	err = db.Domains.FindOne(ctx, bson.M{"domain": domainName}).Decode(&existing)
	if err == nil {
		log.Fatalf("Domain %s already exists", domainName)
	}

	// Create domain
	domain := models.Domain{
		ID:        primitive.NewObjectID(),
		Domain:    domainName,
		Name:      companyName,
		Status:    "active",
		Settings:  models.DefaultDomainSettings(),
		Branding:  models.DefaultDomainBranding(companyName),
		CreatedAt: time.Now(),
		CreatedBy: "system",
	}

	_, err = db.Domains.InsertOne(ctx, domain)
	if err != nil {
		log.Fatalf("Failed to create domain: %v", err)
	}

	log.Printf("✅ Domain created: %s (%s)", domainName, companyName)

	// Create admin user invitation (unless --no-invite flag is used)
	if !noInvite {
		invitationService := services.NewInvitationService(db, cfg)

		invitation, qrCodeURL, err := invitationService.CreateInvitation(ctx, &services.CreateInvitationRequest{
			Domain:          domainName,
			Email:           adminEmail,
			Role:            "admin",
			Type:            "email_with_qr",
			SingleUse:       true,
			ExpiresInHours:  cfg.Invitation.Defaults.EmailExpiryHours,
		})
		if err != nil {
			log.Fatalf("Failed to create admin invitation: %v", err)
		}

		inviteURL := fmt.Sprintf("%s/invite?token=%s", cfg.App.FrontendURL, invitation.Token)

		log.Printf("✅ Admin invitation created for: %s", adminEmail)
		log.Printf("📧 Invitation URL: %s", inviteURL)
		log.Printf("📱 QR Code: %s", qrCodeURL)
		log.Printf("⏰ Expires: %s", invitation.ExpiresAt.Format(time.RFC3339))

		fmt.Println("\n---")
		fmt.Println("Next steps:")
		fmt.Println("1. Send the invitation URL to the admin via email")
		fmt.Println("2. Admin clicks the link or scans the QR code")
		fmt.Println("3. Admin completes signup and gets access")
	} else {
		log.Printf("ℹ️  Skipping invitation creation (--no-invite flag used)")
	}
}

func listDomains() {
	// Connect to database
	db, err := database.Connect(cfg.MongoDB.URI, cfg.MongoDB.Database)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.Disconnect()

	ctx := context.Background()

	cursor, err := db.Domains.Find(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Failed to list domains: %v", err)
	}
	defer cursor.Close(ctx)

	var domains []models.Domain
	if err = cursor.All(ctx, &domains); err != nil {
		log.Fatalf("Failed to decode domains: %v", err)
	}

	if len(domains) == 0 {
		fmt.Println("No domains found")
		return
	}

	fmt.Printf("\nTotal domains: %d\n\n", len(domains))
	fmt.Printf("%-30s %-30s %-10s %-20s\n", "DOMAIN", "NAME", "STATUS", "CREATED")
	fmt.Println("---------------------------------------------------------------------------------------------------")

	for _, d := range domains {
		fmt.Printf("%-30s %-30s %-10s %-20s\n",
			d.Domain,
			d.Name,
			d.Status,
			d.CreatedAt.Format("2006-01-02 15:04"),
		)
	}
}

func deleteDomain(domainName string) {
	// Connect to database
	db, err := database.Connect(cfg.MongoDB.URI, cfg.MongoDB.Database)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.Disconnect()

	ctx := context.Background()

	// Check if domain exists
	var domain models.Domain
	err = db.Domains.FindOne(ctx, bson.M{"domain": domainName}).Decode(&domain)
	if err != nil {
		log.Fatalf("Domain %s not found", domainName)
	}

	// Delete domain
	_, err = db.Domains.DeleteOne(ctx, bson.M{"domain": domainName})
	if err != nil {
		log.Fatalf("Failed to delete domain: %v", err)
	}

	// Delete all users for this domain
	usersResult, err := db.Users.DeleteMany(ctx, bson.M{"domain": domainName})
	if err != nil {
		log.Printf("Warning: Failed to delete users: %v", err)
	}

	log.Printf("✅ Domain deleted: %s", domainName)
	log.Printf("✅ Deleted %d users", usersResult.DeletedCount)
}
