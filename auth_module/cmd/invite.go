package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/sparque/auth_module/internal/database"
	"github.com/sparque/auth_module/internal/services"
	"github.com/spf13/cobra"
)

var inviteCmd = &cobra.Command{
	Use:   "invite",
	Short: "Manage user invitations",
	Long:  `Create, list, and manage user invitations for domain access`,
}

var createInviteCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new user invitation",
	Long:  `Creates a new invitation for a user to join a domain with specified role and permissions`,
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		email, _ := cmd.Flags().GetString("email")
		role, _ := cmd.Flags().GetString("role")
		inviteType, _ := cmd.Flags().GetString("type")
		expiresInHours, _ := cmd.Flags().GetInt("expires-in")

		if domain == "" {
			log.Fatal("❌ --domain is required")
		}

		if email == "" {
			log.Fatal("❌ --email is required")
		}

		if role == "" {
			role = "customer" // Default role
		}

		if inviteType == "" {
			inviteType = "email_with_qr" // Default type
		}

		if expiresInHours == 0 {
			expiresInHours = 24 // Default 24 hours
		}

		createInvitation(domain, email, role, inviteType, expiresInHours)
	},
}

var listInvitesCmd = &cobra.Command{
	Use:   "list",
	Short: "List pending invitations for a domain",
	Long:  `Lists all pending (not yet accepted) invitations for a domain`,
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")

		if domain == "" {
			log.Fatal("❌ --domain is required")
		}

		listInvitations(domain)
	},
}

func init() {
	rootCmd.AddCommand(inviteCmd)

	// Add subcommands
	inviteCmd.AddCommand(createInviteCmd)
	inviteCmd.AddCommand(listInvitesCmd)

	// Flags for create command
	createInviteCmd.Flags().String("domain", "", "Domain name (e.g., example.com)")
	createInviteCmd.Flags().String("email", "", "User email address")
	createInviteCmd.Flags().String("role", "customer", "User role (admin, customer)")
	createInviteCmd.Flags().String("type", "email_with_qr", "Invitation type (email, qr_code, email_with_qr)")
	createInviteCmd.Flags().Int("expires-in", 24, "Expiration time in hours")

	// Flags for list command
	listInvitesCmd.Flags().String("domain", "", "Domain name (e.g., example.com)")
}

func createInvitation(domain, email, role, inviteType string, expiresInHours int) {
	// Connect to database
	db, err := database.Connect(cfg.MongoDB.URI, cfg.MongoDB.Database)
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}
	defer db.Disconnect()

	ctx := context.Background()

	// Create invitation service
	invitationService := services.NewInvitationService(db, cfg)

	// Create invitation request
	req := &services.CreateInvitationRequest{
		Domain:         domain,
		Email:          email,
		Role:           role,
		Type:           inviteType,
		SingleUse:      true,
		ExpiresInHours: expiresInHours,
		CreatedBy:      "cli",
	}

	// Create invitation
	invitation, qrCodeDataURL, err := invitationService.CreateInvitation(ctx, req)
	if err != nil {
		log.Fatalf("❌ Failed to create invitation: %v", err)
	}

	// Build invitation URL
	inviteURL := cfg.App.FrontendURL + "/invite?token=" + invitation.Token

	// Output results
	fmt.Printf("\n✅ Invitation created successfully!\n\n")
	fmt.Printf("Email: %s\n", email)
	fmt.Printf("Role: %s\n", role)
	fmt.Printf("Domain: %s\n", domain)
	fmt.Printf("Expires: %s\n", invitation.ExpiresAt.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("\n📧 Invitation URL: %s\n", inviteURL)

	if qrCodeDataURL != "" {
		fmt.Printf("📱 QR Code: %s\n", qrCodeDataURL)
	}

	fmt.Printf("\n⏰ Expires: %s\n", invitation.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"))
	fmt.Printf("\ntoken: %s\n", invitation.Token)
	fmt.Println("\n---")
	fmt.Println("Next steps:")
	fmt.Println("1. Send the invitation URL to the user via email")
	fmt.Println("2. User clicks the link or scans the QR code")
	fmt.Println("3. User completes signup and gets access")
}

func listInvitations(domain string) {
	// Connect to database
	db, err := database.Connect(cfg.MongoDB.URI, cfg.MongoDB.Database)
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}
	defer db.Disconnect()

	ctx := context.Background()

	// Query pending invitations
	filter := map[string]interface{}{
		"domain": domain,
		"status": "pending",
	}

	cursor, err := db.Invitations.Find(ctx, filter)
	if err != nil {
		log.Fatalf("❌ Failed to query invitations: %v", err)
	}
	defer cursor.Close(ctx)

	// Print header
	fmt.Printf("\n📋 Pending Invitations for domain: %s\n\n", domain)

	count := 0
	for cursor.Next(ctx) {
		var invitation struct {
			Email     *string `bson:"email"`
			Role      string  `bson:"role"`
			Type      string  `bson:"type"`
			CreatedAt string  `bson:"created_at"`
			ExpiresAt string  `bson:"expires_at"`
			Token     string  `bson:"token"`
		}

		if err := cursor.Decode(&invitation); err != nil {
			log.Printf("⚠️  Error decoding invitation: %v", err)
			continue
		}

		count++
		email := "anonymous"
		if invitation.Email != nil {
			email = *invitation.Email
		}

		fmt.Printf("%d. Email: %s\n", count, email)
		fmt.Printf("   Role: %s\n", invitation.Role)
		fmt.Printf("   Type: %s\n", invitation.Type)
		fmt.Printf("   Token: %s\n", invitation.Token)
		fmt.Printf("   Expires: %s\n\n", invitation.ExpiresAt)
	}

	if count == 0 {
		fmt.Println("No pending invitations found.")
	} else {
		fmt.Printf("Total: %d invitation(s)\n", count)
	}
	fmt.Println()
}
